package controller

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/service"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/model"
	_ "github.com/nhatflash/fbchain/docs"
)

type RestaurantController struct {
	UserService 		service.IUserService
	RestaurantService 	service.IRestaurantService
}

func NewRestaurantController(us service.IUserService, rs service.IRestaurantService) *RestaurantController {
	return &RestaurantController{
		UserService: us,
		RestaurantService: rs,
	}
}


// @Summary Create restaurant API
// @Accept json
// @Produce json
// @Param request body client.CreateRestaurantRequest true "CreateRestaurant body"
// @Success 201 {object} client.RestaurantResponse
// @Failure 400 {object} error
// @Security BearerAuth
// @Router /tenant/restaurant [post]
func (rc RestaurantController) CreateRestaurant(c *gin.Context) {
	var createRestaurantReq client.CreateRestaurantRequest
	var err error
	if err = c.ShouldBindJSON(&createRestaurantReq); err != nil {
		c.Error(err)
		return
	}
	var currUser *model.User
	currUser, err = rc.UserService.GetCurrentUser(c)
	if err != nil {
		c.Error(err)
		return
	}
	var res *client.RestaurantResponse
	res, err = rc.RestaurantService.HandleCreateRestaurant(&createRestaurantReq, currUser.Id)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, "Restaurant created successfully.", res, c)
}