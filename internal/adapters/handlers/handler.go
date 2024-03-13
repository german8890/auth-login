package handlers

import "autenticacion-ms/internal/core/services"

type AuthHttp struct {
	service services.AppAuthService
}

func MakeNewAuthHttp(service services.AppAuthService) AuthHttp {
	return AuthHttp{
		service: service,
	}
}
