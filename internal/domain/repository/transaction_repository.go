package repository

import (
	"context"
	"time"

	"github.com/dcastro0/nexus-finance/internal/domain/entity"
)

type TransactionRepository interface {
	Create(ctx context.Context, transaction *entity.Transaction) error
	Transfer(ctx context.Context, transaction *entity.Transaction) error
	FindAllByAccount(ctx context.Context, accountID string, startDate time.Time, limit, offset int) ([]entity.Transaction, error)
}
