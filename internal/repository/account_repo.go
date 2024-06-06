package repository

import (
	"context"

	"github.com/gopay/internal/models"
)

type AccountRepo interface {
	FindAll(context.Context) ([]models.Account, error)
}

var _ AccountRepo = (*accountRepoImpl)(nil)

type accountRepoImpl struct {
	accounts map[string]models.Account
}

func NewAccountRepo() *accountRepoImpl {
	return &accountRepoImpl{
		accounts: make(map[string]models.Account),
	}
}

func (r *accountRepoImpl) FindAll(_ context.Context) ([]models.Account, error) {
	accs := []models.Account{}

	for _, account := range r.accounts {
		accs = append(accs, account)
	}

	return accs, nil
}
