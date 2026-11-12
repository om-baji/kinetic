package shared

import "fmt"

type APIError struct {
	Code    int
	Message string
}

func (e *APIError) Error() string { return e.Message }

func BadRequest(format string, args ...any) *APIError {
	return &APIError{Code: 400, Message: fmt.Sprintf(format, args...)}
}

func NotFound(format string, args ...any) *APIError {
	return &APIError{Code: 404, Message: fmt.Sprintf(format, args...)}
}

func Internal(format string, args ...any) *APIError {
	return &APIError{Code: 500, Message: fmt.Sprintf(format, args...)}
}
