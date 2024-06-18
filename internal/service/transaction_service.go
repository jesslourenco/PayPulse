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

	balance, err := r.transactionRepo.GetBalance(ctx, owner)
	if err != nil {
		return err
	}

	if (balance.Amount + float64(amount)) < 0 {
		return ErrInsufficentBalance
	}

	transactions, err := r.transactionRepo.FindAll(ctx, owner)
	if err != nil {
		return err
	}

	debit := (-1) * amount
	for _, t := range transactions {
		err = r.transactionRepo.MarkAsConsumed(ctx, t.TransactionId)
		if err != nil {
			return err
		}

		remaining := t.Amount - debit

		if remaining == 0 {
			break
		}

		if remaining < 0 {
			debit = debit - t.Amount
			continue
		}

		if remaining > 0 {
			transaction := models.Transaction{
				CreatedAt:  clockNow(),
				IsConsumed: false,
				Owner:      owner,
				Sender:     owner,
				Receiver:   receiver,
				Amount:     t.Amount - debit,
			}
			err = r.transactionRepo.Create(ctx, transaction)
			if err != nil {
				return err
			}

			break
		}
	}

	return nil
}
