package repositories_test

import (
	"database/sql"
	"testing"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	mocks "github.com/nabilwafi/warehouse-management-system/test/mocks"
	"github.com/stretchr/testify/assert"
)

func TestMockLocationRepositorySave_Success(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	location := &dto.Location{
		ID:       uuid.New(),
		Name:     "Test Location",
		Capacity: 5,
	}

	mockRepo.On("Save", location).Return(nil)

	err := mockRepo.Save(location)

	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "Save", location)
}

func TestMockLocationRepositorySave_Error(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	location := &dto.Location{
		ID:       uuid.New(),
		Name:     "Test Location",
		Capacity: 5,
	}

	mockRepo.On("Save", location).Return(assert.AnError)

	err := mockRepo.Save(location)

	assert.Error(t, err)
	assert.Equal(t, assert.AnError, err)
	mockRepo.AssertCalled(t, "Save", location)
}

func TestMockLocationRepositoryFindByName_Success(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	locationName := "Test Location"

	expectedLocation := &dto.Location{
		ID:       uuid.New(),
		Name:     locationName,
		Capacity: 5,
	}

	mockRepo.On("FindByName", locationName).Return(expectedLocation, 200, nil)

	location, code, err := mockRepo.FindByName(locationName)

	assert.NoError(t, err)
	assert.Equal(t, 200, code)
	assert.Equal(t, expectedLocation, location)
	mockRepo.AssertCalled(t, "FindByName", locationName)
}

func TestMockLocationRepositoryFindByName_Error404(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	locationName := "Test Location"

	mockRepo.On("FindByName", locationName).Return((*dto.Location)(nil), 404, sql.ErrNoRows)

	location, code, err := mockRepo.FindByName(locationName)

	assert.Error(t, err)
	assert.Equal(t, 404, code)
	assert.Equal(t, err, sql.ErrNoRows)
	assert.Nil(t, location)
	mockRepo.AssertCalled(t, "FindByName", locationName)
}

func TestMockLocationRepositoryFindByName_Error500(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	locationName := "Test Location"

	mockRepo.On("FindByName", locationName).Return((*dto.Location)(nil), 500, assert.AnError)

	location, code, err := mockRepo.FindByName(locationName)

	assert.Error(t, err)
	assert.Equal(t, 500, code)
	assert.Equal(t, err, assert.AnError)
	assert.Nil(t, location)
	mockRepo.AssertCalled(t, "FindByName", locationName)
}

func TestMockLocationRepositoryGetAllLocation_Success(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	locations := []*dto.Location{
		{ID: uuid.New(), Name: "Location 1", Capacity: 5},
		{ID: uuid.New(), Name: "Location 2", Capacity: 5},
	}

	mockRepo.On("GetAllLocation", pagination).Return(locations, nil)

	locations, err := mockRepo.GetAllLocation(pagination)

	assert.NotNil(t, locations)
	assert.NoError(t, err)
	assert.Len(t, locations, 2)
	assert.Equal(t, locations[0].Name, "Location 1")
	mockRepo.AssertCalled(t, "GetAllLocation", pagination)
}

func TestMockLocationRepositoryGetAllLocation_SuccessNil(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	mockRepo.On("GetAllLocation", pagination).Return(([]*dto.Location)(nil), nil)

	locations, err := mockRepo.GetAllLocation(pagination)

	assert.Nil(t, locations)
	assert.NoError(t, err)
	mockRepo.AssertCalled(t, "GetAllLocation", pagination)
}

func TestMockLocationRepositoryGetAllLocation_Nil(t *testing.T) {
	mockRepo := new(mocks.MockLocationRepository)
	pagination := &web.PaginationRequest{
		Page: 1,
		Size: 10,
	}

	mockRepo.On("GetAllLocation", pagination).Return(([]*dto.Location)(nil), assert.AnError)

	locations, err := mockRepo.GetAllLocation(pagination)

	assert.Nil(t, locations)
	assert.Error(t, err)
	assert.EqualError(t, err, assert.AnError.Error())
	mockRepo.AssertCalled(t, "GetAllLocation", pagination)
}
