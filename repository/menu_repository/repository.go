package menu_repository

import (
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
)



type MenuRepository interface {
	CreateMenu(menu *entity.Menu) errs.MessageErr
	GetMenuBySerial(serial string)(*entity.Menu, errs.MessageErr)
	GetMenusByRestaurantSerial(serial string) ([]*entity.Menu, errs.MessageErr)
}