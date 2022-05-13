package domain

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

const LoginValidateCtXKey  = 1

type Login struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Validate ...
func (l *Login) Validate() error {
	return validation.ValidateStruct(
		l,
		validation.Field(&l.Username, validation.Required, validation.By(IsSQL)),
		validation.Field(&l.Password, validation.Required, validation.Length(5, 25), validation.By(IsSQL)),
	)
}
