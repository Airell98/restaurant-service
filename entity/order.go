package entity

import (
	"time"
)

type OrderStatus string


const (
	PURCHASED OrderStatus = "purchased"
	UNPURCHASED OrderStatus = "unpurchased"
)

type Order struct {
	OrderSerial string `json:"orderSerial"`	
	CustomerSerial string `json:"customerSerial"`
	RestaurantSerial string `json:"restaurantSerial"`
	TotalPrice int32 `json:"totalPrice"`
	Status OrderStatus `json:"orderStatus"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}


