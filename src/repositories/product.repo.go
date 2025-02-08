package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
)

type ProductRepository interface {
	Save(product *dto.Product) error
	Update(product *dto.Product) error
	FindByID(id uuid.UUID) (*dto.Product, int, error)
	FindByName(name string) (*dto.Product, int, error)
	FindBySKU(sku string) (*dto.Product, int, error)
	GetAllProduct(pagination *web.PaginationRequest) ([]*dto.Product, error)
	Delete(id uuid.UUID) (int, error)
	IncreaseStockWithTransaction(tx *sqlx.Tx, productID uuid.UUID, quantity int64) (int, error)
	DecreaseStockWithTransaction(tx *sqlx.Tx, productID uuid.UUID, quantity int64) (int, error)
}

type productRepositoryImpl struct {
	db *sqlx.DB
}

func NewProductRepository(db *sqlx.DB) ProductRepository {
	return &productRepositoryImpl{
		db: db,
	}
}

func (r *productRepositoryImpl) Save(product *dto.Product) error {
	uuid := uuid.New()

	_, err := r.db.Exec("INSERT INTO public.products (id, name, sku, quantity, location_id) VALUES ($1, $2, $3, $4, $5)", uuid, product.Name, product.SKU, product.Quantity, product.LocationID)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) Update(product *dto.Product) error {
	_, err := r.db.Exec("UPDATE public.products SET name = $2, sku = $3, quantity = $4, location_id = $5 WHERE id = $1", product.ID, product.Name, product.SKU, product.Quantity, product.LocationID)
	if err != nil {
		return err
	}

	return nil
}

func (r *productRepositoryImpl) FindByID(id uuid.UUID) (*dto.Product, int, error) {
	var productData dto.Product

	if err := r.db.QueryRow("SELECT id, name, sku, quantity, location_id FROM public.products WHERE id = $1", id).Scan(&productData.ID, &productData.Name, &productData.SKU, &productData.Quantity, &productData.LocationID); err != nil {
		if err == sql.ErrNoRows {
			return nil, 404, sql.ErrNoRows
		}

		return nil, 500, err
	}

	return &productData, 200, nil
}

func (r *productRepositoryImpl) FindByName(name string) (*dto.Product, int, error) {
	var productData dto.Product

	if err := r.db.QueryRow("SELECT id, name, sku, quantity, location_id FROM public.products WHERE name = $1", name).Scan(&productData.ID, &productData.Name, &productData.SKU, &productData.Quantity, &productData.LocationID); err != nil {
		if err == sql.ErrNoRows {
			return nil, 404, sql.ErrNoRows
		}

		return nil, 500, err
	}

	return &productData, 200, nil
}

func (r *productRepositoryImpl) FindBySKU(sku string) (*dto.Product, int, error) {
	var productData dto.Product

	if err := r.db.QueryRow("SELECT id, name, sku, quantity, location_id FROM public.products WHERE sku = $1", sku).Scan(&productData.ID, &productData.Name, &productData.SKU, &productData.Quantity, &productData.LocationID); err != nil {
		if err == sql.ErrNoRows {
			return nil, 404, sql.ErrNoRows
		}

		return nil, 500, err
	}

	return &productData, 200, nil
}

func (r *productRepositoryImpl) GetAllProduct(pagination *web.PaginationRequest) ([]*dto.Product, error) {
	var productData []*dto.Product

	offset := (pagination.Page - 1) * pagination.Size

	rows, err := r.db.Queryx("SELECT id, name, sku, quantity, location_id FROM public.products OFFSET $1 LIMIT $2", offset, pagination.Size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var product dto.Product
		if err := rows.Scan(&product.ID, &product.Name, &product.SKU, &product.Quantity, &product.LocationID); err != nil {
			return nil, err
		}
		productData = append(productData, &product)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return productData, nil
}

func (r *productRepositoryImpl) Delete(id uuid.UUID) (int, error) {
	result, err := r.db.Exec("DELETE FROM public.products WHERE id = $1", id)
	if err != nil {
		return 500, err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 500, err
	}

	if rowsAffected == 0 {
		return 404, sql.ErrNoRows
	}

	return 200, nil
}

func (r *productRepositoryImpl) IncreaseStockWithTransaction(tx *sqlx.Tx, productID uuid.UUID, quantity int64) (int, error) {
	result, err := tx.Exec("UPDATE products SET quantity = quantity + $1 WHERE id = $2", quantity, productID)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 500, err
	}

	if rowsAffected == 0 {
		return 404, sql.ErrNoRows
	}

	return 200, nil
}

func (r *productRepositoryImpl) DecreaseStockWithTransaction(tx *sqlx.Tx, productID uuid.UUID, quantity int64) (int, error) {
	result, err := tx.Exec("UPDATE products SET quantity = quantity - $1 WHERE id = $2", quantity, productID)
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return 500, err
	}

	if rowsAffected == 0 {
		return 404, sql.ErrNoRows
	}

	return 200, nil
}
