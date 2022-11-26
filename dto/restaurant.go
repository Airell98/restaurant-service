package dto

type RestaurantRegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Address string `json:"address"`
}


type RestaurantRegisterResponse struct {
	Message string `json:"message"`
}

type RestaurantLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type RestaurantLoginResponse struct {
	Token string `json:"token"`
}