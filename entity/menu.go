package entity

import (
	"fmt"
	"restaurant-service/dto"
	"restaurant-service/pkg/errs"
	"time"
)



type MenuType string

const (
	FOOD MenuType = "food"
	BEVERAGE MenuType = "beverage"
)

type Menu struct {
	MenuSerial string `json:"menuSerial"`
	Type MenuType `json:"type"`
	Name string `json:"name"`
	Stock int32 `json:"stock"`
	Price int32 `json:"price"`
	RestaurantSerial string `json:"restaurantSerial"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}


func(m *Menu) ValidateStock(requestedStock int32) errs.MessageErr {
	if m.Stock < requestedStock {
		errMessage := fmt.Sprintf("%s doesn't have enough stock", m.Name)
		return errs.NewBadRequest(errMessage)
	}

	return nil
}


func (m *Menu) ToGetMenusByRestaurantSerialResponseDTO(menus []*Menu) []*dto.GetMenusByRestaurantSerialResponse {

	var data = []*dto.GetMenusByRestaurantSerialResponse{}


	for _, eachMenu := range menus {

		item := &dto.GetMenusByRestaurantSerialResponse {
			MenuSerial: eachMenu.MenuSerial,
			Type: string(eachMenu.Type),
			Name: eachMenu.Name,
			Stock: eachMenu.Stock,
			Price: eachMenu.Price,
			CreatedAt: eachMenu.CreatedAt,
			UpdatedAt: eachMenu.UpdatedAt,
		}
		data = append(data, item)
	}

	return data
}	