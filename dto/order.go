package dto



type CreateOrderRequest struct {
	MenuSerial string `json:"menuSerial"`
	Amount int32 `json:"amount"`
}

type DetailedOrderSchemaRequest struct {
	MenuSerial string `json:"menu_serial"`
	Amount int32 `json:"amount"`
	MenuTotalPrice int32 `json:"menuTotalPrice"`
	OrderSerial string `json:"orderSerial"`
	CartSerial string `json:"cartSerial"`
}

type CreateOrderResponse struct {
	Message string `json:"message"`
}


type PurchaseOrderRequest struct {
	OrderSerial string `json:"orderSerial"`
}


type OrderHistoryResponse struct {
	
}

type PurchaseOrderResponse struct {
	Message string `json:"message"`
}


