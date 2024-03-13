package middlewares

import (
	"context"
	"fmt"
	"net/http"

	"autenticacion-ms/cmd/config/http_client"
	http_client2 "autenticacion-ms/cmd/config/http_client/v2"

	"github.com/gin-gonic/gin"

	"autenticacion-ms/cmd/config/errors"
	"autenticacion-ms/cmd/config/masker"
	"autenticacion-ms/cmd/entity"
	"autenticacion-ms/cmd/logging"
	"autenticacion-ms/cmd/utils"
)

type key int

const (
	ContextKeyGenericError key = iota
)

func LoggerInterceptorClientV2(log logging.Logger, nameBackend string, loggingPayload bool) http_client.ClientFinalizerFunc {
	return func(ctx context.Context, err error) {
		var (
			resp         *http.Response
			_            *context.Context
			req          *http.Request
			request      interface{}
			response     interface{}
			timeResponse interface{}
			requestLog   = entity.Logger{Headers: nil, Body: new(interface{})}
			responseLog  = entity.Logger{Headers: nil, Body: new(interface{})}
			httpStatus   int
			attempt      interface{}
			path         string
		)

		if v := ctx.Value(http_client2.ContextKeyTimeResponse); v != nil {
			timeResponse = v
		}

		if val := ctx.Value(http_client2.ContextKeyAttempt); val != nil {
			attempt = val
		}

		if val := ctx.Value(http_client2.ContextKeyPathRequest); val != nil {
			path = val.(string)
		}

		if v := ctx.Value(http_client2.ContextKeyHttpResponse); v != nil {
			resp = v.(*http.Response)
			if resp != nil {
				httpStatus = resp.StatusCode
				responseLog.Headers = masker.MaskerHeaders(resp.Header, utils.HeadersToMaskerIntegration...)
				responseLog.Method = resp.Request.Method
				responseLog.Path = path
				responseLog.UriPath = resp.Request.URL.String()

			}
		}

		if v := ctx.Value(http_client2.ContextKeyHttpRequest); v != nil {
			req = v.(*http.Request)
			requestLog.Headers = masker.MaskerHeaders(req.Header, utils.HeadersToMaskerIntegration...)
			requestLog.Method = req.Method
			requestLog.Path = path
			requestLog.UriPath = req.URL.String()
		}

		if loggingPayload {
			if v := ctx.Value(http_client2.ContextKeyRequest); v != nil {
				request = v
				if masked, error := masker.Marshal(request); error == nil {
					requestLog.Body = masked
				} else {
					log.With(ctx).Error(error.Error(), logging.AnyField("layer", nameBackend))
				}
			}

			if v := ctx.Value(http_client2.ContextKeyResponse); v != nil {
				response = v
				if masked, error := masker.Marshal(response); error == nil {
					responseLog.Body = masked
				} else {
					log.With(ctx).Error(error.Error(), logging.AnyField("layer", nameBackend))
				}
			}

			if err != nil {
				if resp != nil && resp.Request.Context() != nil && resp.Request.Context().Value(ContextKeyGenericError) != nil {
					genericError := resp.Request.Context().Value(ContextKeyGenericError)
					responseLog.Body = genericError
				} else {
					switch errorT := err.(type) {
					case errors.ErrorResponse:
						responseLog.Body = errorT.ErrorBody
					default:
						httpStatus = http.StatusInternalServerError
						responseLog.Body = &entity.Error{Error: fmt.Sprintf("%T", err), Detail: err.Error()}
					}
				}
			}
		}
		if val1 := req.Context().Value(1); val1 != nil {
			val2 := val1.(context.Context)
			ctx = val2
			if val3 := val2.Value(gin.ContextKey); val3 != nil {
				c2 := val3.(*gin.Context)
				ctx = c2.Copy()
			}

		}
		//ctx = logging.WithRequestV2(ctx, req)
		log.With(ctx).Info("",
			logging.AnyField("layer", nameBackend),
			logging.AnyField("timeResponse", timeResponse),
			logging.AnyField("attempt", attempt),
			logging.AnyField("httpStatus", httpStatus),
			logging.AnyField("request", requestLog),
			logging.AnyField("response", responseLog),
		)
	}
}

func LoggerInterceptorClientV3(log logging.Logger, nameBackend string, loggingPayload bool) http_client.ClientFinalizerFunc {
	return func(ctx context.Context, err error) {
		var (
			resp         *http.Response
			_            *context.Context
			req          *http.Request
			c            *gin.Context
			request      interface{}
			response     interface{}
			timeResponse interface{}
			requestLog   = entity.Logger{Headers: nil, Body: new(interface{})}
			responseLog  = entity.Logger{Headers: nil, Body: new(interface{})}
			httpStatus   int
			attempt      interface{}
			path         string
		)

		if v := ctx.Value(http_client.ContextKeyGinContext); v != nil {
			cc, _ := v.(*gin.Context)
			c = cc
		}

		if v := ctx.Value(http_client.ContextKeyTimeResponse); v != nil {
			timeResponse = v
		}

		if v := ctx.Value(http_client.ContextKeyAttempt); v != nil {
			attempt = v
		}

		if val, exist := c.Get(http_client.ContextKeyPathRequest); exist {
			path = val.(string)
		}

		if v := ctx.Value(http_client.ContextKeyHttpResponse); v != nil {
			resp = v.(*http.Response)
			if resp != nil {
				httpStatus = resp.StatusCode
				responseLog.Headers = masker.MaskerHeaders(resp.Header, utils.HeadersToMaskerIntegration...)
				responseLog.Method = resp.Request.Method
				responseLog.Path = path
				responseLog.UriPath = resp.Request.URL.String()

			}
		}

		if v := ctx.Value(http_client.ContextKeyHttpRequest); v != nil {
			req = v.(*http.Request)
			requestLog.Headers = masker.MaskerHeaders(req.Header, utils.HeadersToMaskerIntegration...)
			requestLog.Method = req.Method
			requestLog.Path = path
			requestLog.UriPath = req.URL.String()
		}

		if loggingPayload {
			if v := ctx.Value(http_client.ContextKeyRequest); v != nil {
				request = v
				if masked, error := masker.Marshal(request); error == nil {
					requestLog.Body = masked
				} else {
					log.With(c).Error(error.Error(), logging.AnyField("layer", nameBackend))
				}
			}

			if v := ctx.Value(http_client.ContextKeyResponse); v != nil {
				response = v
				if masked, error := masker.Marshal(response); error == nil {
					responseLog.Body = masked
				} else {
					log.With(c).Error(error.Error(), logging.AnyField("layer", nameBackend))
				}
			}

			if err != nil {
				if resp != nil && resp.Request.Context() != nil && resp.Request.Context().Value(ContextKeyGenericError) != nil {
					genericError := resp.Request.Context().Value(ContextKeyGenericError)
					responseLog.Body = genericError
				} else {
					switch errorT := err.(type) {
					case errors.ErrorResponse:
						responseLog.Body = errorT.ErrorBody
					default:
						httpStatus = http.StatusInternalServerError
						responseLog.Body = &entity.Error{Error: fmt.Sprintf("%T", err), Detail: err.Error()}
					}
				}
			}
		}

		log.With(c).Info("",
			logging.AnyField("layer", nameBackend),
			logging.AnyField("timeResponse", timeResponse),
			logging.AnyField("attempt", attempt),
			logging.AnyField("httpStatus", httpStatus),
			logging.AnyField("request", requestLog),
			logging.AnyField("response", responseLog),
		)
	}
}
