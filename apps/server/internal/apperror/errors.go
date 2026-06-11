package apperror

import (
	"errors"
	"fmt"
	"net/http"
)

type AppError struct {
	Code       string
	Message    string
	HTTPStatus int
	Err        error
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func BadRequest(message string, err error) *AppError {
	return &AppError{
		Code:       "BAD_REQUEST",
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Err:        err,
	}
}

func NotFound(message string, err error) *AppError {
	return &AppError{
		Code:       "NOT_FOUND",
		Message:    message,
		HTTPStatus: http.StatusNotFound,
		Err:        err,
	}
}

func Validation(message string, err error) *AppError {
	return &AppError{
		Code:       "VALIDATION_FAILED",
		Message:    message,
		HTTPStatus: http.StatusBadRequest,
		Err:        err,
	}
}

func Unauthorized(message string, err error) *AppError {
	return &AppError{
		Code:       "UNAUTHORIZED",
		Message:    message,
		HTTPStatus: http.StatusUnauthorized,
		Err:        err,
	}
}

func UnsupportedMediaType(message string, err error) *AppError {
	return &AppError{
		Code:       "UNSUPPORTED_MEDIA_TYPE",
		Message:    message,
		HTTPStatus: http.StatusUnsupportedMediaType,
		Err:        err,
	}
}

func TooManyRequests(message string, err error) *AppError {
	return &AppError{
		Code:       "TOO_MANY_REQUESTS",
		Message:    message,
		HTTPStatus: http.StatusTooManyRequests,
		Err:        err,
	}
}

func Internal(message string, err error) *AppError {
	return &AppError{
		Code:       "INTERNAL_ERROR",
		Message:    message,
		HTTPStatus: http.StatusInternalServerError,
		Err:        err,
	}
}

func FromError(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return Internal("An error occurred, please retry", err)
}
