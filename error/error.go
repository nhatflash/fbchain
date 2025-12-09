package error

import "errors"

var ErrBadRequest = errors.New("BAD_REQUEST")
var ErrNotFound = errors.New("NOT_FOUND")
var ErrInternal = errors.New("INTERNAL_SERVER_ERROR")
var ErrUnauthorized = errors.New("UNAUTHORIZED")
var ErrForbidden = errors.New("FORBIDDEN")