package security

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/nhatflash/fbchain/enum"
	"github.com/nhatflash/fbchain/model"
)

type JwtAccessClaims struct {
	UserId int64  `json:"userId"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

type JwtRefreshClaims struct {
	UserId int64  `json:"userId"`
	Type   string `json:"type"`
	jwt.RegisteredClaims
}

func GenerateJwtAccessToken(u *model.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtAccessExpiration := os.Getenv("JWT_ACCESS_EXPIRATION_MIN")

	expiration, convErr := strconv.Atoi(jwtAccessExpiration)
	if convErr != nil {
		return "", errors.New("access expiration string converted fail")
	}
	userRole := getUserRole(u.Role)
	accessClaims := JwtAccessClaims{
		UserId:           u.Id,
		Email:            u.Email,
		Role:             userRole,
		Type:             "ACCESS",
		RegisteredClaims: buildClaims(u, expiration),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	secret := []byte(jwtSecret)
	accessTokenStr, err := accessToken.SignedString(secret)
	if err != nil {
		return "", errors.New("access token secret signed fail")
	}
	return accessTokenStr, nil
}

func GenerateJwtRefreshToken(u *model.User) (string, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	jwtRefreshExpiration := os.Getenv("JWT_REFRESH_EXPIRATION_MIN")

	expiration, convErr := strconv.Atoi(jwtRefreshExpiration)
	if convErr != nil {
		return "", errors.New("refresh expiration string converted fail")
	}

	refreshClaims := JwtRefreshClaims{
		UserId:           u.Id,
		Type:             "REFRESH",
		RegisteredClaims: buildClaims(u, expiration),
	}
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	secret := []byte(jwtSecret)
	refreshTokenStr, err := refreshToken.SignedString(secret)
	if err != nil {
		return "", errors.New("refresh token secret signed fail")
	}
	return refreshTokenStr, nil
}

func ValidateJwtAccessToken(accessTokenStr string) (*JwtAccessClaims, error) {
	jwtSecret := os.Getenv("JWT_SECRET")
	claims := &JwtAccessClaims{}
	accessToken, err := jwt.ParseWithClaims(accessTokenStr, claims, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})
	if err != nil {
		return nil, err
	}
	if !accessToken.Valid {
		return nil, errors.New("access token invalid")
	}
	return claims, nil
}

func buildClaims(u *model.User, expiration int) jwt.RegisteredClaims {
	return jwt.RegisteredClaims{
		Issuer:    "fbchain",
		Subject:   strconv.FormatInt(u.Id, 10),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		NotBefore: jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(expiration))),
	}
}

func getUserRole(role *enum.Role) string {
	var userRole string
	switch *role {
	case enum.ROLE_ADMIN:
		userRole = "ADMIN"
	case enum.ROLE_MANAGER:
		userRole = "MANAGER"
	case enum.ROLE_STAFF:
		userRole = "STAFF"
	case enum.ROLE_TENANT:
		userRole = "TENANT"
	default:
		userRole = "RESTAURANT_STAFF"
	}
	return userRole
}
