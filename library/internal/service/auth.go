package service

import (
	"auth/domain"
	"auth/internal/store"
	"auth/pkg/logging"
)

type AuthService struct {
	logger  *logging.Logger
	storage store.UserStorage
}

func NewAuthService(logger *logging.Logger, storage store.UserStorage) *AuthService {
	return &AuthService{
		logger:  logger,
		storage: storage,
	}
}

func (s *AuthService) GetByUUID(UUID string) (*domain.Login, error) {
	return s.storage.GetByUUID(UUID)
}

func (s *AuthService) GetByUsername(username string) (*domain.Login, error) {
	return s.storage.GetByUsername(username)
}

func (s *AuthService) GetAll(limit, offset int) ([]*domain.Login, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *AuthService) Delete(UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *AuthService) Create(Login *domain.Login) (string, error) {
	return s.storage.Create(Login)
}

func (s *AuthService) Update(Login *domain.Login) error {
	return s.storage.Update(Login)
}
