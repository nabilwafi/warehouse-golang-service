package web

type LocationCreateRequest struct {
	Name     string `json:"name" binding:"required"`
	Capacity int `json:"capacity" binding:"required"`
}
