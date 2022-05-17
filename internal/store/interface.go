package store

import "auth/domain"

type UserStorage interface {
	GetByUsername(username string) (*domain.Login, error)
	GetByUUID(UUID string) (*domain.Login, error)
	GetAll(limit, offset int) ([]*domain.Login, error)
	Create(login *domain.Login) (string, error)
	Delete(UUID string) error
	Update(login *domain.Login) error
}
