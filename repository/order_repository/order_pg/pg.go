package order_pg

import (
	"database/sql"
	"fmt"
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
	"restaurant-service/repository/order_repository"
)



type orderPG struct {
	db *sql.DB
}

func NewOrderPG(db *sql.DB) order_repository.OrderRepository{
	return &orderPG {
		db: db,
	}
}

func (o *orderPG) GetCustomerOrderHistory(customerSerial string) ([]order_repository.OrderHistory , errs.MessageErr) {
	var orderHistories = []order_repository.OrderHistory{} 

	const getOrderByCustomerSerial = `
		SELECT order_serial, customer_serial, total_price, status, created_at, updated_at from orders
		WHERE customer_serial = $1;
	`

	const getCartByQuery = 
	`	
		SELECT cart_serial, order_serial, menu_serial, amount,total_price, created_at, updated_at from carts
		WHERE order_serial = $1;
	`

	rows, err := o.db.Query(getOrderByCustomerSerial, customerSerial)

	if err != nil {
		return nil, errs.NewInternalServerErrorr("something went wrong")
	}
	
	for rows.Next() {
		var order = order_repository.OrderHistory{} 

		err := rows.Scan(&order.Order.OrderSerial, &order.Order.CustomerSerial, &order.Order.TotalPrice, &order.Order.Status, &order.Order.CreatedAt, &order.Order.UpdatedAt)

		if err != nil {
			return nil, errs.NewInternalServerErrorr("something went wrong")
		}

		orderHistories = append(orderHistories, order)
	}
	for i, eachOrder := range orderHistories {
		var carts = []entity.Cart{}
		rows, err = o.db.Query(getCartByQuery, eachOrder.Order.OrderSerial)

		if err != nil {
			return nil, errs.NewInternalServerErrorr("something went wrong")
		}
		for rows.Next() {
			var cart = entity.Cart{}

			err = rows.Scan(&cart.CartSerial, &cart.OrderSerial, &cart.MenuSerial, &cart.Amount, &cart.TotalPrice, &cart.CreatedAt, &cart.UpdatedAt)

			if err != nil {
				return nil, errs.NewInternalServerErrorr("something went wrong")
			}

			carts = append(carts, cart)
		}

		orderHistories[i].Carts = append(orderHistories[i].Carts, carts...)
	}
	return orderHistories, nil
}

func (o *orderPG) PurchaseOrders(menus []*entity.Menu, orderSerial string) errs.MessageErr {
	tx, err := o.db.Begin()

	if err != nil {
		return errs.NewInternalServerErrorr("something went wrong")
	}

	for _, eachMenu := range menus {
		const updateMenuStockQuery = 
		`
			UPDATE menus
			SET stock = $2
			WHERE menu_serial = $1;
		`
		_ , err := tx.Exec(updateMenuStockQuery, eachMenu.MenuSerial, eachMenu.Stock)

		if err != nil {
			tx.Rollback()
			return errs.NewInternalServerErrorr("something went wrong")
		}
	}

	const updateOrderStatusQuery = `
		UPDATE orders
		SET status = $2
		WHERE order_serial = $1;
	`

	_ , err = tx.Exec(updateOrderStatusQuery, orderSerial, "purchased")

	if err != nil {
		tx.Rollback()
		return errs.NewInternalServerErrorr("something went wrong")
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errs.NewInternalServerErrorr("something went wrong")
	}

	return nil
}

func (o *orderPG) GetCartsByOrderSerial(serial string) ([]*entity.Cart,errs.MessageErr ) {
	var getCartsByOrderSerialQuery = `
		SELECT cart_serial, order_serial, menu_serial, amount,total_price, created_at, updated_at from carts
		WHERE order_serial = $1;
	`
	var carts = []*entity.Cart{}
	rows, err := o.db.Query(getCartsByOrderSerialQuery, serial)

	if err != nil {
		return nil,errs.NewInternalServerErrorr("something went wrong")
	}

	for rows.Next() {
		var cart entity.Cart

		err = rows.Scan(&cart.CartSerial, &cart.OrderSerial, &cart.MenuSerial, &cart.Amount, &cart.TotalPrice, &cart.CreatedAt, &cart.UpdatedAt)
		if err != nil {
			
			return nil, errs.NewInternalServerErrorr("something went wrong")
		}

		carts  = append(carts, &cart)
	}


	return carts, nil
}	

func (o *orderPG) 	GetOrderBySerial(serial string) (*entity.Order,errs.MessageErr ) {
	var order entity.Order

	var getOrderBySerialQuery = `
		SELECT order_serial, customer_serial, total_price, status, created_at, updated_at from orders
		WHERE order_serial = $1;
	`
	
	row := o.db.QueryRow(getOrderBySerialQuery, serial)

	err := row.Scan(&order.OrderSerial, &order.CustomerSerial, &order.TotalPrice, &order.Status, &order.CreatedAt, &order.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			errMsg := fmt.Sprintf("order with serial %s doesn't exist", serial)
			return nil,  errs.NewNotFoundError(errMsg)
		}
			return nil, errs.NewInternalServerErrorr("something went wrong")
	}

	return &order, nil
}

func (o *orderPG) CreateOrder(carts []*entity.Cart, order *entity.Order) errs.MessageErr {
	tx, err := o.db.Begin()

	if err != nil {
		return errs.NewInternalServerErrorr("something went wrong")
	}

	const createOrderQuery = `
		INSERT INTO orders
		(
			order_serial,
			customer_serial,
			total_price,
			status
		)
		VALUES ($1, $2, $3, $4)
	`

	const createCartQuery = 
	`
		INSERT INTO carts
		(
			cart_serial,
			order_serial,
			menu_serial,
			amount,
			total_price
		)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err = tx.Exec(createOrderQuery, order.OrderSerial, order.CustomerSerial, order.TotalPrice, order.Status)

	if err != nil {
		tx.Rollback()
		return errs.NewInternalServerErrorr("something went wrong")
	}

	for _, eachCart := range carts {
		_, err = tx.Exec(createCartQuery, eachCart.CartSerial, eachCart.OrderSerial, eachCart.MenuSerial, eachCart.Amount, eachCart.TotalPrice)

		if err != nil {
			fmt.Println("err createCart =>", err)
			tx.Rollback()
			return errs.NewInternalServerErrorr("something went wrong")
		}
	}

	err = tx.Commit()

	if err != nil {
		tx.Rollback()
		return errs.NewInternalServerErrorr("something went wrong")
	}

	return nil
}

