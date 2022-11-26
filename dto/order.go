package dto

import "time"



type CreateOrderRequest struct {
	MenuSerial string `json:"menuSerial"`
	Amount int32 `json:"amount"`
}

type DetailedOrderSchemaRequest struct {
	MenuSerial string `json:"menu_serial"`
	Amount int32 `json:"amount"`
	MenuTotalPrice int32 `json:"menuTotalPrice"`
	OrderSerial string `json:"orderSerial"`
	CartSerial string `json:"cartSerial"`
}

type CreateOrderResponse struct {
	Message string `json:"message"`
}


type PurchaseOrderRequest struct {
	OrderSerial string `json:"orderSerial"`
}


type OrderHistoryResponse struct {
	
}

type PurchaseOrderResponse struct {
	Message string `json:"message"`
}


type PurchaseHistoryRequest struct {
	Month uint8 `json:"month" valid:"range(1|12)~month has to be in the range of 1 - 12"`
	Year uint32 `json:"year" valid:"required~year cannot be empty"`
}

type PurchaseHistoryResponse struct {
	CartSerial string `json:"cartSerial"`
	OrderSerial string `json:"orderSerial"`
	MenuSerial string `json:"menuSerial"`
	RestaurantSerial string `json:"restaurantSerial"`
	MenuName string `json:"menuName"`
	TotalPrice int32 `json:"totalPrice"`
	Amount int32 `json:"amount"`
	CustomerName string `json:"customerName"`
	CreatedAt time.Time `json:"createdAt"`
}