package dto

import "time"


type CreateMenuRequest struct {
	Type string`json:"type"`
	Stock int32 `json:"stock"`
	Price int32 `json:"price"`
	RestaurantSerial string `json:"restaurantSerial"`
}




type CreateMenuResponse struct {
	Messsage string `json:"message"`
}

type GetMenusByRestaurantSerialResponse struct {
	MenuSerial string `json:"menuSerial"`
	Type string `json:"type"`
	Name string `json:"name"`
	Stock int32 `json:"stock"`
	Price int32 `json:"price"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}