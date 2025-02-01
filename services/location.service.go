package services

import (
	"context"
	"database/sql"
	"errors"

	"github.com/go-playground/validator/v10"
	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/repositories"
)

type locationServiceImpl struct {
	repo      repositories.LocationRepository
	db        *sql.DB
	validator *validator.Validate
}

func NewLocationService(repo repositories.LocationRepository, validator *validator.Validate, db *sql.DB) LocationService {
	return &locationServiceImpl{repo: repo, db: db, validator: validator}
}

func (s *locationServiceImpl) AddLocation(ctx context.Context, location *web.LocationCreateRequest) error {
	if err := s.validator.Struct(location); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.Save(ctx, tx, location); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *locationServiceImpl) GetAllLocations(ctx context.Context) ([]domain.Location, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	locations, err := s.repo.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	return locations, tx.Commit()
}
