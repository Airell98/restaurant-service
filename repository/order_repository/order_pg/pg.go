package order_pg

import (
	"database/sql"
	"fmt"
	"restaurant-service/dto"
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


func (o *orderPG) GetRestaurantPurchaseHistoryByMonthAndYear(restaurantSerial string, month uint8, year uint32) ([]dto.PurchaseHistoryResponse, errs.MessageErr){
	const getPurchaseHistoryByWeekAndYearQuery = `
	SELECT cart_serial, c.order_serial as order_serial, c.menu_serial as menu_serial, m.restaurant_serial as restaurant_serial, m.name as menu_name, c.total_price as total_price, amount, c.created_at as created_at, cust.username as customer_username from carts as c
	LEFT JOIN  menus as m ON m.menu_serial = c.menu_serial 
	LEFT JOIN  orders as o ON o.order_serial = c.order_serial
	LEFT JOIN  customers as cust ON cust.customer_serial = o.customer_serial
	WHERE restaurant_serial = $1 AND EXTRACT(MONTH FROM c.created_at) = $2 AND EXTRACT(YEAR FROM c.created_at) = $3;
	`

	var purchaseHistories = []order_repository.PurchaseHistory{}
	var purchaseHistory = order_repository.PurchaseHistory{}

	rows, err := o.db.Query(getPurchaseHistoryByWeekAndYearQuery, restaurantSerial, month, year)

	if err != nil {
		return nil, errs.NewInternalServerErrorr("something went wrong")
	}

	for rows.Next(){
		err = rows.Scan(&purchaseHistory.Cart.CartSerial, 
			&purchaseHistory.Cart.OrderSerial, 
			&purchaseHistory.Menu.MenuSerial, 
			&purchaseHistory.Menu.RestaurantSerial,
			&purchaseHistory.Menu.Name,
			&purchaseHistory.Cart.TotalPrice,
			&purchaseHistory.Cart.Amount,
			&purchaseHistory.Cart.CreatedAt,
			&purchaseHistory.Customer.Username,
		)

		if err != nil {
			return nil, errs.NewInternalServerErrorr("something went wrong")
		}

		purchaseHistories = append(purchaseHistories, purchaseHistory)
	}


	response := purchaseHistory.ToPurchaseHistoryResponseDTO(purchaseHistories)

	return response, nil

}
