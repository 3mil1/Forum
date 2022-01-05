package appError

import (
	"encoding/json"
	"net/http"
)

var (
	ErrIncorrectEmail = NewAppError(nil, "incorrect e-mail address", http.StatusBadRequest)
	ForbiddenError    = NewAppError(nil, "access forbidden", http.StatusForbidden)
)

type AppError struct {
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
	Err        error  `json:"-"`
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Marshal() []byte {
	marshal, err := json.Marshal(e)
	if err != nil {
		return nil
	}
	return marshal
}

func NewAppError(err error, message string, code int) *AppError {
	return &AppError{
		Err:        err,
		Message:    message,
		StatusCode: code,
	}
}

func SystemError(err error) *AppError {
	return NewAppError(err, "internal system error", http.StatusInternalServerError)
}

func InvalidArgumentError(err error, message string) *AppError {
	return NewAppError(err, message, http.StatusBadRequest)
}

func DataBaseError(err error) *AppError {
	return NewAppError(err, "cannot reach database", http.StatusInternalServerError)
}

func NotFoundError(err error, message string) *AppError {
	return NewAppError(err, message, http.StatusNotFound)
}

func UnsupportedError(err error, message string) *AppError {
	return NewAppError(err, message, http.StatusUnsupportedMediaType)
}
