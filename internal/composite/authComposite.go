package composite

import (
	handler "auth/internal/handlers"
	"auth/internal/handlers/http"
	"auth/internal/service"
	"auth/internal/store"

	"auth/pkg/logging"
)

type AuthComposite struct {
	logger  *logging.Logger
	Storage store.UserStorage
	Service *service.AuthService
	Handler handler.Handler
}

func (ac *AuthComposite) New(logger *logging.Logger, storage store.UserStorage) {
	ac.logger = logger
	ac.Storage = storage
	ac.Service = service.NewAuthService(ac.logger , storage)
	ac.Handler = http.NewAuthHandler(*ac.Service, ac.logger)
}
