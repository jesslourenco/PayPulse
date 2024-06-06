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
		ctx  context.Context
		data map[string]models.Transaction
	}

	var scenarios = map[string]struct {
		given   args
		want    []models.Transaction
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
			wantErr: nil,
		},

		"no transactions": {
			given: args{
				ctx:  context.Background(),
				data: map[string]models.Transaction{},
			},
			want:    []models.Transaction{},
			wantErr: nil,
		},
	}

	for name, tcase := range scenarios {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			repo := setupTransactions(t, tcase.given.data, nil)

			result, err := repo.FindAll(tcase.given.ctx)

			if tcase.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tcase.wantErr.Error())
			}

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

	var scenarios = map[string]struct {
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

func setupTransactions(_ *testing.T, initialData map[string]models.Transaction, idGenerator func() string) *transactionRepoImpl {
	repo := NewTransactionRepo()
	repo.transactions = initialData
	repo.idGenerator = idGenerator
	return repo
}
