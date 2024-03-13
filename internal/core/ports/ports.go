package ports

import (
	"context"

	repository "autenticacion-ms/internal/core/domain/repository"

	"github.com/gin-gonic/gin"
)

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	FindUser(username string) (*repository.User, error)
}
type AuthHandler interface {
	Login(c *gin.Context)
}

type UserRepository interface {
	FindUser(ctx context.Context, username string) (*repository.User, error)
}
