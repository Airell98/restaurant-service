package order_repository

import (
	"restaurant-service/dto"
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
)



type OrderRepository interface {
	CreateOrder(carts []*entity.Cart, order *entity.Order) errs.MessageErr
	GetOrderBySerial(serial string) (*entity.Order,errs.MessageErr ) 
	GetCartsByOrderSerial(serial string) ([]*entity.Cart,errs.MessageErr ) 
	PurchaseOrders(menus []*entity.Menu, orderSerial string) errs.MessageErr
	GetCustomerOrderHistory(customerSerial string) ([]OrderHistory , errs.MessageErr) 
	GetRestaurantPurchaseHistoryByMonthAndYear(restaurantSerial string, month uint8, year uint32) ([]dto.PurchaseHistoryResponse, errs.MessageErr)
}