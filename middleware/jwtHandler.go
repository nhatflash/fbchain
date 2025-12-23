package middleware

import (
	"context"
	"errors"
	"slices"
	"strings"

	"github.com/gin-gonic/gin"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/security"
)

const AUTH_HEADER = "Authorization"

type UserKey struct {

}

func JwtRestHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AUTH_HEADER)
		if authHeader == "" {
			c.Error(appErr.UnauthorizedError("Missing authorization header."))
			c.Abort()
			return
		}
		var err error
		var token string
		token, err = getTokenFromHeader(authHeader)
		if err != nil {
			c.Error(appErr.UnauthorizedError(err.Error()))
			c.Abort()
			return
		}
		var claims *security.JwtAccessClaims
		claims, err = security.ValidateJwtAccessToken(token)
		if err != nil {
			c.Error(appErr.UnauthorizedError(err.Error()))
			c.Abort()
			return
		}
		ctx := context.WithValue(c.Request.Context(), UserKey{}, claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}


func JwtGraphQLHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader(AUTH_HEADER)
		if authHeader == "" {
			c.Next()
			return
		}
		var err error
		var token string
		token, err = getTokenFromHeader(authHeader)
		if err != nil {
			c.Next()
			return
		}
		var claims *security.JwtAccessClaims
		claims, err = security.ValidateJwtAccessToken(token)
		if err != nil {
			c.Next()
			return
		}
		ctx := context.WithValue(c.Request.Context(), UserKey{}, claims)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}


func RoleBasedHandler(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		var claims *security.JwtAccessClaims
		var ok bool
		claims, ok = ctx.Value(UserKey{}).(*security.JwtAccessClaims)
		if !ok || claims == nil {
			c.Error(appErr.UnauthorizedError("Authentication is required."))
			c.Abort()
			return
		}
		if !slices.Contains(roles, claims.Role) {
			c.Error(appErr.ForbiddenError("You are not allowed to perform this action."))
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




