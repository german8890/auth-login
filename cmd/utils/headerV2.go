package utils

import (
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"strings"

	"google.golang.org/grpc/metadata"
)

type Headers struct {
	headers []string
}

func MakeNewHeadersToCopy(extras ...string) *Headers {
	HeadersToCopy := []string{
		CountryHeader,
		ConsumerRefHeader,
		ChannelRefHeader,
		BrandHeader,
		StoreRefHeader,
		UserTrxHeader,
		EnvHeader,
		ConsumerDateTimeHeader,
		ProcessRefHeader,
		TypeProductHeader,
		TypeProcessRefHeader,
		IdempotencyKeyHeader,
		B3TraceIdHeader,
		B3SpanIdHeader,
		Apikey,
		AcceptLanguage,
	}
	HeadersToCopy = append(HeadersToCopy, extras...)
	return &Headers{headers: HeadersToCopy}
}

func (h Headers) GetHeadersInMap(r *http.Request) map[string][]string {
	var headers = make(map[string][]string)
	if span := trace.SpanFromContext(r.Context()); span != nil {
		span := trace.SpanFromContext(r.Context())
		if r.Header.Get(B3TraceIdHeader) != "" {
			r.Header.Del(B3TraceIdHeader)
			r.Header.Add(B3TraceIdHeader, span.SpanContext().TraceID().String())
		}
		if r.Header.Get(B3SpanIdHeader) != "" {
			r.Header.Del(B3SpanIdHeader)
			r.Header.Add(B3SpanIdHeader, span.SpanContext().SpanID().String())
		}
	}

	for _, v := range h.headers {
		val := r.Header.Get(v)
		if val != "" {
			headers[v] = []string{val}
		}
		if val == "" {
			requestHeaders := r.Header
			for keyCode, keyValue := range requestHeaders {
				keyCodeSensitive := strings.ToLower(keyCode)
				headSensitive := strings.ToLower(v)
				if keyCodeSensitive == headSensitive {
					if len(keyValue) > 0 {
						headers[v] = []string{keyValue[0]}
					}
				}
			}
		}
	}
	return headers
}

func (h Headers) HeadersFilterToMap(mapHeaders map[string][]string) map[string][]string {
	var head = new(http.Request)
	head.Header = mapHeaders
	return h.GetHeadersInMap(head)
}

func (h Headers) GetMetadataInMap(md metadata.MD) *map[string][]string {
	var apigeeHeaders = make(map[string][]string)
	for _, header := range h.headers {
		headerVal := md.Get(header)
		if len(headerVal) > 0 {
			apigeeHeaders[header] = []string{headerVal[0]}
		}
	}
	return &apigeeHeaders
}
