package services

import (
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/nabilwafi/warehouse-management-system/src/models/dto"
	"github.com/nabilwafi/warehouse-management-system/src/models/web"
	"github.com/nabilwafi/warehouse-management-system/src/repositories"
)

type OrderService interface {
	ReceiveOrder(order *dto.OrderCreateRequest) (int, error)
	ShipOrder(order *dto.OrderCreateRequest) (int, error)
	GetAllOrders(pagination *web.PaginationRequest) ([]*dto.Order, int, error)
	GetOrderByID(id uuid.UUID) (*dto.Order, int, error)
}

type orderServiceImpl struct {
	order       repositories.OrderRepository
	product     repositories.ProductRepository
	transaction repositories.TransactionRepository
	mu          sync.Mutex
}

func NewOrderService(order repositories.OrderRepository, product repositories.ProductRepository, transaction repositories.TransactionRepository) OrderService {
	return &orderServiceImpl{
		order:       order,
		product:     product,
		transaction: transaction,
	}
}

func (s *orderServiceImpl) ReceiveOrder(order *dto.OrderCreateRequest) (int, error) {
	err := s.transaction.Begin()
	if err != nil {
		return 500, errors.New("error when create tx")
	}

	orderData := &dto.Order{
		ProductID: order.ProductID,
		Type:      dto.OrderTypeReceiving,
		Quantity:  order.Quantity,
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	var mu sync.Mutex
	err = s.transaction.Transaction(func() error {
		tx, _ := s.transaction.GetTx()

		wg.Add(1)
		go func(tx *sqlx.Tx, order *dto.Order) {
			mu.Lock()
			defer wg.Done()
			if err := s.order.SaveWithTransaction(tx, order); err != nil {
				errCh <- err
			}
			defer mu.Unlock()
		}(tx, orderData)

		wg.Add(1)
		go func(tx *sqlx.Tx, productID uuid.UUID, quantity int64) {
			mu.Lock()
			defer wg.Done()
			if _, err := s.product.IncreaseStockWithTransaction(tx, productID, quantity); err != nil {
				errCh <- err
			}
			defer mu.Unlock()
		}(tx, orderData.ProductID, orderData.Quantity)

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			return err
		}

		return nil
	})

	if err != nil {
		return 500, err
	}

	return 200, nil
}

func (s *orderServiceImpl) ShipOrder(order *dto.OrderCreateRequest) (int, error) {
	err := s.transaction.Begin()
	if err != nil {
		return 500, errors.New("error when create tx")
	}

	orderData := &dto.Order{
		ProductID: order.ProductID,
		Type:      dto.OrderTypeShipping,
		Quantity:  order.Quantity,
	}

	var wg sync.WaitGroup
	errCh := make(chan error, 2)
	var mu sync.Mutex
	err = s.transaction.Transaction(func() error {
		tx, _ := s.transaction.GetTx()

		wg.Add(1)
		go func(tx *sqlx.Tx, order *dto.Order) {
			mu.Lock()
			defer wg.Done()
			if err := s.order.SaveWithTransaction(tx, order); err != nil {
				errCh <- err
			}
			defer mu.Unlock()
		}(tx, orderData)

		wg.Add(1)
		go func(tx *sqlx.Tx, productID uuid.UUID, quantity int64) {
			mu.Lock()
			defer wg.Done()
			if _, err := s.product.DecreaseStockWithTransaction(tx, productID, quantity); err != nil {
				errCh <- err
			}
			defer mu.Unlock()
		}(tx, orderData.ProductID, orderData.Quantity)

		go func() {
			wg.Wait()
			close(errCh)
		}()

		for err := range errCh {
			return err
		}

		return nil
	})

	if err != nil {
		return 500, err
	}

	return 200, nil
}

func (s *orderServiceImpl) GetAllOrders(pagination *web.PaginationRequest) ([]*dto.Order, int, error) {
	users, err := s.order.FindAll(pagination)
	if err != nil {
		return nil, 500, err
	}

	return users, 200, nil
}

func (s *orderServiceImpl) GetOrderByID(id uuid.UUID) (*dto.Order, int, error) {
	user, code, err := s.order.FindByID(id)
	if err != nil {
		return nil, code, err
	}

	return user, 200, nil
}
