package repository

import (
	"context"
	"errors"

	"github.com/dcastro0/nexus-finance/internal/domain/entity"
)

var (
	ErrAccountNotFound = errors.New("account not found")
	ErrConcurrency     = errors.New("optimistic locking failure: account modified by another transaction")
)

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	FindByCPF(ctx context.Context, cpf string) (*entity.Account, error)
	FindByID(ctx context.Context, id string) (*entity.Account, error)
	UpdateBalance(ctx context.Context, account *entity.Account) error
}
