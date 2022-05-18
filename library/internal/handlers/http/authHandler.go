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

	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	loginURL        = "/login"
	logoutURL       = "/logout"
	getUserURL      = "/getUser"
	isLoggedInURL   = "/isloggedin"
	refreshTokenURL = "/refreshtoken"
	registrationURL = "/registration"
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
	router.POST(registrationURL, ah.Registration)
}

func (ah *authHandler) Login(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	req := &domain.Login{}
	if err := json.NewDecoder(r.Body).Decode(req); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Error(apperror.NewAppError("Error occured while parsing json", fmt.Sprintf("%d", http.StatusInternalServerError), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("Error occured while parsing json", fmt.Sprintf("%d", http.StatusUnprocessableEntity), err.Error()))
		return
	}

	login, err := ah.Service.GetByUsername(req.Username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apperror.NewAppError("Error occured while getting user", fmt.Sprintf("%d", http.StatusBadRequest), err.Error()))
		return
	}

	err = domain.CheckPasswordHash(login.Password, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Error(apperror.NewAppError("Error occured while checking password", fmt.Sprintf("%d", http.StatusBadRequest), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("wrong password", fmt.Sprintf("%d", http.StatusBadRequest), err.Error()))
		return
	}

	tk, err := jwt.CreateToken(login.Username)
	if err != nil {
		w.WriteHeader(http.StatusUnprocessableEntity)
		ah.logger.Error(apperror.NewAppError("Error occured while createing tokens", fmt.Sprintf("%d", http.StatusInternalServerError), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("Error occured while createing tokens", fmt.Sprintf("%d", http.StatusInternalServerError), err.Error()))
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
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode("Successfully logged out")

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
	w.WriteHeader(http.StatusOK)
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
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode("You are already logged in")

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

func (ah *authHandler) Registration(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	w.Header().Set("Content-Type", "application/json")
	login := &domain.Login{}
	if err := json.NewDecoder(r.Body).Decode(login); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ah.logger.Error(apperror.NewAppError("Error occured while parsing json", fmt.Sprintf("%d", http.StatusInternalServerError), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("Error occured while parsing json", fmt.Sprintf("%d", http.StatusInternalServerError), err.Error()))
		return
	}

	err := login.Validate()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ah.logger.Error(apperror.NewAppError("Error occured while validation data", fmt.Sprintf("%d", http.StatusBadRequest), err.Error()))
		json.NewEncoder(w).Encode(apperror.NewAppError("Error occured while validation data", fmt.Sprintf("%d", http.StatusBadRequest), err.Error()))
		return
	}

	uuid, err := ah.Service.Create(login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(apperror.NewAppError("Error occured while registration", fmt.Sprintf("%d", http.StatusBadRequest), err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode("Created user with uuid " + uuid)
}
