package controller

import (
	"database/sql"
	"net/http"
	"github.com/gin-gonic/gin"
	"github.com/nhatflash/fbchain/client"
	"github.com/nhatflash/fbchain/service"
	"github.com/nhatflash/fbchain/api"
	_ "github.com/nhatflash/fbchain/docs"
)

type RestaurantController struct {
	Db 			*sql.DB
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

	if reqErr := c.ShouldBindJSON(&createRestaurantReq); reqErr != nil {
		c.Error(reqErr)
		return
	}

	currentUser, userErr := service.GetCurrentUser(c, rc.Db)
	if userErr != nil {
		c.Error(userErr)
		return
	}

	res, resErr := service.HandleCreateRestaurant(&createRestaurantReq, currentUser.Id, rc.Db)
	if resErr != nil {
		c.Error(resErr)
		return
	}
	api.SuccessMessage(http.StatusCreated, "Restaurant created successfully.", res, c)
}