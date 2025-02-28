package models

type OrderItem struct {
	ProductID   string  `json:"product_id" validate:"required"`
	Description string  `json:"description" validate:"required"`
	Quantity    int     `json:"quantity" validate:"required"`
	Price       float64 `json:"price" validate:"required"`
}
