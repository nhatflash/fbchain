package graph

import (
	"context"
	"slices"

	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/middleware"
	"github.com/nhatflash/fbchain/security"
)

func CheckAuthenticatedUser(ctx context.Context) (*security.JwtAccessClaims, error) {
	claims, ok := ctx.Value(middleware.UserKey{}).(*security.JwtAccessClaims)
	if !ok || claims == nil {
		return nil, appErr.UnauthorizedError("Authentication is required.")
	}
	return claims, nil
}


func CheckRoleAccessUser(ctx context.Context, roles ...string) (*security.JwtAccessClaims, error) {
	claims , err := CheckAuthenticatedUser(ctx)
	if err != nil {
		return nil, err
	}
	if !slices.Contains(roles, claims.Role) {
		return nil, appErr.ForbiddenError("You are not allowed to perform this action.")
	}
	return claims, nil
}