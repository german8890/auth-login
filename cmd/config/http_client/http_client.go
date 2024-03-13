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

	"github.com/gin-gonic/gin"
)

type key int

const (
	ContextKeyGinContext key = iota
	ContextKeyHttpRequest
	ContextKeyHttpResponse
	ContextKeyRequest
	ContextKeyResponse
	ContextKeyTimeResponse
	ContextKeyAttempt
)

var (
	ContextKeyPathRequest = "ContextKey/Path_Request"
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
	GET(c *gin.Context, url string, headers map[string][]string, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error)
	POST(c *gin.Context, url string, headers map[string][]string, request interface{}, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error)
	PUT(c *gin.Context, url string, headers map[string][]string, request interface{}, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error)
	PATCH(c *gin.Context, url string, headers map[string][]string, request interface{}, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error)
	DELETE(c *gin.Context, url string, headers map[string][]string, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error)
	HEAD(c *gin.Context, url string, headers map[string][]string, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error)
}

type httpClient struct {
	Client    *http.Client
	Client2   http.Client
	attempts  int
	before    []RequestFunc
	after     []ClientResponseFunc
	finalizer []ClientFinalizerFunc
	//mutex     sync.RWMutex
}

func MakeNewHttpClient(client ClientS) HttpClient {
	return &httpClient{Client: client.C, before: client.Before, after: client.After, finalizer: client.Finalizer, attempts: client.Attempts}
}

func (h *httpClient) do(c *gin.Context, method, url string, headers map[string][]string, request interface{}, encode EncodeRequestFunc, decode DecodeResponseFunc, attempt int) (interface{}, error) {
	ctx, cancel := context.WithCancel(c.Request.Context())
	defer cancel()
	//h.mutex.Lock()
	ctx = context.WithValue(ctx, gin.ContextKey, c.Copy())
	//h.mutex.Unlock()
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
			ctx = context.WithValue(ctx, ContextKeyGinContext, c)
			ctx = context.WithValue(ctx, ContextKeyHttpRequest, req)
			ctx = context.WithValue(ctx, ContextKeyHttpRequest, req)
			ctx = context.WithValue(ctx, ContextKeyRequest, request)
			ctx = context.WithValue(ctx, ContextKeyHttpResponse, resp)
			ctx = context.WithValue(ctx, ContextKeyResponse, response)
			ctx = context.WithValue(ctx, ContextKeyTimeResponse, timeResponse)
			ctx = context.WithValue(ctx, ContextKeyAttempt, attempt)
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

	timeResponse = time.Since(beginTime).Milliseconds()
	//timeResponse = time.Now().Sub(beginTime).Milliseconds()
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

func (h *httpClient) doWithRetries(c *gin.Context, method, url string, headers map[string][]string, request interface{}, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error) {
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
		resp, err = h.do(c, method, url, headers, request, encode, decode, attempt)

		if err == nil && resp != nil {
			break Loop
		}

		if err != nil {
			switch merr := err.(type) {
			case *URL.Error:
				_ = merr
				err2 := err.(*URL.Error)
				if err2.Timeout() {
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

func (h *httpClient) GET(c *gin.Context, url string, headers map[string][]string, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(c, http.MethodGet, url, headers, nil, encode, decode)
}

func (h *httpClient) POST(c *gin.Context, url string, headers map[string][]string, request interface{}, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(c, http.MethodPost, url, headers, request, encode, decode)
}

func (h *httpClient) PUT(c *gin.Context, url string, headers map[string][]string, request interface{}, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(c, http.MethodPut, url, headers, request, encode, decode)
}

func (h *httpClient) PATCH(c *gin.Context, url string, headers map[string][]string, request interface{}, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(c, http.MethodPatch, url, headers, request, encode, decode)
}

func (h *httpClient) DELETE(c *gin.Context, url string, headers map[string][]string, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(c, http.MethodDelete, url, headers, nil, encode, decode)
}

func (h *httpClient) HEAD(c *gin.Context, url string, headers map[string][]string, encode EncodeRequestFunc, decode DecodeResponseFunc) (interface{}, error) {
	return h.doWithRetries(c, http.MethodHead, url, headers, nil, encode, decode)
}
