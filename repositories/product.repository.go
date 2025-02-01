package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/utils"
)

type productRepositoryImpl struct{}

func NewProductRepository() ProductRepository {
	return &productRepositoryImpl{}
}

func (r *productRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, product *web.ProductCreateRequest) error {
	uuid := utils.GenerateUUID()

	query := `INSERT INTO products (id, name, sku, quantity, location_id) VALUES ($1, $2, $3, $4, $5)`
	_, err := tx.ExecContext(ctx, query, uuid, product.Name, product.SKU, product.Quantity, product.LocationID)
	if err != nil {
		return errors.New("failed to insert product")
	}
	return nil
}

func (r *productRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Product, error) {
	query := `SELECT id, name, sku, quantity, location_id FROM products`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.New("failed to query products")
	}
	defer rows.Close()

	var products []domain.Product
	for rows.Next() {
		var product domain.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.SKU, &product.Quantity, &product.LocationID); err != nil {
			return nil, errors.New("failed to scan product")
		}
		products = append(products, product)
	}

	return products, nil
}

func (r *productRepositoryImpl) FindByID(ctx context.Context, tx *sql.Tx, id uuid.UUID) (*domain.Product, error) {
	query := `SELECT id, name, sku, quantity, location_id FROM products WHERE id = $1`
	row := tx.QueryRowContext(ctx, query, id)

	var product domain.Product
	if err := row.Scan(&product.ID, &product.Name, &product.SKU, &product.Quantity, &product.LocationID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, errors.New("failed to scan product")
	}

	return &product, nil
}

func (r *productRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, product *domain.Product) error {
	query := `UPDATE products SET name = $1, sku = $2, quantity = $3, location_id = $4 WHERE id = $5`
	_, err := tx.ExecContext(ctx, query, product.Name, product.SKU, product.Quantity, product.LocationID, product.ID)
	if err != nil {
		return errors.New("failed to update product")
	}
	return nil
}

func (r *productRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`
	_, err := tx.ExecContext(ctx, query, id)
	if err != nil {
		return errors.New("failed to delete product")
	}
	return nil
}

func (r *productRepositoryImpl) IncreaseStock(ctx context.Context, tx *sql.Tx, productID uuid.UUID, quantity int64) error {
	query := `UPDATE products SET quantity = quantity + $1 WHERE id = $2`

	result, err := tx.ExecContext(ctx, query, quantity, productID)
	if err != nil {
		return errors.New("failed to update stock: " + err.Error())
	}

	// Memeriksa jumlah baris yang terpengaruh
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to get affected rows: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("no product found with the given ID")
	}

	return nil
}

func (r *productRepositoryImpl) DecreaseStock(ctx context.Context, tx *sql.Tx, productID uuid.UUID, quantity int64) error {
	query := `UPDATE products SET quantity = quantity - $1 WHERE id = $2`
	result, err := tx.ExecContext(ctx, query, quantity, productID)
	if err != nil {
		return errors.New("failed to decrease stock: " + err.Error())
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.New("failed to get affected rows: " + err.Error())
	}

	if rowsAffected == 0 {
		return errors.New("no rows updated, insufficient stock or product may not exist")
	}

	return nil
}
