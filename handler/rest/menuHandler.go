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


// GetMenusByRestaurantSerial godoc
// @Tags menus
// @Description This is used by restaurants to get all of their menu data
// @ID get-menu-by-restaurant-serial
// @Param Authorization header string true "Insert the restaurant token here" default(Bearer <Add access token here>)
// @Produce json
// @Success 200 {array} dto.GetMenusByRestaurantSerialResponse
// @Router /menu/my-menus [get]
func (m menuRestHandler) GetMenusByRestaurantSerial(c *gin.Context) {
	restaurantData :=  c.MustGet("restaurantData").(entity.Restaurant)


	menus, err := m.service.GetMenusByRestaurantSerial(restaurantData.RestaurantSerial)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, menus)
}


// CreateMenu godoc
// @Tags menus
// @Description This is used by restaurants to add a menu data for them
// @ID create-menu
// @Param Authorization header string true "Insert the restaurant token here" default(Bearer <Add access token here>)
// @Accept json
// @Produce json
// @Param RequestBody body dto.CreateMenuRequest true "request body json"
// @Success 201 {object} dto.CreateMenuResponse
// @Router /menu [post]
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


// GetMenus godoc
// @Tags menus
// @Description Get all menu data endpoint
// @ID get-menus
// @Produce json
// @Success 200 {array} entity.Menu
// @Router /menu [get]
func (m menuRestHandler) GetMenus(c *gin.Context)  {
	successRes, err := m.service.GetMenus()

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}
	c.JSON(http.StatusOK, successRes)
}
