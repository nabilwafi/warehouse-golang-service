package repositories_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/lib/pq"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/repositories"
	"github.com/stretchr/testify/assert"
)

func setupDatabase(t *testing.T) *sql.DB {
	db, err := sql.Open("postgres", "user=user dbname=warehouse sslmode=disable password=password")
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	return db
}

func TestSaveLocation(t *testing.T) {
	db := setupDatabase(t)

	repo := repositories.NewLocationRepository()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctx := context.Background()

	locationRequest := web.LocationCreateRequest{
		Name:     "Warehouse A",
		Capacity: 100,
	}

	err = repo.Save(ctx, tx, &locationRequest)
	assert.NoError(t, err, "Save should not return an error")

	var id, name string
	var capacity int
	err = tx.QueryRowContext(ctx, "SELECT id, name, capacity FROM locations WHERE name = $1", "Warehouse A").Scan(&id, &name, &capacity)
	assert.NoError(t, err, "Failed to query location")
	assert.Equal(t, "Warehouse A", name, "Location name should match")
	assert.Equal(t, 100, capacity, "Location capacity should match")
}

func TestFindAllLocations(t *testing.T) {
	db := setupDatabase(t)

	repo := repositories.NewLocationRepository()

	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	ctx := context.Background()

	_, err = tx.ExecContext(ctx, "INSERT INTO locations (id, name, capacity) VALUES ($1, $2, $3)", "fbf9ade9-5be5-47b2-8291-c78a7848d24a", "Warehouse A", 100)
	assert.NoError(t, err, "Failed to insert location")

	_, err = tx.ExecContext(ctx, "INSERT INTO locations (id, name, capacity) VALUES ($1, $2, $3)", "1d3d328e-5196-4f5c-91e0-d79d83f5b96e", "Warehouse B", 200)
	assert.NoError(t, err, "Failed to insert location")

	locations, err := repo.FindAll(ctx, tx)
	assert.NoError(t, err, "FindAll should not return an error")
	assert.Equal(t, 2, len(locations), "Should return 2 locations")

	assert.Equal(t, "Warehouse A", locations[0].Name, "First location name should match")
	assert.Equal(t, int64(100), locations[0].Capacity, "First location capacity should match")
	assert.Equal(t, "Warehouse B", locations[1].Name, "Second location name should match")
	assert.Equal(t, int64(200), locations[1].Capacity, "Second location capacity should match")
}
