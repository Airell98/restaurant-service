package service

import (
	"restaurant-service/dto"
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
	"restaurant-service/pkg/helpers"
	"restaurant-service/repository/menu_repository"
)


type MenuService interface {
	CreateMenu(menuPayload *dto.CreateMenuRequest) (*dto.CreateMenuResponse, errs.MessageErr)
	GetMenusByRestaurantSerial(restaurantSerial string) ([]*dto.GetMenusByRestaurantSerialResponse, errs.MessageErr )
}



type menuService struct {
	menuRepo menu_repository.MenuRepository
}

func NewMenuService(menuRepo menu_repository.MenuRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
	}
}

func (m *menuService) GetMenusByRestaurantSerial(restaurantSerial string) ([]*dto.GetMenusByRestaurantSerialResponse, errs.MessageErr ) {
	menus , err := m.menuRepo.GetMenusByRestaurantSerial(restaurantSerial)

	if err != nil {
		return nil, err
	}


	var menu = entity.Menu{}


	return menu.ToGetMenusByRestaurantSerialResponseDTO(menus), err
}


func (m *menuService) CreateMenu(menuPayload *dto.CreateMenuRequest) (*dto.CreateMenuResponse, errs.MessageErr) {
	err := helpers.ValidateStruct(menuPayload)

	if err != nil {
		return nil, err
	}

	menuSerial := helpers.GenerateBaseSerial("MNU")
	menu := &entity.Menu{
		MenuSerial: menuSerial,
		Type: entity.MenuType(menuPayload.Type),
		Stock: menuPayload.Stock,
		Price: menuPayload.Price,
		RestaurantSerial: menuPayload.RestaurantSerial,
	}

	err = m.menuRepo.CreateMenu(menu)

	if err != nil {
		return nil, err
	}

	response := &dto.CreateMenuResponse{
		Messsage: "your menu has been successfully created",
	}
	return response, nil
}