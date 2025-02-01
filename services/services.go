package services

import (
	"context"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
)

type UserService interface {
	Login(ctx context.Context, user *web.Login) (string, error)
	Register(ctx context.Context, user *web.Register) error
	GetMe(ctx context.Context, id uuid.UUID) (*domain.User, error)
	ListUser(ctx context.Context, userQuery *web.UserQuery) ([]domain.User, error)
}

type ProductService interface {
	Create(ctx context.Context, product *web.ProductCreateRequest) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error)
	GetAll(ctx context.Context) ([]domain.Product, error)
	Update(ctx context.Context, product *domain.Product) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type LocationService interface {
	AddLocation(ctx context.Context, location *web.LocationCreateRequest) error
	GetAllLocations(ctx context.Context) ([]domain.Location, error)
}

type OrderService interface {
	ReceiveOrder(ctx context.Context, order *web.OrderCreateRequest) error
	ShipOrder(ctx context.Context, order *web.OrderCreateRequest) error
	GetAllOrders(ctx context.Context) ([]domain.Order, error)
	GetOrderByID(ctx context.Context, id uuid.UUID) (*domain.Order, error)
}
