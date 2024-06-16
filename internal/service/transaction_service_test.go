package service

import (
	"context"
	"testing"
	"time"

	"github.com/gopay/internal/models"
	"github.com/gopay/internal/repository"
	"github.com/stretchr/testify/assert"
)

func TestTransactionService_Deposit(t *testing.T) {
	now := time.Now()
	setupClock(now)
	defer resetClock()

	var (
		ctx            = context.Background()
		owner          = "0001"
		amount float32 = 50.0
	)

	type args struct {
		owner  string
		amount float32
	}

	scenarios := map[string]struct {
		given   args
		doMocks func(deps transactionServiceDependencies)
		want    models.Transaction
		wantErr error
	}{
		"happy-path": {
			given: args{
				owner:  owner,
				amount: amount,
			},
			doMocks: func(deps transactionServiceDependencies) {
				transaction := models.Transaction{
					CreatedAt:  now,
					IsConsumed: false,
					Owner:      owner,
					Sender:     owner,
					Receiver:   owner,
					Amount:     amount,
				}

				deps.accRepoMock.On("FindOne", ctx, owner).Return(models.Account{
					AccountId: owner,
					Name:      "Shankar",
					LastName:  "Nakai",
				}, nil)
				deps.transRepoMock.On("Create", ctx, transaction).Return(nil)
			},
			wantErr: nil,
		},
		"zero-amount": {
			given: args{
				owner:  owner,
				amount: 0.00,
			},
			wantErr: ErrInvalidAmount,
		},
		"negative-amount": {
			given: args{
				owner:  owner,
				amount: -1.00,
			},
			wantErr: ErrInvalidAmount,
		},
		"invalid-owner": {
			given: args{
				owner:  owner,
				amount: amount,
			},
			doMocks: func(deps transactionServiceDependencies) {
				deps.accRepoMock.On("FindOne", ctx, owner).Return(models.Account{}, repository.ErrAccountNotFound)
			},
			wantErr: repository.ErrAccountNotFound,
		},
	}

	for name, tcase := range scenarios {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			service, deps := setupTransactionService(t)
			if tcase.doMocks != nil {
				tcase.doMocks(deps)
			}

			err := service.Deposit(ctx, tcase.given.owner, tcase.given.amount)

			if tcase.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tcase.wantErr)
			}
		})
	}
}

