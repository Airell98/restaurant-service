package auth_pg

import (
	"database/sql"
	"fmt"
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
	"restaurant-service/repository/auth_repository"
)



type authPG struct {
	db *sql.DB
}

func NewAuthPG(db *sql.DB) auth_repository.AuthRepository {
	return &authPG{
		db: db,
	}
}





func (a *authPG) FindCustomerByUsername(username string) (*entity.Customer, errs.MessageErr) {
	var customer entity.Customer
	const getCustomerByEmailQuery = `SELECT customer_serial, username, password from customers WHERE username = $1;`

	row := a.db.QueryRow(getCustomerByEmailQuery, username)

	err := row.Scan(&customer.CustomerSerial, &customer.Username, &customer.Password)

	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil,  errs.NewNotFoundError("customer doesn't exist")
		}
		return nil, errs.NewInternalServerErrorr("something went wrong")
	}

	return &customer, nil
}


func (a *authPG) FindRestaurantByUsername(username string) (*entity.Restaurant, errs.MessageErr) {
	var restaurant entity.Restaurant
	const getRestaurantByUsernameQuery = `SELECT restaurant_serial, username, password from restaurants WHERE username = $1;`

	row := a.db.QueryRow(getRestaurantByUsernameQuery, username)

	err := row.Scan(&restaurant.RestaurantSerial, &restaurant.Username, &restaurant.Password)

	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil,  errs.NewNotFoundError("restaurant doesn't exist")
		}
		return nil, errs.NewInternalServerErrorr("something went wrong")
	}

	return &restaurant, nil
}


func (a *authPG) CustomerRegister (customer *entity.Customer)  errs.MessageErr {
	const createCustomerQuery = `
		INSERT INTO customers
		(
			customer_serial,
			username,
			password
		)
		VALUES ($1, $2, $3)
	`

	_, err := a.db.Exec(createCustomerQuery, customer.CustomerSerial, customer.Username, customer.Password)

	if err != nil {
		return errs.NewInternalServerErrorr("something went wrong")
	}

	return nil
}


func (a *authPG) RestaurantRegister (restaurant *entity.Restaurant)  errs.MessageErr {
	const createRestaurantQuery = `
	INSERT INTO restaurants
	(
		restaurant_serial,
		username,
		password,
		address
	)
	VALUES ($1, $2, $3, $4)
`

_, err := a.db.Exec(createRestaurantQuery, restaurant.RestaurantSerial, restaurant.Username, restaurant.Password, restaurant.Address)

if err != nil {
	fmt.Println("err =>",err)
	return errs.NewInternalServerErrorr("something went wrong")
}

return nil
}

func(a *authPG) FindRestaurantBySerialAndUsername(serial string, username string) (*entity.Restaurant, errs.MessageErr) {
	var restaurant entity.Restaurant
	
	const getRestaurantBySerialAndUsernameQuery = `
		SELECT restaurant_serial from restaurants WHERE restaurant_serial = $1 AND username = $2;
	`

	row := a.db.QueryRow(getRestaurantBySerialAndUsernameQuery, serial, username)

	err := row.Scan(&restaurant.RestaurantSerial)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil,  errs.NewNotFoundError("restaurant doesn't exist")
		}
		return nil, errs.NewInternalServerErrorr("something went wrong")
	}

	return &restaurant, nil
}