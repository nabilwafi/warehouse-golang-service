package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
)

type OrderRepository interface {
	SaveWithTransaction(tx *sqlx.Tx, order *dto.Order) error
	FindAll(pagination *web.PaginationRequest) ([]*dto.Order, error)
	FindByID(id uuid.UUID) (*dto.Order, int, error)
}

type orderRepositoryImpl struct {
	db *sqlx.DB
}

func NewOrderRepository(db *sqlx.DB) OrderRepository {
	return &orderRepositoryImpl{
		db: db,
	}
}

func (r *orderRepositoryImpl) SaveWithTransaction(tx *sqlx.Tx, order *dto.Order) error {
	uuid := uuid.New()

	_, err := tx.Exec("INSERT INTO public.orders (id, type, product_id, quantity) VALUES ($1, $2, $3, $4)", uuid, order.Type, order.ProductID, order.Quantity)
	if err != nil {
		return err
	}

	return nil
}

func (r *orderRepositoryImpl) FindAll(pagination *web.PaginationRequest) ([]*dto.Order, error) {
	var orderData []*dto.Order

	offset := (pagination.Page - 1) * pagination.Size

	rows, err := r.db.Queryx("SELECT id, type, product_id, quantity FROM public.orders OFFSET $1 LIMIT $2", offset, pagination.Size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order dto.Order
		if err := rows.Scan(&order.ID, &order.Type, &order.ProductID, &order.Quantity); err != nil {
			return nil, err
		}
		orderData = append(orderData, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orderData, nil
}

func (r *orderRepositoryImpl) FindByID(id uuid.UUID) (*dto.Order, int, error) {
	var orderData dto.Order

	if err := r.db.QueryRow("SELECT id, type, product_id, quantity FROM public.orders WHERE id = $1", id).Scan(&orderData.ID, &orderData.Type, &orderData.ProductID, &orderData.Quantity); err != nil {
		if sql.ErrNoRows != nil {
			return nil, 404, sql.ErrNoRows
		}

		return nil, 500, err
	}

	return &orderData, 200, nil
}
