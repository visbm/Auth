package composite

import (
	handler "auth/internal/handlers"
	"auth/internal/handlers/http"
	"auth/internal/service"

	"auth/pkg/logging"
)

type AuthComposite struct {
	logger *logging.Logger
	//Storage store.UserStorage
	Service *service.AuthService
	Handler handler.Handler
}

func (ac *AuthComposite) New(logger *logging.Logger) {
	ac.logger = logger
	//	ac.Storage = storage
	ac.Service = service.NewAuthService(ac.logger)
	ac.Handler = http.NewAuthHandler(*ac.Service , ac.logger)
}
