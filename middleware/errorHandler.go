package middleware

import (
	"errors"
	"strings"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/error"
	"net/http"
	"log"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if (len(c.Errors) > 0) {
			err := c.Errors[0].Err

			code := http.StatusInternalServerError
			message := "An unexpected error occurred"
			switch {
			case errors.Is(err, error.ErrBadRequest):
				code = http.StatusBadRequest
				message = "The requested resource is invalid"

			case errors.Is(err, error.ErrNotFound):
				code = http.StatusNotFound
				message = "The requested resource is not found"

			case errors.Is(err, error.ErrUnauthorized):
				code = http.StatusUnauthorized
				message = "You must be authenticated first"

			case errors.Is(err, error.ErrForbidden):
				code = http.StatusForbidden
				message = "Access denied"
			case strings.Contains(err.Error(), "timeout"):
				code = http.StatusServiceUnavailable
				message = "Service temporarily unavailable"
			default:
				log.Fatalln("Unhandled error details", err.Error())
			}
			api.ErrorMessage(code, err, message, c)
		}
	}
}