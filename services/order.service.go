package services

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
	"github.com/nabilwafi/warehouse-management-system/models/domain"
	"github.com/nabilwafi/warehouse-management-system/models/web"
	"github.com/nabilwafi/warehouse-management-system/repositories"
)

type orderServiceImpl struct {
	repo        repositories.OrderRepository
	productRepo repositories.ProductRepository
	db          *sql.DB
}

func NewOrderService(repo repositories.OrderRepository, productRepo repositories.ProductRepository, db *sql.DB) OrderService {
	return &orderServiceImpl{repo: repo, db: db, productRepo: productRepo}
}

func (s *orderServiceImpl) ReceiveOrder(ctx context.Context, order *web.OrderCreateRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	go func() {
		err := s.productRepo.IncreaseStock(ctx, tx, order.ProductID, order.Quantity)
		if err != nil {
			return
		}
	}()

	orderData := new(domain.Order)
	orderData.ID = uuid.New()
	orderData.Type = "receiving"
	orderData.ProductID = order.ProductID
	orderData.Quantity = order.Quantity

	if err := s.repo.Save(ctx, tx, orderData); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *orderServiceImpl) ShipOrder(ctx context.Context, order *web.OrderCreateRequest) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	go func() {
		err := s.productRepo.DecreaseStock(ctx, tx, order.ProductID, order.Quantity)
		if err != nil {
			tx.Rollback()
			return
		}
	}()

	orderData := new(domain.Order)
	orderData.ID = uuid.New()
	orderData.Type = "shipping"
	orderData.ProductID = order.ProductID
	orderData.Quantity = order.Quantity
	if err := s.repo.Save(ctx, tx, orderData); err != nil {
		return err
	}

	return tx.Commit()
}

func (s *orderServiceImpl) GetAllOrders(ctx context.Context) ([]domain.Order, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	orders, err := s.repo.FindAll(ctx, tx)
	if err != nil {
		return nil, err
	}

	return orders, nil
}

func (s *orderServiceImpl) GetOrderByID(ctx context.Context, id uuid.UUID) (*domain.Order, error) {
	tx, err := s.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	order, err := s.repo.FindByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return order, nil
}
