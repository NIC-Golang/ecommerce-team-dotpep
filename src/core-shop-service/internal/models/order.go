package models

import (
	"time"
)

type OrdersResponse struct {
	Orders []Order `json:"orders"`
}
type Order struct {
	ID          string      `json:"order_id"`
	UserID      string      `json:"user_id"`
	Products    []OrderItem `json:"products"`
	TotalAmount float64     `json:"total_amount"`
	Status      string      `json:"status"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ProductID string  `json:"product_id"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price"`
}
