package services

import (
	"context"

	repository "autenticacion-ms/internal/core/domain/repository"
)

// AppAuthService encapsula la lógica de autenticación.
type AppAuthService struct {
	UserRepository repository.RedisUserRepository
}

func NewAuthService(userRepository repository.RedisUserRepository) *AppAuthService {
	return &AppAuthService{
		UserRepository: userRepository,
	}
}

// Login autentifica a un usuario y genera un token JWT.
func (s *AppAuthService) Login(ctx context.Context, request *repository.User) (string, error) {
	username := request.Username
	password := request.Password
	token, err := s.UserRepository.Login(ctx, username, password)
	if err != nil {
		return "", err
	}

	return token, nil
}
