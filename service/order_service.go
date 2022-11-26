package service

import (
	"net/http"
	"restaurant-service/dto"
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
	"restaurant-service/pkg/helpers"
	"restaurant-service/repository/menu_repository"
	"restaurant-service/repository/order_repository"
)

type OrderService interface {
	CreateOrder(customerSerial string,ordersPayload []dto.CreateOrderRequest) (*dto.CreateOrderResponse, errs.MessageErr)
	PurchaseOrder(orderPayload *dto.PurchaseOrderRequest) (*dto.PurchaseOrderResponse, errs.MessageErr)
	GetCustomerOrderHistory(customerSerial string) ([]order_repository.OrderHistory , errs.MessageErr)
}


type orderService struct {
	menuRepo menu_repository.MenuRepository
	orderRepo order_repository.OrderRepository
}

func NewOrderService(menuRepo menu_repository.MenuRepository, orderRepo order_repository.OrderRepository)OrderService {
	return &orderService{
		menuRepo: menuRepo,
		orderRepo: orderRepo,
	}
}

func (o *orderService) 	GetCustomerOrderHistory(customerSerial string) ([]order_repository.OrderHistory , errs.MessageErr) {
	orderHistory, err := o.orderRepo.GetCustomerOrderHistory(customerSerial)

	if err != nil {
		return nil, err
	}

	return orderHistory, nil
}
 
func (o *orderService) 	CreateOrder(customerSerial string, ordersPayload []dto.CreateOrderRequest) (*dto.CreateOrderResponse, errs.MessageErr) {
	if len(ordersPayload) == 0 {
		return nil, errs.NewBadRequest("you haven't added any order data yet")
	}

	var cart entity.Cart
	var totalOrderPrice int32 = 0


	orderDetailRequest := []*dto.DetailedOrderSchemaRequest{}

	orderSerial := ""

	for {
		orderSerial = helpers.GenerateBaseSerial("ORD")
		_, err := o.orderRepo.GetOrderBySerial(orderSerial)

		if err != nil {
			if err.Status() == http.StatusNotFound {
				break
			}

			if err.Status() == http.StatusInternalServerError {
				return nil, err
			}
		}
	}

	for _, eachOrder := range ordersPayload {
		
		menu, err := o.menuRepo.GetMenuBySerial(eachOrder.MenuSerial)

		if err != nil {
			return nil, err
		}

		err = menu.ValidateStock(eachOrder.Amount)
		
		if err != nil {
			return nil, err
		}


		totalOrderPrice += (menu.Price * eachOrder.Amount)

		orderDetailSerial := helpers.GenerateBaseSerial("CRT")
		orderDetail := &dto.DetailedOrderSchemaRequest{
			MenuSerial: menu.MenuSerial,
			Amount: eachOrder.Amount,
			MenuTotalPrice: menu.Price * eachOrder.Amount,
			OrderSerial: orderSerial,
			CartSerial: orderDetailSerial,
		}

		orderDetailRequest = append(orderDetailRequest, orderDetail)
	}

	carts := cart.GenerateCartRequest(orderDetailRequest)
	
	singleOrder := &entity.Order{
		TotalPrice: totalOrderPrice,
		Status: "unpurchased",
		CustomerSerial: customerSerial,
		OrderSerial: orderSerial,
	}


	err := o.orderRepo.CreateOrder(carts, singleOrder)

	
if err != nil {
		return nil, err
	}


	response := &dto.CreateOrderResponse{
		Message: "you orders have been successfully created",
	}

	return response, nil
}



func (o *orderService) PurchaseOrder(orderPayload *dto.PurchaseOrderRequest) (*dto.PurchaseOrderResponse, errs.MessageErr)  {
	err := helpers.ValidateStruct(orderPayload)


	if err != nil {
		return nil, err
	}

	_, err = o.orderRepo.GetOrderBySerial(orderPayload.OrderSerial)

	if err != nil {
		return nil, err
	}

	carts, err := o.orderRepo.GetCartsByOrderSerial(orderPayload.OrderSerial)


	if err != nil {
		return nil, err
	}

	menus := []*entity.Menu{}

	for _, eachCart := range carts {
		menu, err := o.menuRepo.GetMenuBySerial(eachCart.MenuSerial)

		if err != nil {
			return nil, err
		}
		err = menu.ValidateStock(eachCart.Amount)

		if err != nil{
			return nil, err
		}	

		menu.Stock = menu.Stock - eachCart.Amount
		menus = append(menus, menu)
	}

	err = o.orderRepo.PurchaseOrders(menus, orderPayload.OrderSerial)
	if err != nil {
		return nil, err
	}

	response := &dto.PurchaseOrderResponse{
		Message: "your orders has been successfully purchased",
	}

	return response, nil
}


