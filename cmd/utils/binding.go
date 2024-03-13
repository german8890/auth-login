package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"autenticacion-ms/cmd/config/errors"
	"autenticacion-ms/cmd/entity"

	"github.com/gin-gonic/gin"
)

const (
	RequestOriginalBody = "requestOriginalBody"
	RequestBody         = "requestBody"
)

func ShouldBindJSON(c *gin.Context, obj any) error {
	body, _ := io.ReadAll(c.Request.Body)
	c.Set(RequestOriginalBody, body)
	c.Request.Body = io.NopCloser(bytes.NewReader(body))
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	err := ValidateJSON(body)
	if err != nil {
		_ = c.Error(err)
		return err
	}
	if err := c.ShouldBindJSON(&obj); err != nil {
		_ = c.Error(err)
		var jsonText interface{}
		if err2 := json.Unmarshal(body, &jsonText); err2 != nil {
			c.Set(RequestBody, string(body))

		} else {
			c.Set(RequestBody, jsonText)

		}
		return err
	}

	c.Set(RequestBody, obj)
	return nil
}

func ShouldBindQuery(c *gin.Context, obj interface{}) error {
	c.Set(RequestOriginalBody, []byte(c.Request.URL.RawQuery))
	queryParams := c.Request.URL.RawQuery
	if err := c.ShouldBindQuery(obj); err != nil {
		c.Set(RequestBody, queryParams)
		return err
	}
	c.Set(RequestBody, obj)
	return nil
}

func ValidateJSON(body []byte) error {
	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.DisallowUnknownFields()

	var parsedJSON interface{}
	if err := decoder.Decode(&parsedJSON); err != nil {
		fieldName := findErrorField(err, body)
		var errDetails []entity.Detail
		errDetailsSingle := entity.Detail{
			Message:      "El JSON en su solicitud es inválido o mal formado.",
			InternalCode: "400",
			Detail:       fmt.Sprintf("Error al validar el JSON: [%s]", err.Error()),
		}
		if fieldName != "" {
			errDetailsSingle.Detail = fmt.Sprintf("Error al validar el JSON: [%s]. Error en campo o cerca del campo: [%s]", err.Error(), fieldName)
		}
		errDetails = append(errDetails, errDetailsSingle)
		errResponse := errors.BadRequest(errDetails, errors.CONSUMER)
		return errResponse
	}

	return nil
}

func findErrorField(err error, jsonInput []byte) string {
	if se, ok := err.(*json.SyntaxError); ok {
		return extractField(se, jsonInput)
	}
	return ""
}

func extractField(se *json.SyntaxError, jsonInput []byte) string {
	errorPosition := int(se.Offset)
	jsonString := string(jsonInput)

	// Find the start position of the field name by searching for the nearest '"' before the error position
	quoteStart := strings.LastIndex(jsonString[:errorPosition], `"`)
	if quoteStart >= 0 {
		// Find the nearest ':' after the quoteStart
		colonStart := strings.Index(jsonString[quoteStart:], ":") + quoteStart
		if colonStart >= 0 {
			// Extract the content between the nearest '"' and the ':' to get the field name
			if quoteStart+1 < colonStart {
				fieldContent := jsonString[quoteStart+1 : colonStart]
				// Trim whitespace and quotes to get the field name
				fieldName := strings.TrimSpace(fieldContent)
				return fieldName
			}
		}
	}
	return ""
}
