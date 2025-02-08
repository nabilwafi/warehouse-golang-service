package services

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/repositories"
)

type ProductService interface {
	Create(product *dto.Product) (int, error)
	GetByID(id uuid.UUID) (*dto.Product, int, error)
	GetAll(pagination *web.PaginationRequest) ([]*dto.Product, int, error)
	Update(product *dto.Product) (int, error)
	Delete(id uuid.UUID) (int, error)
}

type productServiceImpl struct {
	product repositories.ProductRepository
}

func NewProductService(product repositories.ProductRepository) ProductService {
	return &productServiceImpl{
		product: product,
	}
}

func (s *productServiceImpl) Create(product *dto.Product) (int, error) {
	productData, code, err := s.product.FindByName(product.Name)
	if err != nil {
		if err != sql.ErrNoRows || code == 500 {
			return code, err
		}
	}

	if productData != nil {
		return 401, errors.New("product name is exists")
	}

	productData, code, err = s.product.FindBySKU(product.SKU)
	if err != nil {
		if err != sql.ErrNoRows || code == 500 {
			return code, err
		}
	}

	if productData != nil {
		return 401, errors.New("product sku is exists")
	}

	err = s.product.Save(product)
	if err != nil {
		return 500, err
	}

	return 201, nil
}

func (s *productServiceImpl) GetByID(id uuid.UUID) (*dto.Product, int, error) {
	product, code, err := s.product.FindByID(id)
	if err != nil {
		return nil, code, err
	}

	return product, 200, nil
}

func (s *productServiceImpl) GetAll(pagination *web.PaginationRequest) ([]*dto.Product, int, error) {
	products, err := s.product.GetAllProduct(pagination)
	if err != nil {
		return nil, 500, err
	}

	return products, 200, nil
}

func (s *productServiceImpl) Update(product *dto.Product) (int, error) {
	_, code, err := s.product.FindByID(product.ID)
	if err != nil {
		return code, err
	}

	err = s.product.Update(product)
	if err != nil {
		return 500, err
	}

	return 200, nil
}

func (s *productServiceImpl) Delete(id uuid.UUID) (int, error) {
	_, code, err := s.product.FindByID(id)
	if err != nil {
		return code, err
	}

	code, err = s.product.Delete(id)
	if err != nil {
		return code, err
	}

	return 200, nil
}
