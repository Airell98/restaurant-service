package rest

import (
	"net/http"
	"restaurant-service/dto"
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
	"restaurant-service/service"

	"github.com/gin-gonic/gin"
)


type orderRestHandler struct {
	service service.OrderService
}


func newOrderHandler(orderService service.OrderService)orderRestHandler {
	return orderRestHandler{
		service: orderService,
	}
}

func (o orderRestHandler) CreateOrder(c *gin.Context) {
	var orders []dto.CreateOrderRequest
	customerData :=  c.MustGet("customerData").(entity.Customer)

	if err := c.ShouldBindJSON(&orders); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errs.NewUnprocessibleEntityError())
		return
	}

	successRes, err := o.service.CreateOrder(customerData.CustomerSerial , orders)


	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusCreated, successRes)
}


func (o orderRestHandler) GetCustomerOrderHistory(c *gin.Context) {
	customerData :=  c.MustGet("customerData").(entity.Customer)


	successRes, err := o.service.GetCustomerOrderHistory(customerData.CustomerSerial)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, successRes)
}

func (o orderRestHandler) PurchaseOrders(c *gin.Context) {
	var order = dto.PurchaseOrderRequest{}

	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusUnprocessableEntity, errs.NewUnprocessibleEntityError())
		return
	}

	successRes, err := o.service.PurchaseOrder(&order)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, successRes)
}