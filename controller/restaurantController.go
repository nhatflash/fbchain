package controller

import (
	"net/http"
	"strconv"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/constant"
	_ "github.com/nhatflash/fbchain/docs"
	"github.com/nhatflash/fbchain/model"
	"github.com/nhatflash/fbchain/service"
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
// @Security BearerAuth
// @Router /restaurant [post]
func (rc *RestaurantController) CreateRestaurant(c *gin.Context) {
	var req client.CreateRestaurantRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
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
	res, err = rc.RestaurantService.HandleCreateRestaurant(c.Request.Context(), &req, currUser.Id)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, constant.RESTAURANT_CREATED_SUCCESS, res, c)
}


// @Summary Add New Restaurant Item API
// @Produce json
// @Accept json
// @Param restaurantId path string true "Restaurant ID"
// @Param request body client.AddRestaurantItemRequest true "AddRestaurantItem body"
// @Security BearerAuth
// @Router /restaurant/{restaurantId}/item [post]
func (rc *RestaurantController) AddNewRestaurantItem(c *gin.Context) {
	var req client.AddRestaurantItemRequest
	var err error
	if err = c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	restaurantIdStr := c.Param("restaurantId")
	var restaurantId int64
	restaurantId, err = strconv.ParseInt(restaurantIdStr, 10, 64)
	if err != nil {
		c.Error(appErr.BadRequestError("Invalid restaurant id."))
		return
	}
	var res *client.RestaurantItemResponse
	res, err = rc.RestaurantService.HandleAddNewRestaurantItem(c.Request.Context(), restaurantId, &req)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, "New item added successfully.", res, c)
}