package middleware

import (
	"auth/domain"
	"context"
	"encoding/json"
	"net/http"
)

// ValidateLogin ...
func ValidateLogin(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		login := &domain.Login{}
		if err := json.NewDecoder(r.Body).Decode(login); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		err := login.Validate()
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(context.Background(), domain.LoginValidateCtXKey, login)))
	})
}
