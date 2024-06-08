package repository

import (
	"context"
	"testing"
	"time"

	"github.com/gopay/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestTransactions_FindAll(t *testing.T) {
	time := time.Now()

	type args struct {
		ctx   context.Context
		accId string
		data  map[string]models.Transaction
	}

	scenarios := map[string]struct {
		given args
		want  []models.Transaction
	}{
		"happy-path": {
			given: args{
				ctx:   context.Background(),
				accId: "0001",
				data: map[string]models.Transaction{
					"1000000": {
						TransactionId: "1000000",
						Owner:         "0001",
						Sender:        "0001",
						Receiver:      "0001",
						CreatedAt:     time,
						Amount:        7000.00,
						IsConsumed:    false,
					},
				},
			},

			want: []models.Transaction{
				{
					TransactionId: "1000000",
					Owner:         "0001",
					Sender:        "0001",
					Receiver:      "0001",
					CreatedAt:     time,
					Amount:        7000.00,
					IsConsumed:    false,
				},
			},
		},

		"no transactions": {
			given: args{
				ctx:   context.Background(),
				accId: "0001",
				data:  map[string]models.Transaction{},
			},
			want: []models.Transaction{},
		},
	}

	for name, tcase := range scenarios {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			repo := setupTransactions(t, tcase.given.data, nil)

			result := repo.FindAll(tcase.given.ctx, tcase.given.accId)
			assert.ElementsMatch(t, tcase.want, result)
		})
	}
}

func TestTransactions_FindOne(t *testing.T) {
	time := time.Now()

	type args struct {
		ctx  context.Context
		id   string
		data map[string]models.Transaction
	}

	scenarios := map[string]struct {
		given   args
		want    models.Transaction
		wantErr error
	}{
		"happy-path": {
			given: args{
				ctx: context.Background(),
				data: map[string]models.Transaction{
					"1000000": {
						TransactionId: "1000000",
						Owner:         "0001",
						Sender:        "0001",
						Receiver:      "0001",
						CreatedAt:     time,
						Amount:        7000.00,
						IsConsumed:    false,
					},
				},
				id: "1000000",
			},

			want: models.Transaction{
				TransactionId: "1000000",
				Owner:         "0001",
				Sender:        "0001",
				Receiver:      "0001",
				CreatedAt:     time,
				Amount:        7000.00,
				IsConsumed:    false,
			},

			wantErr: nil,
		},
		"transaction not found": {
			given: args{
				ctx: context.Background(),
				data: map[string]models.Transaction{
					"1000000": {
						TransactionId: "1000000",
						Owner:         "0001",
						Sender:        "0001",
						Receiver:      "0001",
						CreatedAt:     time,
						Amount:        7000.00,
						IsConsumed:    false,
					},
				},
				id: "2000000",
			},

			want:    models.Transaction{},
			wantErr: ErrTransactionNotFound,
		},
	}

	for name, tcase := range scenarios {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			repo := setupTransactions(t, tcase.given.data, nil)

			result, err := repo.FindOne(tcase.given.ctx, tcase.given.id)

			if tcase.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tcase.wantErr.Error())
			}

			assert.Equal(t, tcase.want, result)
		})
	}
}

func TestTransaction_Create(t *testing.T) {
	time := time.Now()

	id := "0123456789"
	idGenerator := func() string {
		return id
	}

	type args struct {
		ctx         context.Context
		transaction models.Transaction
		data        map[string]models.Transaction
	}

	scenarios := map[string]struct {
		given   args
		want    models.Transaction
		wantErr error
	}{
		"happy-path": {
			given: args{
				ctx: context.Background(),
				transaction: models.Transaction{
					Owner:      "0001",
					Sender:     "0001",
					Receiver:   "0001",
					CreatedAt:  time,
					Amount:     1000,
					IsConsumed: false,
				},
				data: map[string]models.Transaction{},
			},
			want: models.Transaction{
				TransactionId: id,
				Owner:         "0001",
				Sender:        "0001",
				Receiver:      "0001",
				CreatedAt:     time,
				Amount:        1000,
				IsConsumed:    false,
			},
			wantErr: nil,
		},

		"missing owner": {
			given: args{
				ctx: context.Background(),
				transaction: models.Transaction{
					Owner:      "",
					Sender:     "0001",
					Receiver:   "0001",
					CreatedAt:  time,
					Amount:     1000,
					IsConsumed: false,
				},
				data: map[string]models.Transaction{},
			},
			want:    models.Transaction{},
			wantErr: ErrMissingOwnerField,
		},

		"missing sender": {
			given: args{
				ctx: context.Background(),
				transaction: models.Transaction{
					Owner:      "0001",
					Sender:     "",
					Receiver:   "0001",
					CreatedAt:  time,
					Amount:     1000,
					IsConsumed: false,
				},
				data: map[string]models.Transaction{},
			},
			want:    models.Transaction{},
			wantErr: ErrMissingSenderField,
		},

		"missing receiver": {
			given: args{
				ctx: context.Background(),
				transaction: models.Transaction{
					Owner:      "0001",
					Sender:     "0001",
					Receiver:   "",
					CreatedAt:  time,
					Amount:     1000,
					IsConsumed: false,
				},
				data: map[string]models.Transaction{},
			},
			want:    models.Transaction{},
			wantErr: ErrMissingReceiverField,
		},

		"missing owner, receiver, sender": {
			given: args{
				ctx: context.Background(),
				transaction: models.Transaction{
					Owner:      "",
					Sender:     "",
					Receiver:   "",
					CreatedAt:  time,
					Amount:     1000,
					IsConsumed: false,
				},
				data: map[string]models.Transaction{},
			},
			want:    models.Transaction{},
			wantErr: ErrMissingFields,
		},

		"amount is zero": {
			given: args{
				ctx: context.Background(),
				transaction: models.Transaction{
					Owner:      "0001",
					Sender:     "0001",
					Receiver:   "0001",
					CreatedAt:  time,
					Amount:     0,
					IsConsumed: false,
				},
				data: map[string]models.Transaction{},
			},
			want:    models.Transaction{},
			wantErr: ErrZeroAmount,
		},
	}

	for name, tcase := range scenarios {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			repo := setupTransactions(t, tcase.given.data, idGenerator)

			err := repo.Create(tcase.given.ctx, tcase.given.transaction)

			if tcase.wantErr == nil {
				assert.NoError(t, err)
				transaction, err := repo.FindOne(tcase.given.ctx, id)
				assert.NoError(t, err)
				assert.Equal(t, tcase.want, transaction)
			} else {
				assert.ErrorIs(t, err, tcase.wantErr)
			}
		})
	}
}

func setupTransactions(_ *testing.T, initialData map[string]models.Transaction, idGenerator func() string) *transactionRepoImpl {
	repo := NewTransactionRepo()
	repo.transactions = initialData
	repo.idGenerator = idGenerator
	return repo
}
