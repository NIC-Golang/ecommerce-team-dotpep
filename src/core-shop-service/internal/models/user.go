package models

type Role string

const (
	RoleCustomer      Role = "CUSTOMER"
	RoleShopClient    Role = "SHOP_CLIENT"
	RoleAdministrator Role = "ADMIN"
)

type CLient struct {
	CLient_ID        *string `json:"client_id,omitempty" json:"id"`
	Client_name      *string `json:"client_name" validate:"required, min=5, max=100"`
	Client_last_name *string `json:"client_last_name" validate:"required, min=5, max=100"`
	Client_email     *string `json:"client_email" validate:"required, min=5,max=100"`
	Client_phone     *string `json:"client_phone" validate:"required"`
	Client_type      Role    `json:"client_type" validate:"required,oneof=CUSTOMER SHOP_CLIENT ADMIN"`
}
