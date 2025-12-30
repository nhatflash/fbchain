package controller

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/api"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/constant"
	_ "github.com/nhatflash/fbchain/docs"
	appErr "github.com/nhatflash/fbchain/error"
	"github.com/nhatflash/fbchain/helper"
	"github.com/nhatflash/fbchain/service"
)

type RestaurantController struct {
	UserService 		service.IUserService
	RestaurantService 	service.IRestaurantService
	TenantService 		service.ITenantService
}

func NewRestaurantController(us service.IUserService, rs service.IRestaurantService, ts service.ITenantService) *RestaurantController {
	return &RestaurantController{
		UserService: us,
		RestaurantService: rs,
		TenantService: ts,
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
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	ctx := c.Request.Context()
	currUser, err := rc.UserService.FindCurrentUser(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	if !currUser.IsVerified {
		c.Error(appErr.UnauthorizedError("Please verify your account before doing this action."))
		return
	}

	currTenant, err := rc.TenantService.FindTenantByUserId(ctx, currUser.Id)
	if err != nil {
		c.Error(err)
		return
	} 

	res, err := rc.RestaurantService.HandleCreateRestaurant(ctx, &req, currTenant.Id)
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
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	rIdStr := c.Param("restaurantId")
	restaurantId, err := strconv.ParseInt(rIdStr, 10, 64)
	if err != nil {
		c.Error(appErr.BadRequestError("Invalid restaurant id."))
		return
	}

	ctx := c.Request.Context()
	currUser, err := rc.UserService.FindCurrentUser(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	currTenant, err := rc.TenantService.FindTenantByUserId(ctx, currUser.Id)
	if err != nil {
		c.Error(err)
		return
	}

	res, err := rc.RestaurantService.HandleAddNewRestaurantItem(ctx, restaurantId, currTenant.Id, &req)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusCreated, "New item added successfully.", res, c)
}


// @Summary Add New Restaurant Table API
// @Produce json
// @Accept json
// @Param restaurantId path string true "Restaurant ID"
// @Param request body client.AddRestaurantTableRequest true "AddRestaurantTable body"
// @Security BearerAuth
// @Router /restaurant/{restaurantId}/table [post]
func (rc *RestaurantController) AddNewRestaurantTable(c *gin.Context) {
	var req client.AddRestaurantTableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.Error(err)
		return
	}
	rIdStr := c.Param("restaurantId")
	restaurantId, err := strconv.ParseInt(rIdStr, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	ctx := c.Request.Context()
	currUser, err := rc.UserService.FindCurrentUser(ctx)
	if err != nil {
		c.Error(err)
		return
	}
	currTenant, err := rc.TenantService.FindTenantByUserId(ctx, currUser.Id)
	if err != nil {
		c.Error(err)
		return
	}
	res, err := rc.RestaurantService.HandleAddNewRestaurantTable(ctx, currTenant.Id, restaurantId, &req)
	if err != nil {
		c.Error(err)
		return 
	}
	api.SuccessMessage(http.StatusCreated, "Restaurant table added successfully.", res, c)
}


func (rc *RestaurantController) ShowRestaurantItemsViaQRCode(c *gin.Context) {
	tblIdParam := c.Param("tableId")
	tableId, err := strconv.ParseInt(tblIdParam, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}
	items, err := rc.RestaurantService.HandleShowRestaurantItemsViaQRCode(c.Request.Context(), tableId)
	if err != nil {
		c.Error(err)
		return
	}
	res := helper.MapToRestaurantItemsResponse(items)
	api.SuccessMessage(http.StatusOK, "Restaurant items show successfully.", res, c)
}


// @Summary Get Restaurant Table QR Code
// @Param tableId path string true "Table ID"
// @Param restaurantId path string true "Restaurant ID"
// @Security BearerAuth
// @Router /restaurant/{restaurantId}/table/{tableId}/qrCode [get]
func (rc *RestaurantController) GetTableQrCode(c *gin.Context) {
	resIdParam := c.Param("restaurantId")
	tblIdParam := c.Param("tableId")

	restaurantId, err := strconv.ParseInt(resIdParam, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	tableId, err := strconv.ParseInt(tblIdParam, 10, 64)
	if err != nil {
		c.Error(err)
		return
	}

	ctx := c.Request.Context()
	currUser, err := rc.UserService.FindCurrentUser(ctx)
	if err != nil {
		c.Error(err)
		return
	}

	if !currUser.IsVerified {
		c.Error(appErr.UnauthorizedError("Please verify your account before doing this action."))
		return
	}

	currTenant, err := rc.TenantService.FindTenantByUserId(ctx, currUser.Id)
	if err != nil {
		c.Error(err)
		return
	}

	err = rc.RestaurantService.GetQRCodeOnRestaurantTable(ctx, tableId, currTenant.Id, restaurantId)
	if err != nil {
		c.Error(err)
		return
	}
	api.SuccessMessage(http.StatusOK, "QR Code saved successfully.", nil, c)
}