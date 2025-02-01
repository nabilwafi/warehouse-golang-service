package repositories_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/repositories"
	"github.com/stretchr/testify/assert"
)

func TestSaveProduct(t *testing.T) {
	db := setupDatabase(t)

	repo := repositories.NewProductRepository()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctx := context.Background()

	_, err = db.Exec("INSERT INTO locations (id, name, capacity) VALUES ($1, $2, $3)", uuid.New(), "Location A", 100)
	if err != nil {
		t.Fatalf("Failed to insert location: %v", err)
	}

	var locID uuid.UUID
	err = tx.QueryRowContext(ctx, "SELECT id FROM locations WHERE name = $1", "Location A").Scan(&locID)
	assert.NoError(t, err, "Failed to get location ID")

	productRequest := web.ProductCreateRequest{
		Name:       "Product A",
		SKU:        "SKU123",
		Quantity:   100,
		LocationID: locID,
	}

	err = repo.Save(ctx, tx, &productRequest)
	assert.NoError(t, err, "Save should not return an error")

	// Verifikasi data yang disimpan
	var id uuid.UUID
	var name, sku string
	var quantity int
	var locationID uuid.UUID
	err = tx.QueryRowContext(ctx, "SELECT id, name, sku, quantity, location_id FROM products WHERE name = $1", "Product A").Scan(&id, &name, &sku, &quantity, &locationID)
	assert.NoError(t, err, "Failed to query product")
	assert.Equal(t, "Product A", name, "Product name should match")
	assert.Equal(t, "SKU123", sku, "Product SKU should match")
	assert.Equal(t, 100, quantity, "Product quantity should match")
	assert.Equal(t, productRequest.LocationID, locationID, "Product location ID should match")
}

func teardownDatabase(db *sql.DB) {
	db.Exec("TRUNCATE locations CASCADE")
	db.Close()
}

func TestFindAllProducts(t *testing.T) {
	db := setupDatabase(t)
	defer teardownDatabase(db)

	repo := repositories.NewProductRepository()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctx := context.Background()

	// Ambil locationID dari tabel locations
	var locationID uuid.UUID
	err = tx.QueryRowContext(ctx, "SELECT id FROM locations WHERE name = $1", "Location A").Scan(&locationID)
	assert.NoError(t, err, "Failed to get location ID")

	// Insert beberapa data product untuk testing
	_, err = tx.ExecContext(ctx, "INSERT INTO products (id, name, sku, quantity, location_id) VALUES ($1, $2, $3, $4, $5)", uuid.New(), "Product A", "SKU123", 100, locationID)
	assert.NoError(t, err, "Failed to insert product")

	_, err = tx.ExecContext(ctx, "INSERT INTO products (id, name, sku, quantity, location_id) VALUES ($1, $2, $3, $4, $5)", uuid.New(), "Product B", "SKU456", 200, locationID)
	assert.NoError(t, err, "Failed to insert product")

	products, err := repo.FindAll(ctx, tx)
	assert.NoError(t, err, "FindAll should not return an error")
	assert.Equal(t, 2, len(products), "Should return 2 products")

	// Verifikasi data yang diambil
	assert.Equal(t, "Product A", products[0].Name, "First product name should match")
	assert.Equal(t, "SKU123", products[0].SKU, "First product SKU should match")
	assert.Equal(t, int64(100), products[0].Quantity, "First product quantity should match")
	assert.Equal(t, locationID, products[0].LocationID, "First product location ID should match")

	assert.Equal(t, "Product B", products[1].Name, "Second product name should match")
	assert.Equal(t, "SKU456", products[1].SKU, "Second product SKU should match")
	assert.Equal(t, int64(200), products[1].Quantity, "Second product quantity should match")
	assert.Equal(t, locationID, products[1].LocationID, "Second product location ID should match")
}
