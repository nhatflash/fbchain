package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	appError "github.com/nhatflash/fbchain/error"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if (len(c.Errors) > 0) {
			err := c.Errors[0].Err
			if e, ok := err.(*appError.ErrorResponse); ok {
				c.JSON(e.Status, gin.H{
					"status": e.Status,
					"error": e.Code.Error(),
					"message": e.Message,
				})
				return
			}

			c.JSON(http.StatusInternalServerError, gin.H{
				"status": http.StatusInternalServerError,
				"error": err.Error(),
				"message": "An unexpected error occurred",
			})
		}
	}
}