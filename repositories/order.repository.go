package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/models/domain"
)

type orderRepositoryImpl struct{}

func NewOrderRepository() OrderRepository {
	return &orderRepositoryImpl{}
}

func (r *orderRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, order *domain.Order) error {
	query := `INSERT INTO orders(id, type, product_id, quantity) VALUES ($1, $2, $3, $4)`
	_, err := tx.ExecContext(ctx, query, order.ID, order.Type, order.ProductID, order.Quantity)
	return err
}

func (r *orderRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Order, error) {
	query := `SELECT id, type, product_id, quantity FROM orders`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []domain.Order
	for rows.Next() {
		var order domain.Order
		if err := rows.Scan(&order.ID, &order.Type, &order.ProductID, &order.Quantity); err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}
	return orders, nil
}

func (r *orderRepositoryImpl) FindByID(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*domain.Order, error) {
	query := `SELECT id, type, product_id, quantity FROM orders WHERE id = $1`
	row := tx.QueryRowContext(ctx, query, id)

	var order domain.Order
	if err := row.Scan(&order.ID, &order.Type, &order.ProductID, &order.Quantity); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("order not found")
		}
		return nil, err
	}
	return &order, nil
}
