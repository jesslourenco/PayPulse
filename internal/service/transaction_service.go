package service

import (
	"context"
	"errors"
	"time"

	"github.com/gopay/internal/models"
	"github.com/gopay/internal/repository"
)

var ErrInvalidAmount = errors.New("amount cannot be less or equal to zero")

var nowOriginal = func() time.Time {
	return time.Now()
}
var clockNow = nowOriginal

func setupClock(value time.Time) {
	clockNow = func() time.Time {
		return value
	}
}

func resetClock() {
	clockNow = nowOriginal
}

type TransactionService interface {
	Deposit(ctx context.Context, owner string, amount float32) error
	Balance(ctx context.Context, accId string) (models.Balance, error)
}

var _ TransactionService = (*transactionServiceImpl)(nil)

type transactionServiceImpl struct {
	transactionRepo repository.TransactionRepo
	accountRepo     repository.AccountRepo
}

func NewTransactionService(transactionRepo repository.TransactionRepo, accountRepo repository.AccountRepo) *transactionServiceImpl {
	return &transactionServiceImpl{
		transactionRepo: transactionRepo,
		accountRepo:     accountRepo,
	}
}

func (r *transactionServiceImpl) Deposit(ctx context.Context, owner string, amount float32) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	_, err := r.accountRepo.FindOne(ctx, owner)
	if err != nil {
		return err
	}

	transaction := models.Transaction{
		CreatedAt:  clockNow(),
		IsConsumed: false,
		Owner:      owner,
		Sender:     owner,
		Receiver:   owner,
		Amount:     amount,
	}

	return r.transactionRepo.Create(ctx, transaction)
}

func (r *transactionServiceImpl) Balance(ctx context.Context, accId string) (models.Balance, error) {
	_, err := r.accountRepo.FindOne(ctx, accId)
	if err != nil {
		return models.Balance{}, err
	}

	balance := models.Balance{
		AccountId: accId,
		Amount:    0.00,
	}

	transactions, err := r.transactionRepo.FindAll(ctx, accId)
	if err != nil {
		return models.Balance{}, err
	}

	for _, t := range transactions {
		balance.Amount += float64(t.Amount)
	}

	return balance, nil
}
