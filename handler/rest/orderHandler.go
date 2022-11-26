package rest

import (
	"net/http"
	"restaurant-service/dto"
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
	"restaurant-service/pkg/helpers"
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


// CreateOrder godoc
// @Tags orders
// @Description This is used by customers to create their order data
// @ID create-order
// @Param Authorization header string true "Insert the customer token here" default(Bearer <Add access token here>)
// @Accept json
// @Produce json
// @Param RequestBody body []dto.CreateOrderRequest true "request body json"
// @Success 201 {object} dto.CreateOrderResponse
// @Router /order [post]
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

// GetCustomerOrderHistory godoc
// @Tags orders
// @Description This is used by customers to get their order history data
// @ID get-customer-order-history
// @Param Authorization header string true "Insert the customer token here" default(Bearer <Add access token here>)
// @Accept json
// @Produce json
// @Success 200 {array} order_repository.OrderHistory
// @Router /order/customer/history [get]
func (o orderRestHandler) GetCustomerOrderHistory(c *gin.Context) {
	customerData :=  c.MustGet("customerData").(entity.Customer)


	successRes, err := o.service.GetCustomerOrderHistory(customerData.CustomerSerial)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, successRes)
}


// PurchaseOrders godoc
// @Tags orders
// @Description This is used by customers to purchase their order data
// @ID purchase-orders
// @Param Authorization header string true "Insert the customer token here" default(Bearer <Add access token here>)
// @Accept json
// @Produce json
// @Param RequestBody body dto.PurchaseOrderRequest true "request body json"
// @Success 200 {object} dto.PurchaseOrderResponse
// @Router /order/purchase [put]
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


// GetRestaurantPurchaseHistoryByMonthAndYear godoc
// @Tags orders
// @Description This is used by restaurants to get their purchase history data
// @ID get-restaurant-purchase-history-by-month-and-year
// @Param Authorization header string true "Insert the restaurant token here" default(Bearer <Add access token here>)
// @Accept json
// @Produce json
// @Param month query int true "month query"
// @Param year query int true "year query"
// @Success 200 {array} dto.PurchaseHistoryResponse
// @Router /order/restaurant/history [get]
func (o orderRestHandler) GetRestaurantPurchaseHistoryByMonthAndYear(c *gin.Context) {
	restaurantData :=  c.MustGet("restaurantData").(entity.Restaurant)

	monthParam, err := helpers.GetQueryParam(c,"month")

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	yearParam,err := helpers.GetQueryParam(c,"year")

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	request := &dto.PurchaseHistoryRequest{
		Month: uint8(monthParam),
		Year: uint32(yearParam),
	}

	response, err := o.service.GetRestaurantPurchaseHistoryByMonthAndYear(restaurantData.RestaurantSerial, request)

	if err != nil {
		c.JSON(err.Status(), err)
		return
	}

	c.JSON(http.StatusOK, response)
}
