package errors

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"autenticacion-ms/cmd/entity"

	"github.com/go-playground/validator/v10"
	"github.com/jackc/pgconn"
)

var (
	CONSUMER         = "CONSUMER"
	INTERNALDATABASE = "INTERNAL DATABASE"
	SYSTEM           = "SYSTEM"
)

type ErrorResponse struct {
	Status    int
	ErrorBody entity.Response
}

func (e ErrorResponse) Body() entity.Response {
	return e.ErrorBody
}

// Error is required by the error interface.
func (e ErrorResponse) Error() string {
	data, _ := json.Marshal(e.Body())
	return fmt.Sprint(string(data))
}

// StatusCode is required by routing.HTTPError interface.
func (e ErrorResponse) StatusCode() int {
	return e.Status
}

// InternalServerError creates a new error response representing an integral server error (HTTP 500).
func InternalServerError(details []entity.Detail, source string) ErrorResponse {
	if details == nil {
		details = make([]entity.Detail, 1)
		details[0].Message = "We encountered an error while processing your request."
	}
	if source == "" {
		source = SYSTEM
	}
	return fillErrorResponse(http.StatusInternalServerError, details, source)
}

// NotFound creates a new error response representing a resource-not-found error (HTTP 400).
func NotFound(details []entity.Detail, source string) ErrorResponse {
	if details == nil {
		details = make([]entity.Detail, 1)
		details[0].Message = "The requested resource was not found."
	}
	return fillErrorResponse(http.StatusNotFound, details, source)
}

// Unauthorized creates a new error response representing an authentication/authorization failure (HTTP 401)
func Unauthorized(details []entity.Detail, source string) ErrorResponse {
	if details == nil {
		details = make([]entity.Detail, 1)
		details[0].Message = "You are not authenticated to perform the requested action."
	}

	return fillErrorResponse(http.StatusUnauthorized, details, source)
}

// Forbidden creates a new error response representing an authorization failure (HTTP 403)
func Forbidden(details []entity.Detail, source string) ErrorResponse {
	if details == nil {
		details = make([]entity.Detail, 1)
		details[0].Message = "You are not authorized to perform the requested action."
	}
	return fillErrorResponse(http.StatusForbidden, details, source)
}

// BadRequest creates a new error response representing a bad request (HTTP 400)
func BadRequest(details []entity.Detail, source string) ErrorResponse {
	if details == nil {
		details = make([]entity.Detail, 1)
		details[0].Message = "Your request is in a bad format."
	}
	return fillErrorResponse(http.StatusBadRequest, details, source)
}

func ClientError(httpStatus int, details []entity.Detail, source string) ErrorResponse {
	if details == nil {
		details = make([]entity.Detail, 1)
		details[0].Message = "Client Error"
	}
	return fillErrorResponse(httpStatus, details, source)
}

// InvalidInput creates a new error response representing a data validation error (HTTP 400).
func InvalidInput(errs validator.ValidationErrors) ErrorResponse {
	var details []entity.Detail
	for _, field := range errs {
		details = append(details, entity.Detail{Message: field.Error()})
	}
	return fillErrorResponse(http.StatusBadRequest, details, CONSUMER)
}

// InvalidJsonInput creates a new error response representing a data validation error (HTTP 400).
func InvalidJsonInput(errs json.UnmarshalTypeError) ErrorResponse {
	var details []entity.Detail
	details = append(details, entity.Detail{Message: errs.Error()})
	return fillErrorResponse(http.StatusBadRequest, details, CONSUMER)
}

// PgError creates a new error response
func PgError(err *pgconn.PgError) ErrorResponse {
	details := make([]entity.Detail, 1)
	details[0].InternalCode = err.Code
	details[0].Message = err.Detail
	return fillErrorResponse(http.StatusInternalServerError, details, INTERNALDATABASE)
}

func UrlErrorHandle(err *url.Error) ErrorResponse {
	details := make([]entity.Detail, 1)
	details[0].InternalCode = strconv.Itoa(http.StatusInternalServerError)
	details[0].Message = err.Unwrap().Error()
	details[0].Detail = err.Error()
	return fillErrorResponse(http.StatusInternalServerError, details, SYSTEM)
}

func fillErrorResponse(httpStatus int, details []entity.Detail, source string) ErrorResponse {
	if source == "" {
		source = CONSUMER
	}

	for i := 0; i < len(details); i++ {
		if details[i].InternalCode == "" {
			details[i].InternalCode = strconv.Itoa(httpStatus)
		}
	}

	return ErrorResponse{
		Status: httpStatus,
		ErrorBody: entity.Response{
			Result: entity.Result{
				Details: details,
				Source:  source,
			},
		},
	}
}
