package repository

import (
	"context"
	"errors"
	"fmt"
	"sort"

	"github.com/gopay/internal/models"
	"github.com/gopay/internal/utils"
)

var (
	ErrTransactionNotFound  = errors.New("transaction not found")
	ErrAccountMissing       = errors.New("account id was not provided")
	ErrMissingFields        = errors.New("transaction is missing a field")
	ErrMissingSenderField   = fmt.Errorf("sender: %w", ErrMissingFields)
	ErrMissingReceiverField = fmt.Errorf("receiver: %w", ErrMissingFields)
	ErrMissingOwnerField    = fmt.Errorf("owner: %w", ErrMissingFields)
	ErrZeroAmount           = errors.New("transaction amount cannot be zero")
	ErrNegativeBalance      = errors.New("negative balance")
)

type TransactionRepo interface {
	FindAll(ctx context.Context, accId string) ([]models.Transaction, error)
	FindOne(ctx context.Context, id string) (models.Transaction, error)
	Create(ctx context.Context, transaction models.Transaction) error
	MarkAsConsumed(ctx context.Context, id string) error
	GetBalance(ctx context.Context, id string) (models.Balance, error)
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

func (r *transactionRepoImpl) GetBalance(ctx context.Context, id string) (models.Balance, error) {
	balance := models.Balance{
		AccountId: id,
		Amount:    0.0,
	}

	for _, t := range r.transactions {
		if t.Owner == id && !t.IsConsumed {
			balance.Amount += float64(t.Amount)
		}
	}

	if balance.Amount < 0 {
		return balance, ErrNegativeBalance
	}

	return balance, nil
}

func (r *transactionRepoImpl) FindAll(_ context.Context, accId string) ([]models.Transaction, error) {
	transactions := []models.Transaction{}

	for _, t := range r.transactions {
		if t.Owner == accId {
			transactions = append(transactions, t)
		}
	}

	sort.Slice(transactions, func(i, j int) bool {
		return transactions[i].CreatedAt.Before(transactions[j].CreatedAt)
	})

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
	if transaction.Sender == "" {
		return ErrMissingSenderField
	}
	if transaction.Receiver == "" {
		return ErrMissingReceiverField
	}

	if transaction.Owner == "" {
		return ErrMissingOwnerField
	}

	if transaction.Amount == 0 {
		return ErrZeroAmount
	}

	id := r.idGenerator()
	transaction.TransactionId = id

	r.transactions[id] = transaction

	return nil
}

func (r *transactionRepoImpl) MarkAsConsumed(ctx context.Context, id string) error {
	transaction, err := r.FindOne(ctx, id)
	if err != nil {
		return err
	}

	transaction.IsConsumed = true
	r.transactions[id] = transaction

	return nil
}
