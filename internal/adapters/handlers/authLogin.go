package handlers

import (
	"context"
	"net/http"

	"autenticacion-ms/cmd/utils"

	"github.com/gin-gonic/gin"
)

func (h *AuthHttp) Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
		var request interface{}
		if err := utils.ShouldBindJSON(c, &request); err != nil {
			return
		}
		response, err := h.service.Login(context.Background(), username, password)
		if err != nil {
			_ = c.Error(err)
			return
		}
		c.JSON(http.StatusOK, response)
	}
}
