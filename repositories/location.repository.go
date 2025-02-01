package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/utils"
)

type locationRepositoryImpl struct{}

func NewLocationRepository() LocationRepository {
	return &locationRepositoryImpl{}
}

func (r *locationRepositoryImpl) Save(ctx context.Context, tx *sql.Tx, location *web.LocationCreateRequest) error {
	uuid := utils.GenerateUUID()

	query := `INSERT INTO locations (id, name, capacity) VALUES ($1, $2, $3)`
	_, err := tx.ExecContext(ctx, query, uuid, location.Name, location.Capacity)
	if err != nil {
		return errors.New("failed to insert location")
	}
	return nil
}

func (r *locationRepositoryImpl) FindAll(ctx context.Context, tx *sql.Tx) ([]domain.Location, error) {
	query := `SELECT id, name, capacity FROM locations`
	rows, err := tx.QueryContext(ctx, query)
	if err != nil {
		return nil, errors.New("failed to query locations")
	}
	defer rows.Close()

	var locations []domain.Location
	for rows.Next() {
		var location domain.Location
		if err := rows.Scan(&location.ID, &location.Name, &location.Capacity); err != nil {
			return nil, errors.New("failed to scan location")
		}
		locations = append(locations, location)
	}

	return locations, nil
}
