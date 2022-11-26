package order_repository

import "restaurant-service/entity"

type OrderHistory struct {
	Order entity.Order
	Carts []entity.Cart
}




