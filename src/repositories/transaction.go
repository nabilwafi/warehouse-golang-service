package repositories

import (
	"errors"

	"github.com/jmoiron/sqlx"
)

type TransactionRepository interface {
	Begin() error
	Commit() error
	Rollback() error
	Transaction(callback func() error) error
	GetTx() (*sqlx.Tx, error)
}

type transactionRepositoryImpl struct {
	db *sqlx.DB
	tx *sqlx.Tx
}

func NewTransactionRepository(db *sqlx.DB) TransactionRepository {
	return &transactionRepositoryImpl{db: db}
}

func (r *transactionRepositoryImpl) Begin() (err error) {
	r.tx, err = r.db.Beginx()
	if err != nil {
		err = errors.New("can't create transcation")
		return
	}

	return err
}

func (r *transactionRepositoryImpl) Commit() error {
	err := r.tx.Commit()

	if err != nil {
		return err
	}

	return nil
}

func (r *transactionRepositoryImpl) Rollback() error {
	err := r.tx.Rollback()

	if err != nil {
		return err
	}

	return nil
}

func (r *transactionRepositoryImpl) Transaction(callback func() error) error {
	var err error

	defer func() {
		if p := recover(); p != nil {
			err = r.Rollback()
		} else if err != nil {
			err = r.Rollback()
		} else {
			r.Commit()
		}
	}()

	err = callback()

	return err
}

func (r *transactionRepositoryImpl) GetTx() (*sqlx.Tx, error) {
	if r.tx == nil {
		return nil, errors.New("tx not created")
	}

	return r.tx, nil
}
