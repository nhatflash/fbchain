package middleware

import (
	"slices"
	"errors"
	"strings"
	"github.com/gin-gonic/gin"
	appError "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/security"
)

func JwtAccessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.Error(appError.UnauthorizedError("Missing authorization header."))
			c.Abort()
			return
		}
		token, tokenErr := getTokenFromHeader(authHeader)
		if tokenErr != nil {
			c.Error(appError.UnauthorizedError(tokenErr.Error()))
			c.Abort()
			return
		}
		claims, claimsErr := security.ValidateJwtAccessToken(token)
		if claimsErr != nil {
			c.Error(appError.UnauthorizedError(claimsErr.Error()))
			c.Abort()
			return
		}
		c.Set("user", claims)
		c.Next()
	}
}


func RoleBasedHandler(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		defaultClaims, exists := c.Get("user")
		if !exists {
			c.Error(appError.UnauthorizedError("User is not authenticated"))
			c.Abort()
			return
		}
		claims := defaultClaims.(*security.JwtAccessClaims)
		if !slices.Contains(roles, claims.Role) {
			c.Error(appError.ForbiddenError("You are not allowed to perform this action."))
			c.Abort()
			return
		}
		c.Next()
	}
}

func getTokenFromHeader(authHeader string) (string, error) {
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}
	token := parts[1]
	return token, nil
}

