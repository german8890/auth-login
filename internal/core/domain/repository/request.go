package repository_models

type User struct {
	Username string `redis:"username"`
	Password string `redis:"password"`
}
