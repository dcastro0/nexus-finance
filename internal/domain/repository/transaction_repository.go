package repository

import (
	"context"
	"errors"

	"github.com/dcastro0/nexus-finance/internal/domain/entity"
)

var (
	ErrInsufficientFunds = errors.New("insufficient funds")
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	Transfer(ctx context.Context, transaction *entity.Transaction) error
}
