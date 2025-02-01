package domain

import (
	"github.com/google/uuid"
)

type OrderType string

const (
	OrderTypeReceiving OrderType = "receiving"
	OrderTypeShipping  OrderType = "shipping"
)

type Order struct {
	ID        uuid.UUID `json:"id" form:"id" binding:"required,uuid"`
	Type      OrderType `json:"type" form:"type" binding:"required,oneof=receiving shipping"`
	ProductID uuid.UUID `json:"product_id" form:"product_id" binding:"required,uuid"`
	Quantity  int64     `json:"quantity" form:"quantity" binding:"required,min=0"`
	Product   *Product  `json:"product,omitempty"`
}
