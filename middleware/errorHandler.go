package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
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
			vErrs, ok := err.(validator.ValidationErrors)
			if ok {
				var errMsg []string
				for _, e := range vErrs {
					var msg string

					switch e.Tag() {
						case "required":
							msg = fmt.Sprintf("The field '%s' is required", e.Field())
						case "email":
							msg = "Invalid email format"
						case "name":
							msg = fmt.Sprintf("Person '%s' must not include numbers or special characters", e.Field())
						case "phone":
							msg = "Phone number is not valid with Vietnamese phone number"
						case "identity":
							msg = "Identity number is not valid"
						case "postalcode":
							msg = "Postal code is not valid"
						default:
							msg = fmt.Sprintf("Field '%s' validation failed on tag '%s'", e.Field(), e.Tag())
					}
					errMsg = append(errMsg, msg)
				}
				c.JSON(http.StatusBadRequest, gin.H{
					"status": http.StatusBadRequest,
					"error": "VALIDATION_ERROR",
					"message": errMsg,
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