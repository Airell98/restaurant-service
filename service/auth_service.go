package service

import (
	"restaurant-service/dto"
	"restaurant-service/entity"
	"restaurant-service/pkg/errs"
	"restaurant-service/pkg/helpers"
	"restaurant-service/repository/auth_repository"

	"github.com/gin-gonic/gin"
)


type AuthService interface {
	CustomerRegister(customerPayload *dto.CustomerRegisterRequest) (*dto.CustomerRegisterResponse, errs.MessageErr)
	CustomerLogin(customerPayload *dto.CustomerLoginRequest) (*dto.CustomerLoginResponse, errs.MessageErr)
	RestaurantRegister(restaurantPayload *dto.RestaurantRegisterRequest) (*dto.RestaurantRegisterResponse,  errs.MessageErr)
	RestaurantLogin(restaurantPayload *dto.RestaurantLoginRequest) (*dto.RestaurantLoginResponse,  errs.MessageErr)
	RestaurantAuthentication() gin.HandlerFunc
	CustomerAuthentication() gin.HandlerFunc
}


type authService struct {
	authRepo auth_repository.AuthRepository
}


func NewAuthService(authRepo auth_repository.AuthRepository) AuthService  {
	return &authService{
		authRepo: authRepo,
	}
}


func (a *authService) RestaurantAuthentication() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context)  {
		var restaurant *entity.Restaurant = &entity.Restaurant{}

		errData := errs.NewNotAuthenticated("you're not authenticated")

		tokenStr := ctx.Request.Header.Get("Authorization")

		err := restaurant.VerifyToken(tokenStr)


		
		if err != nil {
			
			ctx.AbortWithStatusJSON(errData.Status(), errData)
			return
		}

		if restaurant.Role != "restaurant" {
			ctx.AbortWithStatusJSON(errData.Status(), errData)
			return
		}

		
		_, err = a.authRepo.FindRestaurantBySerialAndUsername(restaurant.RestaurantSerial, restaurant.Username)


		if err != nil {
			ctx.AbortWithStatusJSON(errData.Status(), errData)
			return
		}



		ctx.Set("restaurantData", *restaurant)
		ctx.Next()
	})
}


func (a *authService) CustomerAuthentication() gin.HandlerFunc {
	return gin.HandlerFunc(func(ctx *gin.Context)  {
		var customer *entity.Customer = &entity.Customer{}

		errData := errs.NewNotAuthenticated("you're not authenticated")

		tokenStr := ctx.Request.Header.Get("Authorization")

		err := customer.VerifyToken(tokenStr)


		if err != nil {
			
			ctx.AbortWithStatusJSON(errData.Status(), errData)
			return
		}

		if customer.Role != "customer" {
			ctx.AbortWithStatusJSON(errData.Status(), errData)
			return
		}

		
		ctx.Set("customerData", *customer)
		ctx.Next()
	})
}
func (a *authService) CustomerLogin(customerPayload *dto.CustomerLoginRequest) (*dto.CustomerLoginResponse,  errs.MessageErr)  {
	err := helpers.ValidateStruct(customerPayload)
	if err != nil {
		return nil, err
	}

	customer, err := a.authRepo.FindCustomerByUsername(customerPayload.Username)

	if err != nil {
		return nil, errs.NewNotAuthenticated("invalid username/password")
	}

	isValidPassword := customer.ComparePassword(customerPayload.Password)

	if !isValidPassword {
		return nil, errs.NewNotAuthenticated("invalid username/password")
	}

	token := customer.GenerateToken()

	response := &dto.CustomerLoginResponse{
		Token: token,
	}

	return response, nil
}


func (a *authService) CustomerRegister(customerPayload *dto.CustomerRegisterRequest) (*dto.CustomerRegisterResponse,  errs.MessageErr) {
	err := helpers.ValidateStruct(customerPayload)

	if err != nil {
		return nil, err
	}

	customerSerial := helpers.GenerateBaseSerial("CST")
	customer := &entity.Customer{
		Username: customerPayload.Username,
		Password: customerPayload.Password,
		CustomerSerial: customerSerial,
	}

	err = customer.HashPass()

	if err != nil {
		return nil, err
	}

	err = a.authRepo.CustomerRegister(customer)

	if err != nil {
		return nil, err
	}

	response := dto.CustomerRegisterResponse{
		Message: "Your data has been successfully registered",
	}

	return &response, nil
}


func (a *authService) RestaurantRegister(restaurantPayload *dto.RestaurantRegisterRequest) (*dto.RestaurantRegisterResponse,  errs.MessageErr) {
	err := helpers.ValidateStruct(restaurantPayload)

	if err != nil {
		return nil, err
	}

	restaurantSerial := helpers.GenerateBaseSerial("RST")
	restaurant := &entity.Restaurant{
		Username: restaurantPayload.Username,
		Password: restaurantPayload.Password,
		Address: restaurantPayload.Address,
		RestaurantSerial: restaurantSerial,
	}

	err = restaurant.HashPass()

	if err != nil {
		return nil, err
	}

	err = a.authRepo.RestaurantRegister(restaurant)

	if err != nil {
		return nil, err
	}

	response := dto.RestaurantRegisterResponse{
		Message: "Your data has been successfully registered",
	}

	return &response, nil
}


func (a *authService) RestaurantLogin(restaurantPayload *dto.RestaurantLoginRequest) (*dto.RestaurantLoginResponse,  errs.MessageErr) {
	err := helpers.ValidateStruct(restaurantPayload)

	if err != nil {
		return nil, err
	}

	restaurant, err := a.authRepo.FindRestaurantByUsername(restaurantPayload.Username)

	if err != nil {
		return nil, errs.NewNotAuthenticated("invalid username/password")
	}

	isValidPassword := restaurant.ComparePassword(restaurantPayload.Password)

	if !isValidPassword {
		return nil, errs.NewNotAuthenticated("invalid username/password")
	}

	token := restaurant.GenerateToken()

	response := &dto.RestaurantLoginResponse{
		Token: token,
	}

	return response, nil
}







