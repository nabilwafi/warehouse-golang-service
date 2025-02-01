package web

import (
	"github.com/google/uuid"
)

type OrderCreateRequest struct {
	ProductID uuid.UUID `json:"product_id" form:"product_id" binding:"required,uuid"`
	Quantity  int64     `json:"quantity" form:"quantity" binding:"required,min=0"`
}
