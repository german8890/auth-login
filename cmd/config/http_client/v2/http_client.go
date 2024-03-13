package http_client

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	URL "net/url"
	"strings"
	"time"

	"autenticacion-ms/cmd/config/errors"
	"autenticacion-ms/cmd/config/http_client"
)

type key int

type keyString string

const (
	ContextKeyGinContext key = iota
	ContextKey
	ContextKeyHttpRequest
	ContextKeyHttpResponse
	ContextKeyRequest
	ContextKeyResponse
	ContextKeyTimeResponse
	ContextKeyAttempt

	ContextKeyPathRequest keyString = "ContextKey/Path_Request"
)

func defaultEncode(_ context.Context, r *http.Request, request interface{}) error {
	r.Header.Set("Content-Type", "application/json")
	if request != nil {
		req, err := json.Marshal(request)
		if err != nil {
			return err
		}
		r.Body = io.NopCloser(strings.NewReader(string(req)))
	}
	return nil
}

func defaultDecode(_ context.Context, response *http.Response) (interface{}, error) {
	return response, nil
}

type HttpClient interface {
	GET(ctx context.Context, url string, headers map[string][]string, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error)
	POST(ctx context.Context, url string, headers map[string][]string, request interface{}, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error)
	PUT(ctx context.Context, url string, headers map[string][]string, request interface{}, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error)
	PATCH(ctx context.Context, url string, headers map[string][]string, request interface{}, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error)
	DELETE(ctx context.Context, url string, headers map[string][]string, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error)
	HEAD(ctx context.Context, url string, headers map[string][]string, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error)
}

type httpClient struct {
	Client    *http.Client
	Client2   http.Client
	attempts  int
	before    []http_client.RequestFunc
	after     []http_client.ClientResponseFunc
	finalizer []http_client.ClientFinalizerFunc
}

func MakeNewHttpClient(client http_client.ClientS) HttpClient {
	return &httpClient{Client: client.C, before: client.Before, after: client.After, finalizer: client.Finalizer, attempts: client.Attempts}
}

func (h httpClient) do(ctx context.Context, method, url string, headers map[string][]string, request interface{}, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc, attempt int) (interface{}, error) {
	var path string
	if v := ctx.Value(ContextKeyPathRequest); v != nil {
		path = v.(string)
	}

	ctx = context.WithValue(context.TODO(), ContextKey, ctx)
	//ctx, cancel := context.WithCancel(ctx)
	var (
		resp         *http.Response
		req          *http.Request
		response     interface{}
		err          error
		timeResponse int64
	)
	if encode == nil {
		encode = defaultEncode
	}
	if decode == nil {
		decode = defaultDecode
	}
	if h.finalizer != nil {
		defer func() {
			//if resp != nil {
			ctx = context.WithValue(ctx, ContextKey, &ctx)
			ctx = context.WithValue(ctx, ContextKeyHttpRequest, req)
			ctx = context.WithValue(ctx, ContextKeyRequest, request)
			ctx = context.WithValue(ctx, ContextKeyHttpResponse, resp)
			ctx = context.WithValue(ctx, ContextKeyResponse, response)
			ctx = context.WithValue(ctx, ContextKeyTimeResponse, timeResponse)
			ctx = context.WithValue(ctx, ContextKeyAttempt, attempt)
			ctx = context.WithValue(ctx, ContextKeyPathRequest, path)
			for _, f := range h.finalizer {
				f(ctx, err)
			}
			//}

		}()
	}
	uri, err := URL.Parse(url)
	if err != nil {
		//cancel()
		return nil, err
	}

	var body []byte = nil
	if request != nil {
		body, err = json.Marshal(request)
		if err != nil {
			return nil, err
		}
	}

	req, err = http.NewRequestWithContext(ctx, method, uri.String(), bytes.NewBuffer(body))

	if err != nil {
		//cancel()
		return nil, err
	}

	if headers != nil {
		req.Header = headers
	}

	for _, f := range h.before {
		ctx = f(ctx, req)
	}
	if err := encode(ctx, req, nil); err != nil {
		return nil, err
	}

	beginTime := time.Now()

	resp, err = h.Client2.Do(req.WithContext(ctx))

	//timeResponse = time.Now().Sub(beginTime).Milliseconds()
	timeResponse = time.Since(beginTime).Milliseconds()
	if err != nil {
		//cancel()
		return nil, err
	}

	defer resp.Body.Close()

	for _, f := range h.after {
		ctx = f(ctx, resp)
	}

	response, err = decode(ctx, resp)
	if err != nil {
		//cancel()
		return nil, err
	}
	return response, err

}

func (h httpClient) doWithRetries(ctx context.Context, method, url string, headers map[string][]string, request interface{}, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error) {
	var (
		resp interface{}
		err  error
	)
	h.Client2 = *h.Client
Loop:
	for i := 0; i <= h.attempts; i++ {
		if i > 0 {
			if i == h.attempts {
				h.Client2.Timeout = 0
			} else {
				// h.Client2.Timeout = time.Duration(int64(i+1)) * h.Client.Timeout
				h.Client2.Timeout = h.Client.Timeout
			}
		}
		attempt := i + 1
		resp, err = h.do(ctx, method, url, headers, request, encode, decode, attempt)

		if err == nil && resp != nil {
			break Loop
		}

		if err != nil {
			switch merr := err.(type) {
			case *URL.Error:
				if merr.Timeout() {
					continue
				}
			case errors.ErrorResponse:
				err2 := err.(errors.ErrorResponse)
				if err2.Status <= 499 && err2.Status > 399 {
					break Loop
				}
			}
			if i == 1 {
				break Loop
			}
		}

	}
	return resp, err
}

func (h httpClient) GET(ctx context.Context, url string, headers map[string][]string, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(ctx, http.MethodGet, url, headers, nil, encode, decode)
}

func (h httpClient) POST(ctx context.Context, url string, headers map[string][]string, request interface{}, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(ctx, http.MethodPost, url, headers, request, encode, decode)
}

func (h httpClient) PUT(ctx context.Context, url string, headers map[string][]string, request interface{}, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(ctx, http.MethodPut, url, headers, request, encode, decode)
}

func (h httpClient) PATCH(ctx context.Context, url string, headers map[string][]string, request interface{}, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(ctx, http.MethodPatch, url, headers, request, encode, decode)
}

func (h httpClient) DELETE(ctx context.Context, url string, headers map[string][]string, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(ctx, http.MethodDelete, url, headers, nil, encode, decode)
}

func (h httpClient) HEAD(ctx context.Context, url string, headers map[string][]string, encode http_client.EncodeRequestFunc, decode http_client.DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(ctx, http.MethodHead, url, headers, nil, encode, decode)
}
