package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role string

const (
	RoleCustomer      Role = "CUSTOMER"
	RoleShopClient    Role = "SHOP_CLIENT"
	RoleAdministrator Role = "ADMIN"
)

type CLient struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Client_name      *string            `json: "name_client" validate:"required, min=5, max=100"`
	Client_last_name *string            `json: "last_name_client" validate:"required, min=5, max=100"`
	Client_email     *string            `json:"email_client" validate:"required, min=5,max=100"`
	Client_phone     *string            `json:"phone_client" validate:"required"`
	Client_type      Role               `json:"client_type" validate:"required,oneof=CUSTOMER SHOP_CLIENT ADMIN"`
}
