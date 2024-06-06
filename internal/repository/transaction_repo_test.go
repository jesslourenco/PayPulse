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
			repo := setupTransactions(t, tcase.given.data)

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

/*func TestFindOne(t *testing.T) {
	type args struct {
		ctx  context.Context
		id   string
		data map[string]models.Account
	}

	var scenarios = map[string]struct {
		given   args
		want    models.Account
		wantErr error
	}{
		"happy-path": {
			given: args{
				ctx: context.Background(),
				data: map[string]models.Account{
					"0001": {
						AccountId: "0001",
						Name:      "Shankar",
						LastName:  "Nakai",
					},
					"0002": {
						AccountId: "0002",
						Name:      "Jessica",
						LastName:  "Lourenco",
					}},
				id: "0002",
			},

			want: models.Account{
				AccountId: "0002",
				Name:      "Jessica",
				LastName:  "Lourenco",
			},

			wantErr: nil,
		},
		"account not found": {
			given: args{
				ctx: context.Background(),
				data: map[string]models.Account{
					"0001": {
						AccountId: "0001",
						Name:      "Shankar",
						LastName:  "Nakai",
					},
					"0002": {
						AccountId: "0002",
						Name:      "Jessica",
						LastName:  "Lourenco",
					}},
				id: "0003",
			},

			want:    models.Account{},
			wantErr: ErrAccountNotFound,
		},
	}

	for name, tcase := range scenarios {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			repo := setup(t, tcase.given.data)

			result, err := repo.FindOne(tcase.given.ctx, tcase.given.id)

			if tcase.wantErr == nil {
				assert.NoError(t, err)
			} else {
				assert.EqualError(t, err, tcase.wantErr.Error())
			}

			assert.Equal(t, tcase.want, result)

		})
	}

}*/

func setupTransactions(_ *testing.T, initialData map[string]models.Transaction) *transactionRepoImpl {
	repo := NewTransactionRepo()
	repo.transactions = initialData
	return repo
}
