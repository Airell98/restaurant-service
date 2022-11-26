package auth_repository

import (
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
)


type AuthRepository interface {
	CustomerRegister (customer *entity.Customer) errs.MessageErr
	FindCustomerByUsername (username string) ( *entity.Customer, errs.MessageErr)
	RestaurantRegister (customer *entity.Restaurant) errs.MessageErr
	FindRestaurantByUsername (username string) (*entity.Restaurant, errs.MessageErr)
	FindRestaurantBySerialAndUsername (serial string, username string) (*entity.Restaurant, errs.MessageErr)
}