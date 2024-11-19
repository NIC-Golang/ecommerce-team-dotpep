package models

import (
	"time"
)

type Product struct {
	Name        *string   `json:"product_name" validate:"required,min=2,max=100"`
	ID          *int      `json:"id" validate:"required,min=2,max=100"`
	Description *string   `json:"product_decription" validate:"min=3,max=100"`
	Price       *float64  `json:"product_price" validate:"min=1"`
	SKU         *string   `json:"product_sku" validate:"required"`
	Quantity    *int      `json:"product_quantity" validate:"required"`
	Created_at  time.Time `json:"product_created_at"`
	Update_at   time.Time `json:"product_updated_at"`
}
