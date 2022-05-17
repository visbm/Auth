package domain

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"golang.org/x/crypto/bcrypt"
)

const LoginValidateCtXKey = 1

type Login struct {
	UUID     string `json:"uuid,omitempty"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate ...
func (l *Login) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Username, validation.Required, validation.By(IsSQL)),
		validation.Field(&l.Password, validation.Required, validation.Length(8, 15), validation.By(PassReq), validation.By(IsSQL)),
	)
}

// EncryptPassword ...
func EncryptPassword(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

// CheckPasswordHash if passwords are same err=nil
func CheckPasswordHash(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
