package services

import (
	"database/sql"
	"errors"

	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/repositories"
)

type LocationService interface {
	Save(location *dto.Location) (int, error)
	GetAll(pagination *web.PaginationRequest) ([]*dto.Location, int, error)
}

type locationServiceImpl struct {
	location repositories.LocationRepository
}

func NewLocationService(location repositories.LocationRepository) LocationService {
	return &locationServiceImpl{
		location: location,
	}
}

func (s *locationServiceImpl) Save(location *dto.Location) (int, error) {
	locationData, code, err := s.location.FindByName(location.Name)
	if err != nil {
		if err != sql.ErrNoRows {
			return code, err
		}
	}

	if locationData != nil {
		return 401, errors.New("location name is exists")
	}

	if err := s.location.Save(location); err != nil {
		return 500, err
	}

	return 201, nil
}

func (s *locationServiceImpl) GetAll(pagination *web.PaginationRequest) ([]*dto.Location, int, error) {
	locations, err := s.location.GetAllLocation(pagination)
	if err != nil {
		return nil, 500, err
	}

	return locations, 200, nil
}
