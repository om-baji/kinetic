package shared

import (
	"fmt"
	"net/http"
	"runtime"
)

func Panic(message string) {
	panic(message)
}

func Panicf(format string, args ...any) {
	panic(fmt.Sprintf(format, args...))
}

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func PanicIfNil(value any, name string) {
	if value == nil {
		panic(fmt.Sprintf("%s must not be nil", name))
	}
}

func Unreachable() {
	_, file, line, _ := runtime.Caller(1)
	panic(fmt.Sprintf("unreachable code reached at %s:%d", file, line))
}

func Must[T any](v T, err error) T {
	HandleErr(err)
	return v
}

type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *APIError) Error() string {
	return e.Message
}

func NewNotFoundError(message string) *APIError {
	return &APIError{Code: http.StatusNotFound, Message: message}
}

func NewConflictError(message string) *APIError {
	return &APIError{Code: http.StatusConflict, Message: message}
}

func NewValidationError(message string) *APIError {
	return &APIError{Code: http.StatusBadRequest, Message: message}
}

func NewInternalError(message string) *APIError {
	return &APIError{Code: http.StatusInternalServerError, Message: message}
}
