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


type ErrorResponse struct {
	Status		int				
	Code		error			
	Message		string			
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

func (e *ErrorResponse) Error() string {
	return e.Message
}