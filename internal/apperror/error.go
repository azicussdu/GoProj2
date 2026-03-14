package apperror

import "net/http"

type AppError struct {
	Code    int
	Message string
	Err     error
}

func (e *AppError) Error() string {
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code int, message string, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

func BadRequest(message string, err error) *AppError {
	return New(http.StatusBadRequest, message, err)
}

func NotFound(msg string, err error) *AppError {
	return New(http.StatusNotFound, msg, err)
}

func Conflict(msg string, err error) *AppError {
	return New(http.StatusConflict, msg, err)
}

func Internal(msg string, err error) *AppError {
	return New(http.StatusInternalServerError, msg, err)
}
