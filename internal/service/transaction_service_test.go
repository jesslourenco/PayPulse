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

func TestTransactionService_Balance(t *testing.T) {
	now := time.Now()
	setupClock(now)
	defer resetClock()

	var (
		ctx   = context.Background()
		accId = "0001"
	)

	type args struct {
		accId string
	}

	scenarios := map[string]struct {
		given   args
		doMocks func(deps transactionServiceDependencies)
		want    models.Balance
		wantErr error
	}{
		"happy-path": {
			given: args{
				accId: accId,
			},
			doMocks: func(deps transactionServiceDependencies) {
				deps.accRepoMock.On("FindOne", ctx, accId).Return(models.Account{
					AccountId: accId,
					Name:      "Shankar",
					LastName:  "Nakai",
				}, nil)
				deps.transRepoMock.On("FindAll", ctx, accId).Return([]models.Transaction{
					{
						TransactionId: "1000000",
						Owner:         accId,
						Sender:        accId,
						Receiver:      accId,
						CreatedAt:     time.Now(),
						Amount:        7000.00,
						IsConsumed:    false,
					},
					{
						TransactionId: "2000000",
						Owner:         accId,
						Sender:        "0002",
						Receiver:      accId,
						CreatedAt:     time.Now(),
						Amount:        3000.00,
						IsConsumed:    false,
					},
				}, nil)
			},
			want: models.Balance{
				AccountId: accId,
				Amount:    10000.00,
			},
			wantErr: nil,
		},
		"account-not-found": {
			given: args{
				accId: accId,
			},
			doMocks: func(deps transactionServiceDependencies) {
				deps.accRepoMock.On("FindOne", ctx, accId).Return(models.Account{}, repository.ErrAccountNotFound)
			},
			want:    models.Balance{},
			wantErr: repository.ErrAccountNotFound,
		},
		"no-transactions": {
			given: args{
				accId: accId,
			},
			doMocks: func(deps transactionServiceDependencies) {
				deps.accRepoMock.On("FindOne", ctx, accId).Return(models.Account{
					AccountId: accId,
					Name:      "Shankar",
					LastName:  "Nakai",
				}, nil)
				deps.transRepoMock.On("FindAll", ctx, accId).Return([]models.Transaction{}, nil)
			},
			want: models.Balance{
				AccountId: accId,
				Amount:    0.0,
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

			balance, err := service.Balance(ctx, tcase.given.accId)

			if tcase.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.ErrorIs(t, err, tcase.wantErr)
			}

			assert.Equal(t, tcase.want, balance)
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
