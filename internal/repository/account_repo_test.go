package repository

import (
	"context"
	"testing"

	"github.com/gopay/internal/models"
)

func TestFindAll(t *testing.T) {
	initialData := map[string]models.Account{
		"0001": {
			AccountId: "0001",
			Name:      "Shankar",
			LastName:  "Nakai",
		},
	}

	t.Run("happy path", func(t *testing.T) {
		repo := setup(t, initialData)

		accounts, err := repo.FindAll(context.Background())
		if err != nil {
			t.Error(err)
		}

		if len(accounts) == 0 {
			t.Error("accounts is empty")
		}
	})
}

func setup(_ *testing.T, initialData map[string]models.Account) *accountRepoImpl {
	repo := NewAccountRepo()
	repo.accounts = initialData
	return repo
}
