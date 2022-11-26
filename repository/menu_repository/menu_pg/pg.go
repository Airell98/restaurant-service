package menu_pg

import (
	"database/sql"
	"fmt"
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
	"restaurant-service/repository/menu_repository"
)


type menuPG struct {
	db *sql.DB
}

func NewMenuPG(db *sql.DB) menu_repository.MenuRepository {
	return &menuPG {
		db: db,
	}
}


func (m *menuPG) GetMenusByRestaurantSerial(serial string) ([]*entity.Menu, errs.MessageErr) {
	const getAllMenusByRestaurantSerial = 
	`
		SELECT menu_serial, name, type, stock, price, created_at, updated_at from menus
		WHERE restaurant_serial = $1;
	`
	rows, err := m.db.Query(getAllMenusByRestaurantSerial, serial)


	var menus = []*entity.Menu{}

	if err != nil {
		return nil,errs.NewInternalServerErrorr("something went wrong")
	}

	for rows.Next() {
		var menu entity.Menu

		err := rows.Scan(&menu.MenuSerial, &menu.Name, &menu.Type, &menu.Stock, &menu.Price, &menu.CreatedAt, &menu.UpdatedAt)

		if err != nil {
			
			return nil, errs.NewInternalServerErrorr("something went wrong")
		}

		menus = append(menus, &menu)
	}

	return menus, nil
}

func(m *menuPG) GetMenuBySerial(serial string) (*entity.Menu, errs.MessageErr) {

	var menu entity.Menu
	const getMenuBySerialQuery = `
		SELECT menu_serial, name, type, stock, price, restaurant_serial, created_at, updated_at from menus
		WHERE menu_serial = $1;
	`

	row := m.db.QueryRow(getMenuBySerialQuery, serial)

	err := row.Scan(&menu.MenuSerial, &menu.Name, &menu.Type, &menu.Stock, &menu.Price, &menu.RestaurantSerial, &menu.CreatedAt, &menu.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := fmt.Sprintf("menu with serial %s doesn't exist", serial)
			return nil,  errs.NewNotFoundError(errMsg)
		}
			return nil, errs.NewInternalServerErrorr("something went wrong")
	}

	return &menu, nil
}

func (m *menuPG) CreateMenu(menu *entity.Menu) errs.MessageErr {
	const createMenuQuery = `
	INSERT INTO menus
	(
		menu_serial,
		type,
		stock,
		price,
		restaurant_serial,
		name
	)
	VALUES ($1, $2, $3, $4, $5, $6)
`	
	_, err := m.db.Exec(createMenuQuery, menu.MenuSerial, menu.Type, menu.Stock, menu.Price, menu.RestaurantSerial, menu.Name)

	if err != nil {
		return errs.NewInternalServerErrorr("something went wrong")
	}
	return nil
}



func (m *menuPG) GetMenus() ([]*entity.Menu, errs.MessageErr) {
	const getMenusQuery = `SELECT menu_serial, name, type, stock, price, restaurant_serial, created_at, updated_at from menus ORDER BY created_at ASC;`

	var menus = []*entity.Menu{}

	rows, err := m.db.Query(getMenusQuery)

	if err != nil {
		return nil,errs.NewInternalServerErrorr("something went wrong")
	}


	for rows.Next() {
		var menu entity.Menu

		err = rows.Scan(&menu.MenuSerial, &menu.Name, &menu.Type, &menu.Stock, &menu.Price, &menu.RestaurantSerial, &menu.CreatedAt, &menu.UpdatedAt)

		if err != nil {
			return nil,errs.NewInternalServerErrorr("something went wrong")
		}
		menus = append(menus, &menu)
	}

	return menus, nil

}