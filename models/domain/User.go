package domain

import "github.com/google/uuid"

type UserRole string

const (
	UserRoleAdmin UserRole = "admin"
	UserRoleStaff UserRole = "staff"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	Email    string    `json:"email"`
	Password string    `json:"-"`
	Name     string    `json:"name"`
	Role     UserRole  `json:"role"`
}
