package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
)

var (

	db *sql.DB
	err error
)

func createRequiredTables() {
	customerTable := `
		CREATE TABLE IF NOT EXISTS customers (
			customer_serial VARCHAR(191) PRIMARY KEY,
			username varchar(191) UNIQUE NOT NULL,
			password TEXT NOT NULL,
			created_at timestamptz NOT NULL DEFAULT (now()),
			updated_at timestamptz NOT NULL DEFAULT (now())
		);
	`

	restaurantTable := `
	CREATE TABLE IF NOT EXISTS restaurants (
		restaurant_serial VARCHAR(191) PRIMARY KEY,
		username VARCHAR(191) UNIQUE NOT NULL,
		address TEXT NOT NULL,
		password TEXT NOT NULL,
		created_at timestamptz NOT NULL DEFAULT (now()),
		updated_at timestamptz NOT NULL DEFAULT (now())
	);
`

menuTable := `
CREATE TABLE IF NOT EXISTS menus (
	menu_serial VARCHAR(191) PRIMARY KEY,
	type VARCHAR(10) NOT NULL,
	stock SMALLINT DEFAULT 0,
	name VARCHAR(191) NOT NULL,
	price INTEGER NOT NULL DEFAULT 0,
	restaurant_serial VARCHAR(191) NOT NULL,
	created_at timestamptz NOT NULL DEFAULT (now()),
	updated_at timestamptz NOT NULL DEFAULT (now()),
	CONSTRAINT menu_restaurant_fk
		FOREIGN KEY(restaurant_serial) 
			REFERENCES restaurants(restaurant_serial)
			  ON DELETE SET NULL
);
`

	orderTable := `
	CREATE TABLE IF NOT EXISTS orders (
		order_serial VARCHAR(191) PRIMARY KEY,
		customer_serial VARCHAR(191) NOT NULL,
		restaurant_serial VARCHAR(191) NOT NULL,
		total_price SMALLINT NOT NULL,
		status VARCHAR(11) NOT NULL,
		created_at timestamptz NOT NULL DEFAULT (now()),
		updated_at timestamptz NOT NULL DEFAULT (now()),
		CONSTRAINT orders_customer_fk
			FOREIGN KEY(customer_serial) 
				REFERENCES customers(customer_serial)
			  		ON DELETE SET NULL,
		CONSTRAINT orders_restaurant_fk
			FOREIGN KEY(restaurant_serial) 
				REFERENCES restaurants(restaurant_serial)
					ON DELETE SET NULL
	);
`

	cartTable := `
	CREATE TABLE IF NOT EXISTS carts (
		cart_serial VARCHAR(191) PRIMARY KEY,
		order_serial VARCHAR(191) NOT NULL,
		menu_serial VARCHAR(191) NOT NULL,
		amount SMALLINT NOT NULL,
		total_price INTEGER NOT NULL,
		created_at timestamptz NOT NULL DEFAULT (now()),
		updated_at timestamptz NOT NULL DEFAULT (now()),
		CONSTRAINT carts_order_fk
			FOREIGN KEY(order_serial) 
				REFERENCES orders(order_serial)
				  ON DELETE SET NULL,
		CONSTRAINT orders_menu_fk
			FOREIGN KEY(menu_serial) 
			  REFERENCES menus(menu_serial)
				ON DELETE SET NULL
);	
`


		
	createTableQueries := fmt.Sprintf("%s %s %s %s %s", customerTable, restaurantTable,menuTable , orderTable, cartTable)
	_, err = db.Exec(createTableQueries)

	if err != nil {
		log.Fatal("error while required tables =>", err.Error())
	}
}



func InitializeDB() {

	username := os.Getenv("PGUSER")
	password := os.Getenv("PGPASSWORD")
	host := os.Getenv("PGHOST")
	dbName := 	os.Getenv("PGDATABASE")
	dbPort := os.Getenv("PGPORT")

	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable",  username, password, host,  dbPort, dbName,)

	db, err = sql.Open("postgres", dsn)

	if err != nil {
		log.Fatal("error connecting to database",err.Error())
	}

	err = db.Ping()

	if err != nil {
		log.Fatal("error while trying to ping the database connection",err.Error())
	}

	fmt.Println("successfully connected to my database")
	createRequiredTables()
}


func GetDB() *sql.DB {
	return db
}