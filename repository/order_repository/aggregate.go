package order_repository

import (
	"restaurant-service/dto"
	"restaurant-service/entity"
)

type OrderHistory struct {
	Order entity.Order
	Carts []entity.Cart
}


type PurchaseHistory struct {
	Cart entity.Cart
	Customer entity.Customer
	Menu entity.Menu
}


func (p PurchaseHistory) ToPurchaseHistoryResponseDTO(purchaseHistoryData []PurchaseHistory) []dto.PurchaseHistoryResponse {
	var result = []dto.PurchaseHistoryResponse {}
	for _, eachItem := range purchaseHistoryData  {
			data := dto.PurchaseHistoryResponse {
				CartSerial: eachItem.Cart.CartSerial,
				OrderSerial: eachItem.Cart.OrderSerial,
				RestaurantSerial: eachItem.Menu.RestaurantSerial,
				MenuSerial: eachItem.Menu.MenuSerial,
				MenuName: eachItem.Menu.Name,
				TotalPrice: eachItem.Cart.TotalPrice,
				Amount: eachItem.Cart.Amount,
				CreatedAt: eachItem.Cart.CreatedAt,
				CustomerName: eachItem.Customer.Username,
			}

			result = append(result, data)
	}

	return result
}
