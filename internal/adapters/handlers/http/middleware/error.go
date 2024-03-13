package middlewares

import (
	"autenticacion-ms/cmd/config/errors"

	"github.com/gin-gonic/gin"
)

func ErrorHandler() gin.HandlerFunc {
	return errors.Handler()
}
func HandlerPanic() gin.HandlerFunc {
	return errors.HandlePanic()
}
func Handler404() gin.HandlerFunc {
	return errors.Handler404()
}
