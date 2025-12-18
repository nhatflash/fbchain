package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/model"
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
func (uc *UserController) ChangeProfile(c *gin.Context) {
	var updtProfileReq client.UpdateProfileRequest
	var err error
	if err = c.ShouldBindJSON(&updtProfileReq); err != nil {
		c.Error(err)
		return 
	}
	var updatedUser *model.User
	updatedUser, err = uc.UserService.ChangeProfile(c.Request.Context(), updtProfileReq.FirstName, updtProfileReq.LastName, updtProfileReq.Birthdate, updtProfileReq.Gender, updtProfileReq.Phone, updtProfileReq.Identity, updtProfileReq.Address, updtProfileReq.PostalCode, updtProfileReq.ProfileImage)
	if err != nil {
		c.Error(err)
		return
	}
	res := helper.MapToUserResponse(updatedUser)
	api.SuccessMessage(http.StatusOK, "Profile updated successfully.", res, c)
}
