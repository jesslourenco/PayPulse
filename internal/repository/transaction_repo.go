package repository

import (
	"context"
	"errors"

	"github.com/gopay/internal/models"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	//ErrMissingParams   = errors.New("must provide a name and last name for an account")
)

type TransactionRepo interface {
	FindAll(ctx context.Context) ([]models.Transaction, error)
}

var _ TransactionRepo = (*transactionRepoImpl)(nil)

type transactionRepoImpl struct {
	transactions map[string]models.Transaction
	transaction  models.Transaction
}

func NewTransactionRepo() *transactionRepoImpl {
	return &transactionRepoImpl{
		transactions: make(map[string]models.Transaction),
		transaction:  models.Transaction{},
	}
}

func (r *transactionRepoImpl) FindAll(_ context.Context) ([]models.Transaction, error) {
	transactions := []models.Transaction{}

	for _, t := range r.transactions {
		transactions = append(transactions, t)
	}

	return transactions, nil
}

/*func (r *accountRepoImpl) FindOne(_ context.Context, id string) (models.Account, error) {
	account, found := r.accounts[id]

	if !found {
		return models.Account{}, ErrAccountNotFound
	}

	return account, nil
}*/
