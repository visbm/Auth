package composite

import (
	"auth/internal/service"
	"auth/internal/store"

	"auth/pkg/logging"
)

type AuthComposite struct {
	logger  *logging.Logger
	Storage store.UserStorage
	Service *service.AuthService
}

func (ac *AuthComposite) New(storage store.UserStorage, logger *logging.Logger) {
	ac.logger = logger
	ac.Storage = storage
	ac.Service = service.NewAuthService(ac.Storage, ac.logger)
}
