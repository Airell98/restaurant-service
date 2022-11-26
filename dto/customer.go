package dto


type CustomerRegisterRequest struct {
	Username string `json:"username" valid:"required~username cannot be empty"`
	Password string `json:"password" valid:"required~password cannot be empty"`
}


type CustomerRegisterResponse struct {
	Message string `json:"message"`
}


type CustomerLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type CustomerLoginResponse struct {
	Token string `json:"token"`
}