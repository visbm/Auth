package service

import (
	"auth/internal/store"
	"auth/pkg/logging"
)

type AuthService struct {
	logger *logging.Logger
	storage  store.UserStorage
}

func NewAuthService(storage store.UserStorage, logger *logging.Logger) *AuthService {
	return &AuthService{
		logger: logger ,
		storage: storage,
	}
}