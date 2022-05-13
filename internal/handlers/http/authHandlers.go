package http

import (
	"auth/domain"

	handler "auth/internal/handlers"
	"auth/internal/service"
	"auth/pkg/jwt"
	"auth/pkg/logging"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	loginURL = "/login"
)

type authHandler struct {
	Service service.AuthService
	logger  *logging.Logger
}

func NewAuthHandler(service service.AuthService, logger *logging.Logger) handler.Handler {
	return &authHandler{
		Service: service,
		logger:  logger,
	}
}

func (ah *authHandler) Register(router *httprouter.Router) {
	router.GET(loginURL, ah.Login)
}

func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	login := r.Context().Value(domain.LoginValidateCtXKey).(*domain.Login)

	tk, err := jwt.CreateToken(login.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		//json.NewEncoder(w).Encode(apperror.NewAppError("Eror during createing tokens", fmt.Sprintf("%d", http.StatusUnprocessableEntity), err.Error()))
		return
	}

	c := http.Cookie{
		Name:     "Refresh-Token",
		Value:    tk.RefreshToken,
		HttpOnly: true,
	}

	http.SetCookie(w, &c)

	w.Header().Set("Access-Token", tk.AccessToken)
	w.WriteHeader(http.StatusOK)
	//json.NewEncoder(w).Encode(user)
}
