package rest

import (
	"restaurant-service/database"
	"restaurant-service/repository/auth_repository/auth_pg"
	"restaurant-service/repository/menu_repository/menu_pg"
	"restaurant-service/repository/order_repository/order_pg"
	"restaurant-service/service"

	"github.com/gin-gonic/gin"
)


const port = ":8080"

func StartApp() {
	database.InitializeDB()

	db := database.GetDB()

	authRepo := auth_pg.NewAuthPG(db)

	authService := service.NewAuthService(authRepo)

	authHandler := newAuthHandler(authService)

	menuRepo := menu_pg.NewMenuPG(db)

	menuService := service.NewMenuService(menuRepo)

	menuHandler := newMenuHandler(menuService)

	orderRepo := order_pg.NewOrderPG(db)

	orderService := service.NewOrderService(menuRepo, orderRepo)

	orderHandler := newOrderHandler(orderService)

	route := gin.Default()

	authRoute := route.Group("/auth") 
	{
		authRoute.POST("/customer-register", authHandler.CustomerRegister)
		authRoute.POST("/customer-login", authHandler.CustomerLogin)
		authRoute.POST("/restaurant-register", authHandler.RestaurantRegister)
		authRoute.POST("/restaurant-login", authHandler.RestaurantLogin)
	}

	menuRoute := route.Group("/menu")
	{	
		menuRoute.POST("/", authService.RestaurantAuthentication(), menuHandler.CreateMenu)
		menuRoute.GET("/my-menus", authService.RestaurantAuthentication(), menuHandler.GetMenusByRestaurantSerial)
	}

	orderRoute := route.Group("/order")
	{
		orderRoute.POST("/", orderHandler.CreateOrder)
		orderRoute.PUT("/purchase", orderHandler.PurchaseOrders);
		orderRoute.GET("/history", orderHandler.GetCustomerOrderHistory)
	}

	route.Run(port)
}