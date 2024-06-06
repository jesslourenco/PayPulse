package repository

import (
	"context"
	"testing"

	"github.com/gopay/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestFindAll(t *testing.T) {
	type args struct {
		ctx  context.Context
		data map[string]models.Account
	}

	var scenarios = map[string]struct {
		given   args
		want    []models.Account
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
			},

			want: []models.Account{
				{
					AccountId: "0001",
					Name:      "Shankar",
					LastName:  "Nakai",
				},
				{
					AccountId: "0002",
					Name:      "Jessica",
					LastName:  "Lourenco",
				}},
			wantErr: nil,
		},
		"no accounts": {
			given: args{
				ctx:  context.Background(),
				data: map[string]models.Account{},
			},

			want:    []models.Account{},
			wantErr: nil,
		},
	}

	for name, tcase := range scenarios {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			repo := setup(t, tcase.given.data)

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

func setup(_ *testing.T, initialData map[string]models.Account) *accountRepoImpl {
	repo := NewAccountRepo()
	repo.accounts = initialData
	return repo
}
