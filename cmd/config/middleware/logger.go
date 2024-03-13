package middlewares

import (
	"bytes"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/gin-gonic/gin"

	"autenticacion-ms/cmd/config/errors"
	"autenticacion-ms/cmd/config/masker"
	"autenticacion-ms/cmd/config/model"
	"autenticacion-ms/cmd/entity"
	"autenticacion-ms/cmd/logging"
	"autenticacion-ms/cmd/utils"
)

var (
	LoggingInternalCallSecondaries = os.Getenv("ENABLE_INTERNAL_CALL_LOGGING")
	API_LAYER                      = "API"
)

func Logger(log logging.Logger, artifact *model.ArtifactResources, loggingPayload bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		if !strings.Contains(c.Request.URL.Path, "/healthCheck") {
			c = logging.AddOperationInContext(c, artifact)
			logging.WithRequest(c, c.Request)
			if c.Request.Header.Get(utils.OriginRequested) != utils.InternalCall || LoggingInternalCallSecondaries == "true" {

				blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
				c.Writer = blw
				defer func() {
					latency := time.Since(t)
					var (
						requestMasked    interface{}
						responseMasked   interface{}
						dataBodyOriginal []uint8
						entityResponse   entity.Response
						err2             error
					)
					if loggingPayload {

						if requestOriginalBody, exist := c.Get(utils.RequestOriginalBody); exist {
							dataBodyOriginal = requestOriginalBody.([]uint8)
						}

						requestBody, _ := c.Get(utils.RequestBody)

						if requestBody != nil {
							requestMasked, err2 = masker.Marshal(requestBody)
							if err2 != nil {
								log.Error(err2.Error(), logging.AnyField("layer", API_LAYER))
							}
						}

						entityResponseCtx, _ := c.Get("entityResponse")

						if entityResponseCtx != nil {
							if reflect.TypeOf(entityResponseCtx) == reflect.TypeOf(errors.ErrorResponse{}) {
								errorResponse := entityResponseCtx.(errors.ErrorResponse)
								entityResponse.Data = errorResponse.ErrorBody.Data
								entityResponse.Result = entity.Result{Source: errorResponse.ErrorBody.Result.Source, Details: errorResponse.ErrorBody.Result.Details}
							} else {
								entityResponse.Data = entityResponseCtx.(entity.Response).Data
								entityResponse.Result = entityResponseCtx.(entity.Response).Result
							}

						}

						responseMasked, err2 = masker.Marshal(entityResponse)
						if err2 != nil {
							log.With(c).Error(err2.Error(), logging.AnyField("layer", API_LAYER))
						}
					}

					request := entity.Logger{Headers: masker.MaskerHeadersV2(c.Request.Header, utils.HeadersToMaskerIntegration...), Method: c.Request.Method, Path: c.FullPath(), UriPath: c.Request.URL.String(), Body: string(dataBodyOriginal), BodyFormatted: requestMasked}
					if val, _ := c.Get(utils.IdempotencyKey); val != nil {
						idempotencyKeyPrimary := val.(string)
						request.Headers.Set(utils.IdempotencyKeyHeader, idempotencyKeyPrimary)
					}

					response := entity.Logger{Headers: masker.MaskerHeadersV2(c.Writer.Header(), utils.HeadersToMaskerIntegration...), Method: c.Request.Method, Path: c.FullPath(), UriPath: c.Request.URL.String(), Body: responseMasked}
					log.With(c).Info("",
						logging.AnyField("layer", API_LAYER),
						logging.AnyField("timeResponse", latency.Milliseconds()),
						logging.AnyField("httpStatus", c.Writer.Status()),
						logging.AnyField("request", request),
						logging.AnyField("response", response),
					)
				}()

			}
		}
		c.Next()
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
