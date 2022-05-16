package apperror

import (
	"encoding/json"
	"fmt"
)

// AppError ...
type AppError struct {
	Err              error  `json:"-"`
	Message          string `json:"message,omitempty"`
	DeveloperMessage string `json:"developer_message,omitempty"`
	Code             string `json:"code,omitempty"`
}

// NewAppError ...
func NewAppError(message, code, developerMessage string) *AppError {
	return &AppError{
		Err:              fmt.Errorf(message),
		Code:             code,
		Message:          message,
		DeveloperMessage: developerMessage,
	}
}

func (e *AppError) Error() string {
	return e.Err.Error()
}

func (e *AppError) Unwrap() error { return e.Err }

// Marshal ...
func (e *AppError) Marshal() []byte {
	bytes, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return bytes
}

// UnauthorizedError ...
func UnauthorizedError(message string) *AppError {
	return NewAppError("you are unauthorized", "401", message)
}

// BadRequestError ...
func BadRequestError(message string) *AppError {
	return NewAppError("wrong data", "400", message)
}

// APIError ...
func APIError(code, message, developerMessage string) *AppError {
	return NewAppError(message, code, developerMessage)
}
