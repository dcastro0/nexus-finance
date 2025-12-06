package repository

import (
	"context"

	"github.com/dcastro0/nexus-finance/internal/domain/entity"
)

type AccountRepository interface {
	Create(ctx context.Context, account *entity.Account) error
	FindByCPF(ctx context.Context, cpf string) (*entity.Account, error)
	FindByID(ctx context.Context, id string) (*entity.Account, error) // Novo método
	UpdateBalance(ctx context.Context, account *entity.Account) error // Novo método
}