func TestTransactionService_Withdraw(t *testing.T) {
	now := time.Now()
	setupClock(now)
	defer resetClock()

	var (
		ctx            = context.Background()
		owner          = "0001"
		amount float32 = -1000.0
	)

	type args struct {
		owner  string
		amount float32
	}

	scenarios := map[string]struct {
		given   args
		doMocks func(deps transactionServiceDependencies)
		wantErr error
	}{
		"happy-path": {
			given: args{
				owner:  owner,
				amount: amount,
			},
			doMocks: func(deps transactionServiceDependencies) {
				transactions := []models.Transaction{
					{
						TransactionId: "1000000",
						CreatedAt:     now,
						IsConsumed:    false,
						Owner:         owner,
						Sender:        owner,
						Receiver:      owner,
						Amount:        7000.0,
					},
				}

				debitTransaction := models.Transaction{
					CreatedAt:  now,
					IsConsumed: true,
					Owner:      owner,
					Sender:     owner,
					Receiver:   owner,
					Amount:     amount,
				}

				transaction := models.Transaction{
					CreatedAt:  now,
					IsConsumed: false,
					Owner:      owner,
					Sender:     owner,
					Receiver:   owner,
					Amount:     7000 + amount,
				}

				deps.accRepoMock.On("FindOne", ctx, owner).Return(models.Account{
					AccountId: owner,
					Name:      "Shankar",
					LastName:  "Nakai",
				}, nil)

				deps.transRepoMock.On("GetBalance", ctx, owner).Return(models.Balance{
					AccountId: owner,
					Amount:    7000,
				}, nil)
				deps.transRepoMock.On("FindAll", ctx, owner).Return(transactions, nil)
				deps.transRepoMock.On("MarkAsConsumed", ctx, transactions[0].TransactionId).Return(nil)
				deps.transRepoMock.On("Create", ctx, transaction).Return(nil)
				deps.transRepoMock.On("Create", ctx, debitTransaction).Return(nil)
			},
			wantErr: nil,
		},
		"invalid-owner": {
			given: args{
				owner:  owner,
				amount: amount,
			},
			doMocks: func(deps transactionServiceDependencies) {
				deps.accRepoMock.On("FindOne", ctx, owner).Return(models.Account{}, repository.ErrAccountNotFound)
			},
			wantErr: repository.ErrAccountNotFound,
		},
		"invalid-amount": {
			given: args{
				owner:  owner,
				amount: amount * -1,
			},
			doMocks: func(deps transactionServiceDependencies) {
				deps.accRepoMock.On("FindOne", ctx, owner).Return(models.Account{
					AccountId: owner,
					Name:      "Shankar",
					LastName:  "Nakai",
				}, nil)
			},
			wantErr: ErrInvalidAmount,
		},
		"insufficient-balance": {
			given: args{
				owner:  owner,
				amount: amount,
			},
			doMocks: func(deps transactionServiceDependencies) {
				deps.accRepoMock.On("FindOne", ctx, owner).Return(models.Account{
					AccountId: owner,
					Name:      "Shankar",
					LastName:  "Nakai",
				}, nil)

				deps.transRepoMock.On("GetBalance", ctx, owner).Return(models.Balance{
					AccountId: owner,
					Amount:    500,
				}, nil)
			},
			wantErr: ErrInsufficentBalance,
		},
		"multi-transaction-consumption": {
			given: args{
				owner:  owner,
				amount: -400.0,
			},
			doMocks: func(deps transactionServiceDependencies) {
				transactions := []models.Transaction{
					{
						TransactionId: "1000000",
						CreatedAt:     now,
						IsConsumed:    false,
						Owner:         owner,
						Sender:        owner,
						Receiver:      owner,
						Amount:        200.0,
					},
					{
						TransactionId: "2000000",
						CreatedAt:     now.Add(10),
						IsConsumed:    false,
						Owner:         owner,
						Sender:        owner,
						Receiver:      owner,
						Amount:        100.0,
					},
					{
						TransactionId: "3000000",
						CreatedAt:     now.Add(50),
						IsConsumed:    false,
						Owner:         owner,
						Sender:        owner,
						Receiver:      owner,
						Amount:        300.0,
					},
				}

				debitTransaction := models.Transaction{
					CreatedAt:  now,
					IsConsumed: true,
					Owner:      owner,
					Sender:     owner,
					Receiver:   owner,
					Amount:     -400,
				}

				transaction := models.Transaction{
					CreatedAt:  now,
					IsConsumed: false,
					Owner:      owner,
					Sender:     owner,
					Receiver:   owner,
					Amount:     200,
				}

				deps.accRepoMock.On("FindOne", ctx, owner).Return(models.Account{
					AccountId: owner,
					Name:      "Shankar",
					LastName:  "Nakai",
				}, nil)

				deps.transRepoMock.On("GetBalance", ctx, owner).Return(models.Balance{
					AccountId: owner,
					Amount:    600,
				}, nil)
				deps.transRepoMock.On("FindAll", ctx, owner).Return(transactions, nil)
				deps.transRepoMock.On("MarkAsConsumed", ctx, transactions[0].TransactionId).Return(nil)
				deps.transRepoMock.On("MarkAsConsumed", ctx, transactions[1].TransactionId).Return(nil)
				deps.transRepoMock.On("MarkAsConsumed", ctx, transactions[2].TransactionId).Return(nil)
				deps.transRepoMock.On("Create", ctx, transaction).Return(nil)
				deps.transRepoMock.On("Create", ctx, debitTransaction).Return(nil)
			},
			wantErr: nil,
		},
	}

	for name, tcase := range scenarios {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			service, deps := setupTransactionService(t)
			if tcase.doMocks != nil {
				tcase.doMocks(deps)
			}

			err := service.Withdraw(ctx, tcase.given.owner, tcase.given.amount)

			if tcase.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tcase.wantErr)
			}
		})
	}
}

type transactionServiceDependencies struct {
	transRepoMock *repository.MockTransactionRepo
	accRepoMock   *repository.MockAccountRepo
}

func setupTransactionService(t *testing.T) (*transactionServiceImpl, transactionServiceDependencies) {
	deps := transactionServiceDependencies{
		transRepoMock: repository.NewMockTransactionRepo(t),
		accRepoMock:   repository.NewMockAccountRepo(t),
	}

	return NewTransactionService(deps.transRepoMock, deps.accRepoMock), deps
}
