package service

import (
	"auth/pkg/logging"
)

type AuthService struct {
	logger *logging.Logger
	//storage  store.UserStorage
}

func NewAuthService(logger *logging.Logger) *AuthService {
	return &AuthService{
		logger: logger ,
		//storage: storage,
	}
}