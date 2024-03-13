package repository_models

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
)

type RedisUserRepository struct {
	RedisClient *redis.Client
}

func NewRedisUserRepository(client *redis.Client) *RedisUserRepository {
	// This is the constructor you were looking for
	return &RedisUserRepository{
		RedisClient: client,
	}
}

func (r *RedisUserRepository) FindUser(ctx context.Context, username string) (*User, error) {
	userBytes, err := r.RedisClient.Get(ctx, username).Bytes()
	if err != nil {
		return nil, err
	}

	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *RedisUserRepository) Login(ctx context.Context, username, password string) (string, error) {
	_, err := r.FindUser(ctx, username)
	if err != nil {
		return "", err
	}

	if !r.IsValidPassword(ctx, username, password) {
		return "", errors.New("invalid password")
	}

	// Generar el token JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	// Guardar el token en Redis
	err = r.RedisClient.Set(ctx, tokenString, "valid", time.Hour*24).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (r *RedisUserRepository) IsValidPassword(ctx context.Context, username, password string) bool {
	userBytes, err := r.RedisClient.Get(ctx, username).Bytes()
	if err != nil {
		return false
	}

	var user User
	if err := json.Unmarshal(userBytes, &user); err != nil {
		return false
	}

	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) == nil
}
