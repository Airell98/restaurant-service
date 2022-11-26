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

// CustomerLogin godoc
// @Tags customers
// @Description Customer login endpoint
// @ID customer-login
// @Accept json
// @Produce json
// @Param RequestBody body dto.CustomerLoginRequest true "request body json"
// @Success 200 {object} dto.CustomerLoginResponse
// @Router /auth/customer-login [post]
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


// CustomerRegister godoc
// @Tags customers
// @Description Customer register endpoint
// @ID customer-register
// @Accept json
// @Produce json
// @Param RequestBody body dto.CustomerRegisterRequest true "request body json"
// @Success 201 {object} dto.CustomerRegisterResponse
// @Router /auth/customer-register [post]
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


// RestaurantRegister godoc
// @Tags restaurants
// @Description Restaurant register endpoint
// @ID restaurant-register
// @Accept json
// @Produce json
// @Param RequestBody body dto.RestaurantRegisterRequest true "request body json"
// @Success 201 {object} dto.RestaurantRegisterResponse
// @Router /auth/restaurant-register [post]
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


// RestaurantLogin godoc
// @Tags restaurants
// @Description Restaurant Login endpoint
// @ID restaurant-login
// @Accept json
// @Produce json
// @Param RequestBody body dto.RestaurantLoginRequest true "request body json"
// @Success 200 {object} dto.RestaurantLoginResponse
// @Router /auth/restaurant-login [post]
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