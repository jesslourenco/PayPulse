package repository

import (
	"context"
	"errors"

	"github.com/gopay/internal/models"
)

var (
	ErrAccountNotFound = errors.New("account not found")
)

type AccountRepo interface {
	FindAll(ctx context.Context) ([]models.Account, error)
	FindOne(ctx context.Context, id string) (models.Account, error)
}

var _ AccountRepo = (*accountRepoImpl)(nil)

type accountRepoImpl struct {
	accounts map[string]models.Account
	account  models.Account
}

func NewAccountRepo() *accountRepoImpl {
	return &accountRepoImpl{
		accounts: make(map[string]models.Account),
		account:  models.Account{},
	}
}

func (r *accountRepoImpl) FindAll(_ context.Context) ([]models.Account, error) {
	accs := []models.Account{}

	for _, account := range r.accounts {
		accs = append(accs, account)
	}

	return accs, nil
}

func (r *accountRepoImpl) FindOne(_ context.Context, id string) (models.Account, error) {
	account, found := r.accounts[id]

	if !found {
		return models.Account{}, ErrAccountNotFound
	}

	return account, nil
}
