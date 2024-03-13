package errors

import (
	"encoding/json"
	"fmt"
	"net/url"

	"autenticacion-ms/cmd/entity"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
)

func CustomRecoveryGinPanic() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, err any) {
		_ = c.Error(err.(error))
	})
}

func Handler404() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		details := []entity.Detail{{Message: "Resource Not Found"}}
		ErrorWrapper(ctx, NotFound(details, ""))
	}
}

func HandleDecodeRequest(requestBody interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		_ = ctx.ShouldBindJSON(&requestBody)
		ctx.Set("requestBody", requestBody)
		ctx.Next()
	}
}

func HandleEncodeResponse(response interface{}) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()
	}
}

func HandlePanic() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if e := recover(); e != nil {
				var (
					ok  bool
					err error
				)
				if err, ok = e.(error); !ok {
					err = fmt.Errorf("%v", e)
				}
				ErrorWrapper(ctx, err)
			}
		}()
		ctx.Next()
	}
}

// TODO: Falta agregar logs
// Handler creates a middleware that handles panics and errors encountered during HTTP request processing.
// func Handler(logger log.Logger) routing.Handler {
func Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			for _, err := range ctx.Errors {
				ErrorWrapper(ctx, err)
			}
		}()
		ctx.Next()
	}
}

func ErrorWrapper(c *gin.Context, err error) {
	res := buildErrorResponse(err)
	c.Set("entityResponse", res.Body())
	c.Header("Content-Type", "application/json")
	c.AbortWithStatusJSON(res.StatusCode(), res.Body())
	c.Done()
}

// buildErrorResponse builds an error response from an error.
func buildErrorResponse(err error) ErrorResponse {
	switch merr := err.(type) {
	case ErrorResponse:
		return err.(ErrorResponse)
	case *gin.Error:
		err2 := err.(*gin.Error)
		var errGin *gin.Error
		errGin = err2
		switch errGin.Err.(type) {
		case ErrorResponse:
			return errGin.Err.(ErrorResponse)
		case *json.UnmarshalTypeError:
			validationErrors := errGin.Err.(*json.UnmarshalTypeError)
			return InvalidJsonInput(*validationErrors)
		case validator.ValidationErrors:
			validationErrors := errGin.Err.(validator.ValidationErrors)
			return InvalidInput(validationErrors)
		case *pgconn.PgError:
			pgErr := errGin.Err.(*pgconn.PgError)
			return PgError(pgErr)
		case *url.Error:
			err3 := errGin.Err.(*url.Error)
			switch err3.Err.(type) {
			case ErrorResponse:
				return err3.Err.(ErrorResponse)
			}
			return UrlErrorHandle(err3)
		}
		_ = merr
	}

	details := []entity.Detail{{Message: err.Error()}}
	return InternalServerError(details, "")
}
