package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/service"
)

type AdminController struct {
	AdminService 			service.IAdminService
}

func NewAdminController(as service.IAdminService) *AdminController {
	return &AdminController{
		AdminService: as,
	}
}

// @Summary Create New Staff API
// @Produce json
// @Accept json
// @Param request body client.CreateStaffRequest true "CreateStaff body"
// @Security BearerAuth
// @Router /admin/signup/staff [post]
func (ac *AdminController) CreateNewStaff(c *gin.Context) {
	var req client.CreateStaffRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}

	res, err := ac.AdminService.HandleCreateNewStaff(c.Request.Context(), &req)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, "Staff registered successfully.", res, c)
}