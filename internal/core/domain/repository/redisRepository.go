package repository_models

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"crypto/subtle"

	"github.com/dgrijalva/jwt-go"
	"github.com/redis/go-redis/v9"
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
	fmt.Println(username)
	userBytes, err := r.RedisClient.Get(ctx, username).Bytes()
	if err != nil {
		println("se quedo aqui !!! ojo")
		return nil, err
	}

	var user User
	err = json.Unmarshal(userBytes, &user)
	if err != nil {
		fmt.Println(err, userBytes, "vacio")
		return nil, err
	}
	fmt.Println("paso")

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
		"exp":      time.Now().Add(time.Minute * 30).Unix(),
	})
	tokenString, err := token.SignedString([]byte("s3cr3tK3yF0rJWT!"))
	if err != nil {
		return "", err
	}

	// Guardar el token en Redis
	err = r.RedisClient.Set(ctx, tokenString, "valid", time.Minute*30).Err()
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (r *RedisUserRepository) IsValidPassword(ctx context.Context, username, password string) bool {
	// Obtener los datos del usuario desde Redis
	userBytes, err := r.RedisClient.Get(ctx, username).Bytes()
	if err != nil {
		return false
	}
	fmt.Println("lego a valid")
	// Deserializar los datos del usuario
	var user User
	if err := json.Unmarshal(userBytes, &user); err != nil {
		return false
	}

	// Obtener el hash almacenado
	storedHash := []byte(user.Password)
	fmt.Println("validacion de store hash")
	// Comparar el hash almacenado con la contraseña proporcionada
	a := CompareHashes(storedHash, []byte(password))
	if a == true {
		return err == nil // Si err es nil, la contraseña es válida
	}
	return false
}

func CompareHashes(hash1, hash2 []byte) bool {
	// Verificar que los hash tengan la misma longitud
	if len(hash1) != len(hash2) {
		return false
	}
	fmt.Println("comparo hash")
	// Comparar byte a byte usando subtle.ConstantTimeCompare
	return subtle.ConstantTimeCompare(hash1, hash2) == 1
}
