package entity

import (
	"restaurant-service/dto"
	"time"
)


type Cart struct {
	CartSerial string `json:"cartSerial"`
	OrderSerial string `json:"orderSerial"`
	MenuSerial string `json:"menuSerial"`
	Amount int32 `json:"amount"`
	TotalPrice int32 `json:"totalPrice"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}


func (c *Cart) GenerateCartRequest(orderRequests []*dto.DetailedOrderSchemaRequest) []*Cart {
	var carts []*Cart

	for _, eachOrder := range orderRequests {
		cartTemp := &Cart{
			Amount: eachOrder.Amount,
			TotalPrice: eachOrder.MenuTotalPrice,
			MenuSerial: eachOrder.MenuSerial,
			OrderSerial: eachOrder.OrderSerial,
			CartSerial: eachOrder.CartSerial,
		}

		carts = append(carts, cartTemp)
	}

	return carts
}

