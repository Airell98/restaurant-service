package rest

import (
	"net/http"
	"restaurant-service/dto"
	"restaurant-service/pkg/errs"
	"restaurant-service/service"

	"github.com/gin-gonic/gin"
)




type authRestHandler struct {
	service service.AuthService
}

func newAuthHandler(authService service.AuthService) authRestHandler {
	return authRestHandler{
		service: authService,
	}
}


func (a authRestHandler) CustomerLogin(c *gin.Context) {
	var customer dto.CustomerLoginRequest

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errs.NewUnprocessibleEntityError())
		return
	}


	successRes, err := a.service.CustomerLogin(&customer)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, successRes)
}

func (a authRestHandler) CustomerRegister(c *gin.Context) {
	var customer dto.CustomerRegisterRequest

	if err := c.ShouldBindJSON(&customer); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errs.NewUnprocessibleEntityError())
		return
	}

	successRes, err := a.service.CustomerRegister(&customer)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, successRes)
}



func (a authRestHandler) RestaurantRegister(c *gin.Context) {
	var restaurant dto.RestaurantRegisterRequest

	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errs.NewUnprocessibleEntityError())
		return
	}

	successRes, err := a.service.RestaurantRegister(&restaurant)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, successRes)
}

func (a authRestHandler) RestaurantLogin(c *gin.Context) {
	var restaurant dto.RestaurantLoginRequest

	if err := c.ShouldBindJSON(&restaurant); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errs.NewUnprocessibleEntityError())
		return
	}


	successRes, err := a.service.RestaurantLogin(&restaurant)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, successRes)
}