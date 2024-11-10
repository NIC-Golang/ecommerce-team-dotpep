package models

type Role string

const (
	RoleCustomer      Role = "CUSTOMER"
	RoleShopClient    Role = "SHOP_CLIENT"
	RoleAdministrator Role = "ADMIN"
)

type CLient struct {
	ID           *string `json:"client_id" validate:"required"`
	Name         *string `json:"client_name" validate:"required, min=5, max=100"`
	LastName     *string `json:"client_last_name" validate:"required, min=5, max=100"`
	Email        *string `json:"client_email" validate:"required, min=5,max=100"`
	Phone        *string `json:"client_phone" validate:"required"`
	Type         Role    `json:"client_type" validate:"required,oneof=CUSTOMER SHOP_CLIENT ADMIN"`
	Token        *string `json:"token"`
	RefreshToken *string `json:"refresh_token"`
}
