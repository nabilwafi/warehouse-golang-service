package repositories

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
)

type UserRepository interface {
	Save(ctx context.Context, tx *sql.Tx, user *domain.User) error
	FindByID(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*domain.User, error)
	FindByEmail(ctx context.Context, tx *sql.Tx, email string) (*domain.User, error)
	FindAll(ctx context.Context, tx *sql.Tx, userQuery *web.UserQuery) ([]domain.User, error)
}

type ProductRepository interface {
	Save(ctx context.Context, tx *sql.Tx, product *web.ProductCreateRequest) error
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Product, error)
	FindByID(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*domain.Product, error)
	Update(ctx context.Context, tx *sql.Tx, product *domain.Product) error
	Delete(ctx context.Context, tx *sql.Tx, id uuid.UUID) error
	IncreaseStock(ctx context.Context, tx *sql.Tx, productID uuid.UUID, quantity int64) error
	DecreaseStock(ctx context.Context, tx *sql.Tx, productID uuid.UUID, quantity int64) error
}

type LocationRepository interface {
	Save(ctx context.Context, tx *sql.Tx, location *web.LocationCreateRequest) error
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Location, error)
}

type OrderRepository interface {
	Save(ctx context.Context, tx *sql.Tx, order *domain.Order) error
	FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Order, error)
	FindByID(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*domain.Order, error)
}
