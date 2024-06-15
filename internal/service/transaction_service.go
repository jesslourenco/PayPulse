package service

import (
	"context"
	"errors"
	"time"

	"github.com/gopay/internal/models"
	"github.com/gopay/internal/repository"
)

var (
	ErrInvalidAmount      = errors.New("amount cannot be less or equal to zero")
	ErrInsufficentBalance = errors.New("insufficient balance")
	ErrNoTransactions     = errors.New("account has no transactions")
)

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
	Withdraw(ctx context.Context, owner string, amount float32) error
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

func (r *transactionServiceImpl) Withdraw(ctx context.Context, owner string, amount float32) error {
	_, err := r.accountRepo.FindOne(ctx, owner)
	if err != nil {
		return err
	}

	err = r.debit(ctx, owner, owner, amount)
	if err != nil {
		return err
	}

	transaction := models.Transaction{
		CreatedAt:  clockNow(),
		IsConsumed: true,
		Owner:      owner,
		Sender:     owner,
		Receiver:   owner,
		Amount:     amount,
	}
	return r.transactionRepo.Create(ctx, transaction)
}

func (r *transactionServiceImpl) debit(ctx context.Context, owner string, receiver string, amount float32) error {
	if amount >= 0 {
		return ErrInvalidAmount
	}

	var oldest models.Transaction
	transactions, err := r.transactionRepo.FindAll(ctx, owner)
	if err != nil {
		return err
	}

	if len(transactions) == 0 {
		return ErrNoTransactions
	}

	var balance float64
	for _, t := range transactions {
		if t.Owner == owner && !t.IsConsumed {
			balance += float64(t.Amount)
			if (models.Transaction{} == oldest) || t.CreatedAt.Before(oldest.CreatedAt) {
				oldest = t
			}
		}
	}

	if (balance + float64(amount)) < 0 {
		return ErrInsufficentBalance
	}

	err = r.transactionRepo.MarkAsConsumed(ctx, oldest.TransactionId)
	if err != nil {
		return err
	}

	if (oldest.Amount + amount) != 0 {
		transaction := models.Transaction{
			CreatedAt:  clockNow(),
			IsConsumed: false,
			Owner:      owner,
			Sender:     owner,
			Receiver:   receiver,
			Amount:     oldest.Amount + amount,
		}
		return r.transactionRepo.Create(ctx, transaction)
	}

	return nil
}
