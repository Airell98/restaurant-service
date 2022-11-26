package rest

import (
	"net/http"
	"restaurant-service/dto"
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
	"restaurant-service/service"

	"github.com/gin-gonic/gin"
)

type menuRestHandler struct {
	service service.MenuService
}


func newMenuHandler(menuService service.MenuService) menuRestHandler {
	return menuRestHandler{
		service: menuService,
	}
}



func (m menuRestHandler) GetMenusByRestaurantSerial(c *gin.Context) {
	restaurantData :=  c.MustGet("restaurantData").(entity.Restaurant)


	menus, err := m.service.GetMenusByRestaurantSerial(restaurantData.RestaurantSerial)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, menus)
}

func (m menuRestHandler) CreateMenu(c *gin.Context) {
	restaurantData :=  c.MustGet("restaurantData").(entity.Restaurant)
	var menu  dto.CreateMenuRequest

	if err := c.ShouldBindJSON(&menu); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errs.NewUnprocessibleEntityError())
		return
	}

	successRes, err := m.service.CreateMenu(restaurantData.RestaurantSerial,&menu)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, successRes)
}