package repository

import (
	"context"
	"testing"

	"github.com/gopay/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestAccount_Create(t *testing.T) {
	id := "0003"
	idGenerator := func() string {
		return id
	}

	type args struct {
		ctx      context.Context
		name     string
		lastname string
		data     map[string]models.Account
	}

	var scenarios = map[string]struct {
		given   args
		want    models.Account
		wantErr error
	}{
		"happy-path": {
			given: args{
				ctx:      context.Background(),
				name:     "Caio",
				lastname: "Henrique",
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

			want: models.Account{
				AccountId: id,
				Name:      "Caio",
				LastName:  "Henrique",
			},
			wantErr: nil,
		},
		"missing params": {
			given: args{
				ctx:      context.Background(),
				name:     "",
				lastname: "Henrique",
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

			want:    models.Account{},
			wantErr: ErrMissingParams,
		},
	}

	for name, tcase := range scenarios {
		tcase := tcase
		t.Run(name, func(t *testing.T) {
			repo := setup(t, tcase.given.data, idGenerator)

			id, err := repo.Create(tcase.given.ctx, tcase.given.name, tcase.given.lastname)

			if tcase.wantErr == nil {
				assert.NoError(t, err)
				acc, err := repo.FindOne(tcase.given.ctx, id)
				assert.NoError(t, err)
				assert.Equal(t, tcase.want, acc)
			} else {
				assert.EqualError(t, err, tcase.wantErr.Error())
			}

		})
	}
}

func TestAccount_FindAll(t *testing.T) {
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
			repo := setup(t, tcase.given.data, nil)

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

func TestAccount_FindOne(t *testing.T) {
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
			repo := setup(t, tcase.given.data, nil)

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

func setup(_ *testing.T, initialData map[string]models.Account, idGenerator func() string) *accountRepoImpl {
	repo := NewAccountRepo()
	repo.accounts = initialData
	repo.idGenerator = idGenerator
	return repo
}
