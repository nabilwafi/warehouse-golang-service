package services_test

import (
	"database/sql"
	"testing"

	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/services"
	mocks "github.com/nabilwafi/warehouse-management-system/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestSaveLocation(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	service := services.NewLocationService(mockRepo)

	location := &dto.Location{
		Name: "Warehouse A",
	}

	t.Run("Success", func(t *testing.T) {
		mockRepo.On("FindByName", location.Name).Return((*dto.Location)(nil), 0, sql.ErrNoRows).Once()
		mockRepo.On("Save", location).Return(nil).Once()

		statusCode, err := service.Save(location)

		assert.NoError(t, err)
		assert.Equal(t, 201, statusCode)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Location Exists", func(t *testing.T) {
		mockRepo.On("FindByName", location.Name).Return(location, 200, nil).Once()

		statusCode, err := service.Save(location)

		assert.Error(t, err)
		assert.Equal(t, 401, statusCode)
		assert.Equal(t, "location name is exists", err.Error())
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Error", func(t *testing.T) {
		mockRepo.On("FindByName", location.Name).Return((*dto.Location)(nil), 500, assert.AnError).Once()

		statusCode, err := service.Save(location)

		assert.Error(t, err)
		assert.Equal(t, 500, statusCode)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestGetAllLocation(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	service := services.NewLocationService(mockRepo)

	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	t.Run("Success", func(t *testing.T) {
		locations := []*dto.Location{
			{Name: "Warehouse A"},
			{Name: "Warehouse B"},
		}

		mockRepo.On("GetAllLocation", pagination).Return(locations, nil).Once()

		result, statusCode, err := service.GetAll(pagination)

		assert.NoError(t, err)
		assert.Equal(t, 200, statusCode)
		assert.Equal(t, 2, len(result))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Internal Error", func(t *testing.T) {
		mockRepo.On("GetAllLocation", pagination).Return(([]*dto.Location)(nil), assert.AnError).Once()

		result, statusCode, err := service.GetAll(pagination)

		assert.Error(t, err)
		assert.Equal(t, 500, statusCode)
		assert.Nil(t, result)
		mockRepo.AssertExpectations(t)
	})
}
