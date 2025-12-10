package service

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	appError "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/repository"
	"github.com/nhatflash/fbchain/security"
)

func GetCurrentUser(c *gin.Context, db *sql.DB) (*model.User, error) {
	defaultClaims, exists := c.Get("user")
	if !exists {
		return nil, appError.UnauthorizedError("User is not authenticated")
	}
	claims := defaultClaims.(*security.JwtAccessClaims)
	email := claims.Email
	user, err := repository.GetUserByEmail(email, db)
	if err != nil {
		return nil, appError.NotFoundError("User not found")
	}
	return user, nil
}