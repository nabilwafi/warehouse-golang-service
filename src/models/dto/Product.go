package dto

import (
	"github.com/google/uuid"
)

type Product struct {
	ID         uuid.UUID `json:"id" binding:"uuid"`
	Name       string    `json:"name" binding:"required,max=100"`
	SKU        string    `json:"sku" binding:"required,max=100"`
	Quantity   int64     `json:"quantity" binding:"required,min=0"`
	LocationID uuid.UUID `json:"location_id" binding:"required,uuid"`
	Location   *Location `json:"location,omitempty"`
}
