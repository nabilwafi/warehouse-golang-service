package repositories

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
)

type LocationRepository interface {
	Save(location *dto.Location) error
	FindByName(name string) (*dto.Location, int, error)
	GetAllLocation(pagination *web.PaginationRequest) ([]*dto.Location, error)
}

type locationRepositoryImpl struct {
	db *sqlx.DB
}

func NewLocationRepository(db *sqlx.DB) LocationRepository {
	return &locationRepositoryImpl{
		db: db,
	}
}

func (r *locationRepositoryImpl) Save(location *dto.Location) error {
	uuid := uuid.New()

	_, err := r.db.Exec("INSERT INTO public.locations (id, name, capacity) VALUES ($1, $2, $3)", uuid, location.Name, location.Capacity)
	if err != nil {
		return err
	}

	return nil
}

func (r *locationRepositoryImpl) FindByName(name string) (*dto.Location, int, error) {
	var locationData dto.Location

	if err := r.db.QueryRow("SELECT id, name, capacity FROM public.locations WHERE name = $1", name).Scan(&locationData.ID, &locationData.Name, &locationData.Capacity); err != nil {
		if err == sql.ErrNoRows {
			return nil, 404, sql.ErrNoRows
		}

		return nil, 500, err
	}

	return &locationData, 200, nil
}

func (r *locationRepositoryImpl) GetAllLocation(pagination *web.PaginationRequest) ([]*dto.Location, error) {
	var locationData []*dto.Location

	offset := (pagination.Page - 1) * pagination.Size

	rows, err := r.db.Queryx("SELECT id, name, capacity FROM public.locations OFFSET $1 LIMIT $2", offset, pagination.Size)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var location dto.Location
		if err := rows.Scan(&location.ID, &location.Name, &location.Capacity); err != nil {
			return nil, err
		}
		locationData = append(locationData, &location)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return locationData, nil
}
