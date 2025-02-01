package web

import "github.com/nabilwafi/warehouse-management-system/models/domain"

type UserQuery struct {
	Size int
	Page int
	Name string
}

type Login struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Register struct {
	Email    string      `json:"email" binding:"required,email"`
	Password string      `json:"password" binding:"required"`
	Role     domain.UserRole `json:"role,omitempty" binding:"required,oneof=admin staff"`
	Name     string      `json:"name" binding:"required"`
}
