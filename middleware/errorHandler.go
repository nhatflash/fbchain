package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	appErr "github.com/nhatflash/fbchain/error"
)

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if (len(c.Errors) > 0) {
			err := c.Errors[0].Err
			if e, ok := err.(*appErr.ErrorResponse); ok {
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
							msg = fmt.Sprintf("The field '%s' is not in valid email format", e.Field())
						case "name":
							msg = fmt.Sprintf("Field '%s' must not include numbers or special characters", e.Field())
						case "phone":
							msg = fmt.Sprintf("Field '%s' should be a valid Vietnamese phone number", e.Field())
						case "identity":
							msg = fmt.Sprintf("Field '%s' must be a valid Vietnamese identity number", e.Field())
						case "postalcode":
							msg = fmt.Sprintf("Field '%s' must be a valid postal code", e.Field())
						case "price":
							msg = fmt.Sprintf("Field '%s' should has be in valid format e.g. '100.99'", e.Field())
						case "number":
							msg = fmt.Sprintf("Field '%s' must not be lower than 0", e.Field())
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