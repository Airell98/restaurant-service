package rest

import (
	"log"
	"os"
	"restaurant-service/database"
	"restaurant-service/docs"
	"restaurant-service/repository/auth_repository/auth_pg"
	"restaurant-service/repository/menu_repository/menu_pg"
	"restaurant-service/repository/order_repository/order_pg"
	"restaurant-service/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"     // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger" // gin-swagger middleware
)




func StartApp() {
	if os.Getenv("APP_ENV") != "production" {
		err := godotenv.Load()

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	} 
	var port = os.Getenv("PORT")




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

	docs.SwaggerInfo.Title = "Restaurant Service"
	docs.SwaggerInfo.Description = "Rest API for a restaurant app"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "restaurant-service-production.up.railway.app"
	docs.SwaggerInfo.Schemes = []string{"https"}

	route.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

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
		menuRoute.GET("/", menuHandler.GetMenus)
		menuRoute.GET("/my-menus", authService.RestaurantAuthentication(), menuHandler.GetMenusByRestaurantSerial)
	}

	orderRoute := route.Group("/order")
	{
		
		orderRoute.POST("/", authService.CustomerAuthentication(),  orderHandler.CreateOrder)
		orderRoute.PUT("/purchase", authService.CustomerAuthentication(),orderHandler.PurchaseOrders);
		orderRoute.GET("/customer/history", authService.CustomerAuthentication(), orderHandler.GetCustomerOrderHistory)
		orderRoute.GET("/restaurant/history", authService.RestaurantAuthentication(), orderHandler.GetRestaurantPurchaseHistoryByMonthAndYear)
	}

	route.Run(":" + port)
}