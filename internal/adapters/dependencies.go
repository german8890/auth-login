package adapters

import (
	"autenticacion-ms/cmd/config"
	auth "autenticacion-ms/internal/adapters/repository/auth"
	repository "autenticacion-ms/internal/core/domain/repository"
	"autenticacion-ms/internal/core/services"
	"fmt"
	"net/http"

	"autenticacion-ms/cmd/logging"
)

type Dependencies struct {
	AuthService *services.AppAuthService
}

func InitDependencies(cfg *config.Config, logger logging.Logger, httpClient *http.Client) *Dependencies {
	redisClient, err := auth.NewRedisClient(cfg.Redis.RedisAddr, cfg.Redis.RedisPassword)
	if err != nil {
		logger.Error("Error al crear el cliente Redis")
		fmt.Println(err)
		// Aquí puedes agregar cualquier otro manejo de errores necesario
	} else {
		logger.Info("Cliente Redis creado exitosamente")
	}
	userRepository := repository.NewRedisUserRepository(redisClient)
	authService := services.NewAuthService(*userRepository)

	return &Dependencies{
		AuthService: authService, // Only include the AuthService dependency
	}
}
