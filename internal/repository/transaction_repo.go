package repository

import (
	"context"
	"errors"

	"github.com/gopay/internal/models"
	"github.com/gopay/internal/utils"
)

var (
	ErrTransactionNotFound = errors.New("transaction not found")
	ErrMissingFields       = errors.New("transaction is missing one or more of these: owner, sender, receiver, amount")
	ErrZeroAmount          = errors.New("transaction amount cannot be zero")
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

func (r *transactionRepoImpl) Create(_ context.Context, transaction models.Transaction) error {
	if transaction.Sender == "" || transaction.Receiver == "" || transaction.Owner == "" {
		return ErrMissingFields
	}

	if transaction.Amount == 0 {
		return ErrZeroAmount
	}

	id := r.idGenerator()
	transaction.TransactionId = id

	return nil
}
