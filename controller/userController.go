package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/constant"
	_ "github.com/nhatflash/fbchain/docs"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/service"
)

type UserController struct {
	UserService 	service.IUserService	
}


func NewUserController(us service.IUserService) *UserController {
	return &UserController{
		UserService: us,
	}
}

// @Summary Change Profile API
// @Produce json
// @Accept json
// @Param request body client.UpdateProfileRequest true "UpdateProfile body"
// @Success 200 {object} client.UserResponse
// @Failure 400 {object} error
// @Security BearerAuth
// @Router /profile [patch]
func (uc *UserController) ChangeUserProfile(c *gin.Context) {
	var req client.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return 
	}
	ctx := c.Request.Context()
	currUser, err := uc.UserService.FindCurrentUser(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	if !currUser.IsVerified {
		c.Error(appErr.UnauthorizedError("You need to verify your account before doing this action."))
		return
	}
	updatedUser, err := uc.UserService.HandleChangeUserProfile(c.Request.Context(), currUser, &req)
	if err != nil {
		c.Error(err)
		return
	}
	res := helper.MapToUserResponse(updatedUser)
	api.SuccessMessage(http.StatusOK, constant.PROFILE_UPDATED_SUCCESS, res, c)
}
