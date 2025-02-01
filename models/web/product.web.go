package web

import "github.com/google/uuid"

type ProductCreateRequest struct {
	Name       string    `json:"name" form:"name" binding:"required,max=100"`
	SKU        string    `json:"sku" form:"sku" binding:"required,max=100"`
	Quantity   int64     `json:"quantity" form:"quantity" binding:"required,min=0"`
	LocationID uuid.UUID `json:"location_id" form:"location_id" binding:"required,uuid"`
}

type ProductUpdateRequest struct {
	Name       string    `json:"name" form:"name" binding:"required,max=100"`
	SKU        string    `json:"sku" form:"sku" binding:"required,max=100"`
	Quantity   int64     `json:"quantity" form:"quantity" binding:"required,min=0"`
	LocationID uuid.UUID `json:"location_id" form:"location_id" binding:"required,uuid"`
}
