package service

import (
	"auth/domain"
	"auth/internal/store"
	"auth/pkg/logging"
	"context"
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

func (s *AuthService) GetByUUID(ctx context.Context,UUID string) (*domain.Login, error) {
	return s.storage.GetByUUID(UUID)
}

func (s *AuthService) GetByUsername(ctx context.Context,username string) (*domain.Login, error) {
	return s.storage.GetByUsername(username)
}

func (s *AuthService) GetAll(ctx context.Context, limit, offset int) ([]*domain.Login, error) {
	return s.storage.GetAll(limit, offset)
}

func (s *AuthService) Delete(ctx context.Context, UUID string) error {
	return s.storage.Delete(UUID)
}

func (s *AuthService) Create(ctx context.Context, Login *domain.Login) (string, error) {
	return s.storage.Create(Login)
}

func (s *AuthService) Update(ctx context.Context, Login *domain.Login) error {
	return s.storage.Update(Login)
}
