package repository

import (
	"context"
	"errors"

	"github.com/gopay/internal/models"
	"github.com/gopay/internal/utils"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrMissingParams   = errors.New("must provide name and last name")
)

type AccountRepo interface {
	FindAll(ctx context.Context) ([]models.Account, error)
	FindOne(ctx context.Context, id string) (models.Account, error)
	Create(ctx context.Context, name string, lastname string) (string, error)
}

var _ AccountRepo = (*accountRepoImpl)(nil)

type accountRepoImpl struct {
	accounts    map[string]models.Account
	idGenerator func() string
}

func NewAccountRepo() *accountRepoImpl {
	return &accountRepoImpl{
		accounts:    make(map[string]models.Account),
		idGenerator: utils.GetAccountUUID,
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

func (r *accountRepoImpl) Create(_ context.Context, name string, lastname string) (string, error) {
	if name == "" || lastname == "" {
		return "", ErrMissingParams
	}

	id := r.idGenerator()

	acc := models.Account{
		AccountId: id,
		Name:      name,
		LastName:  lastname,
	}
	r.accounts[id] = acc

	return id, nil
}
