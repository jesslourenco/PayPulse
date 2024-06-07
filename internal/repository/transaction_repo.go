package repository

import (
	"context"
	"errors"

	"github.com/gopay/internal/models"
	"github.com/gopay/internal/utils"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	//ErrMissingParams   = errors.New("must provide a name and last name for an account")
)

type TransactionRepo interface {
	FindAll(ctx context.Context) ([]models.Transaction, error)
	FindOne(ctx context.Context, id string) (models.Transaction, error)
}

var _ TransactionRepo = (*transactionRepoImpl)(nil)

type transactionRepoImpl struct {
	transactions map[string]models.Transaction
	idGenerator  func() string
}

func NewTransactionRepo() *transactionRepoImpl {
	return &transactionRepoImpl{
		transactions: make(map[string]models.Transaction),
		idGenerator:  utils.GetTransactionUUID,
	}
}

func (r *transactionRepoImpl) FindAll(_ context.Context) ([]models.Transaction, error) {
	transactions := []models.Transaction{}

	for _, t := range r.transactions {
		transactions = append(transactions, t)
	}

	return transactions, nil
}

func (r *transactionRepoImpl) FindOne(_ context.Context, id string) (models.Transaction, error) {
	transaction, found := r.transactions[id]

	if !found {
		return models.Transaction{}, ErrTransactionNotFound
	}

	return transaction, nil
}
