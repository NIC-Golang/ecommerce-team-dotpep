package models

type OrdersResponse struct {
	Orders []Order `json:"orders"`
}
type Order struct {
	ID          int         `json:"order_id"`
	UserID      int         `json:"user_id"`
	Products    []OrderItem `json:"products"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
