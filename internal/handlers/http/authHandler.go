package http

import (
	"auth/domain"
	"encoding/json"
	"fmt"

	handler "auth/internal/handlers"
	"auth/internal/service"

	"auth/pkg/apperror"
	"auth/pkg/jwt"
	"auth/pkg/logging"
	"auth/pkg/response"

	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	loginURL        = "/login"
	logoutURL       = "/logout"
	getUserURL      = "/getUser"
	isLoggedInURL   = "/isloggedin"
	refreshTokenURL = "/refreshtoken"
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
	router.POST(loginURL, ah.Login)
	router.POST(logoutURL, ah.Logout)
	router.GET(getUserURL, ah.GetUser)
	router.GET(isLoggedInURL, ah.IsLoggedIn)
	router.POST(refreshTokenURL, ah.RefreshToken)

}

func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	login := &domain.Login{}
	if err := json.NewDecoder(r.Body).Decode(login); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Error(apperror.NewAppError("Error occured while createing tokens", fmt.Sprintf("%d", http.StatusInternalServerError), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("Error occured while createing tokens", fmt.Sprintf("%d", http.StatusUnprocessableEntity), err.Error()))
		return
	}
	tk, err := jwt.CreateToken(login.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		ah.logger.Error(apperror.NewAppError("Error occured while createing tokens", fmt.Sprintf("%d", http.StatusInternalServerError), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("Error occured while createing tokens", fmt.Sprintf("%d", http.StatusUnprocessableEntity), err.Error()))
		return
	}

	c := http.Cookie{
		Name:     "Refresh-Token",
		Value:    tk.RefreshToken,
		HttpOnly: true,
		MaxAge:   int(tk.RtExpires),
	}

	http.SetCookie(w, &c)

	w.Header().Set("Authorization", "Bearer "+tk.AccessToken)
	w.WriteHeader(http.StatusOK)
}

func (ah *authHandler) Logout(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	_, err := jwt.ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		ah.logger.Error(apperror.NewAppError("you are unauthorized", fmt.Sprintf("%d", http.StatusUnauthorized), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("you are unauthorized", fmt.Sprintf("%d", http.StatusUnauthorized), err.Error()))
		return
	}

	c := http.Cookie{
		Name:     "Refresh-Token",
		Value:    "",
		HttpOnly: true,
	}

	w.Header().Set("Authorization", "")
	http.SetCookie(w, &c)
	json.NewEncoder(w).Encode(response.Info{"Successfully logged out"})

}

func (ah *authHandler) GetUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	accessDetails, err := jwt.ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		ah.logger.Error(apperror.NewAppError("you are unauthorized", fmt.Sprintf("%d", http.StatusUnauthorized), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("you are unauthorized", fmt.Sprintf("%d", http.StatusUnauthorized), err.Error()))
		return
	}
	json.NewEncoder(w).Encode(accessDetails)

}

func (ah *authHandler) IsLoggedIn(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	_, err := jwt.ExtractTokenMetadata(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		ah.logger.Error(apperror.NewAppError("you are unauthorized", fmt.Sprintf("%d", http.StatusUnauthorized), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("you are unauthorized", fmt.Sprintf("%d", http.StatusUnauthorized), err.Error()))
		return
	}
	json.NewEncoder(w).Encode(response.Info{"You are already logged in"})

}

func (ah *authHandler) RefreshToken(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")

	newToken, err := jwt.RefreshToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		ah.logger.Error(apperror.NewAppError("you are unauthorized", fmt.Sprintf("%d", http.StatusUnauthorized), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("you are unauthorized", fmt.Sprintf("%d", http.StatusUnauthorized), err.Error()))
		return
	}

	c := http.Cookie{
		Name:     "Refresh-Token",
		Value:    newToken.RefreshToken,
		HttpOnly: true,
		MaxAge:   int(newToken.RtExpires),
	}

	w.Header().Set("Authorization", "Bearer "+newToken.AccessToken)
	http.SetCookie(w, &c)

	w.WriteHeader(http.StatusOK)
}
