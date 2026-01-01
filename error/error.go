package error

import (
	"errors"
	"net/http"
)

var ErrBadRequest = errors.New("BAD_REQUEST")
var ErrNotFound = errors.New("NOT_FOUND")
var ErrInternal = errors.New("INTERNAL_SERVER_ERROR")
var ErrUnauthorized = errors.New("UNAUTHORIZED")
var ErrForbidden = errors.New("FORBIDDEN")
var ErrValidationFailed = errors.New("VALIDATION_FAILED")
var ErrServiceUnavailable = errors.New("SERVICE_UNAVAILABLE")

type ErrorResponse struct {
	Status		int				`json:"status"`
	Code		error			`json:"code"`
	Message		string			`json:"message"`
}


func BadRequestError(message string) error {
	return &ErrorResponse{
		Status: http.StatusBadRequest,
		Code: ErrBadRequest,
		Message: message,
	}
}

func NotFoundError(message string) error {
	return &ErrorResponse{
		Status: http.StatusNotFound,
		Code: ErrNotFound,
		Message: message,
	}
}

func InternalError(message string) error {
	return &ErrorResponse{
		Status: http.StatusInternalServerError,
		Code: ErrInternal,
		Message: message,
	}
}

func UnauthorizedError(message string) error {
	return &ErrorResponse{
		Status: http.StatusUnauthorized,
		Code: ErrUnauthorized,
		Message: message,
	}
}

func ForbiddenError(message string) error {
	return &ErrorResponse{
		Status: http.StatusForbidden,
		Code: ErrForbidden,
		Message: message,
	}
}

func ServiceUnavailableError(message string) error {
	return &ErrorResponse{
		Status: http.StatusServiceUnavailable,
		Code: ErrServiceUnavailable,
		Message: message,
	}
}

func (e *ErrorResponse) Error() string {
	return e.Message
}