package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/exception"
	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/repositories"
)

type productServiceImpl struct {
	repo      repositories.ProductRepository
	validator *validator.Validate
	db        *sql.DB
}

func NewProductService(repo repositories.ProductRepository, validator *validator.Validate, db *sql.DB) ProductService {
	return &productServiceImpl{
		repo:      repo,
		validator: validator,
		db:        db,
	}
}

func (s *productServiceImpl) Create(ctx context.Context, product *web.ProductCreateRequest) error {
	if err := s.validator.Struct(product); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.Save(ctx, tx, product); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return exception.NewCustomError(500, "failed to commit transaction: "+err.Error())
	}

	return nil
}

func (s *productServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	if err := s.validator.Var(id, "required,uuid"); err != nil {
		return nil, errors.New("validation failed: " + err.Error())
	}

	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	product, err := s.repo.FindByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, exception.NewCustomError(500, "failed to commit transaction: "+err.Error())
	}

	return product, nil
}

func (s *productServiceImpl) GetAll(ctx context.Context) ([]domain.Product, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	products, err := s.repo.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, exception.NewCustomError(500, "failed to commit transaction: "+err.Error())
	}

	return products, nil
}

func (s *productServiceImpl) Update(ctx context.Context, product *domain.Product) error {
	fmt.Println(product, "a")

	if err := s.validator.Struct(product); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.Update(ctx, tx, product); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return exception.NewCustomError(500, "failed to commit transaction: "+err.Error())
	}

	return nil
}

func (s *productServiceImpl) Delete(ctx context.Context, id uuid.UUID) error {
	if err := s.validator.VarCtx(ctx, id, "required,uuid"); err != nil {
		return errors.New("validation failed: " + err.Error())
	}

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := s.repo.Delete(ctx, tx, id); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return exception.NewCustomError(500, "failed to commit transaction: "+err.Error())
	}

	return nil
}
