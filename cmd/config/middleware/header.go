package middlewares

import (
	"autenticacion-ms/cmd/utils"
	"fmt"
	"strings"

	"autenticacion-ms/cmd/config/errors"
	"autenticacion-ms/cmd/entity"

	"github.com/gin-gonic/gin"
)

type HeaderRequired struct {
	Name     string
	Message  string
	Required bool
}

type EndpointsException struct {
	Endpoint string
}

var HeadersRequiredDefault = []HeaderRequired{
	{Name: utils.CountryHeader, Message: fmt.Sprintf("Required parameter '%v' is not present in request header", utils.CountryHeader), Required: true},
	//{Name: utils.AcceptLanguage, Message: fmt.Sprintf("Required parameter '%v' is not present in request header", utils.AcceptLanguage), Required: true},
}

var EndpointsExceptionsDefault = []EndpointsException{
	{Endpoint: "/healthCheck"},
	{Endpoint: "/engine-rest"},
}

func HeaderValidation(headers ...HeaderRequired) gin.HandlerFunc {
	var headersRequired []HeaderRequired
	headersRequired = append(headersRequired, HeadersRequiredDefault...)
	if len(headers) > 0 {
		headersRequired = append(headersRequired, headers...)
	}
	return func(ctx *gin.Context) {
		//  if !strings.Contains(ctx.Request.URL.Path, "/healthCheck") || !strings.Contains(ctx.Request.URL.Path, "/engine-rest") {
		if !strings.Contains(ctx.Request.URL.Path, "/healthCheck") {
			var details []entity.Detail
			requestHeaders := ctx.Request.Header
			for _, head := range headersRequired {
				if head.Required && requestHeaders.Get(head.Name) == "" {
					details = append(details, entity.Detail{Message: head.Message})
				}
			}
			if len(details) > 0 {
				errors.ErrorWrapper(ctx, errors.BadRequest(details, ""))
			}
		}
		ctx.Next()
	}
}

func HeaderValidationV2(headers ...HeaderRequired) gin.HandlerFunc {
	var headersRequired []HeaderRequired
	headersRequired = append(headersRequired, HeadersRequiredDefault...)
	if len(headers) > 0 {
		headersRequired = append(headersRequired, headers...)
	}
	return func(ctx *gin.Context) {
		if !strings.Contains(ctx.Request.URL.Path, "/healthCheck") && !strings.Contains(ctx.Request.URL.Path, "/engine-rest") {
			var details []entity.Detail
			requestHeaders := ctx.Request.Header
			for _, head := range headersRequired {
				if head.Required && requestHeaders.Get(head.Name) == "" {
					details = append(details, entity.Detail{Message: head.Message})
				}
			}
			if len(details) > 0 {
				errors.ErrorWrapper(ctx, errors.BadRequest(details, ""))
			}
		}
		ctx.Next()
	}
}

func HeaderValidationV3(headers *[]HeaderRequired, exceptions *[]EndpointsException, context string) gin.HandlerFunc {
	var headersRequired []HeaderRequired
	var endpointsException []EndpointsException
	headersRequired = append(headersRequired, HeadersRequiredDefault...)
	if headers != nil && len(*headers) > 0 {
		headersRequired = append(headersRequired, *headers...)
	}
	endpointsException = append(endpointsException, EndpointsExceptionsDefault...)
	if exceptions != nil && len(*exceptions) > 0 {
		endpointsException = append(endpointsException, *exceptions...)
	}
	return func(ctx *gin.Context) {
		exists := false
		for _, endpoint := range endpointsException {
			if strings.Compare(ctx.FullPath(), context+endpoint.Endpoint) == 0 {
				exists = true
			}
		}
		if !exists {
			var details []entity.Detail
			requestHeaders := ctx.Request.Header
			for _, head := range headersRequired {
				if head.Required && requestHeaders.Get(head.Name) == "" {
					details = append(details, entity.Detail{Message: head.Message})
				}
			}
			if len(details) > 0 {
				errors.ErrorWrapper(ctx, errors.BadRequest(details, ""))
			}
		}
		ctx.Next()
	}
}
