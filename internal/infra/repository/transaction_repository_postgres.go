package repository

import (
	"context"

	"github.com/dcastro0/nexus-finance/internal/domain/entity"
	domainRepo "github.com/dcastro0/nexus-finance/internal/domain/repository"
	"gorm.io/gorm"
)

type TransactionRepositoryPostgres struct {
	DB *gorm.DB
}

func NewTransactionRepositoryPostgres(db *gorm.DB) *TransactionRepositoryPostgres {
	return &TransactionRepositoryPostgres{DB: db}
}

func (r *TransactionRepositoryPostgres) Create(ctx context.Context, transaction *entity.Transaction) error {
	return r.DB.WithContext(ctx).Create(transaction).Error
}

func (r *TransactionRepositoryPostgres) Transfer(ctx context.Context, t *entity.Transaction) error {
	return r.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var fromAccount entity.Account
		// Verifica conta origem
		if err := tx.First(&fromAccount, "id = ?", t.FromAccountID).Error; err != nil {
			return err
		}

		if fromAccount.Balance < t.Amount {
			return domainRepo.ErrInsufficientFunds
		}

		// Verifica conta destino
		var toAccount entity.Account
		if err := tx.First(&toAccount, "id = ?", t.ToAccountID).Error; err != nil {
			return err
		}

		// Débito Origem
		oldVersionFrom := fromAccount.Version
		fromAccount.Balance -= t.Amount
		fromAccount.Version++
		resFrom := tx.Model(&entity.Account{}).
			Where("id = ? AND version = ?", fromAccount.ID, oldVersionFrom).
			Updates(map[string]interface{}{"balance": fromAccount.Balance, "version": fromAccount.Version})

		if resFrom.RowsAffected == 0 {
			return domainRepo.ErrConcurrency
		}

		// Crédito Destino
		oldVersionTo := toAccount.Version
		toAccount.Balance += t.Amount
		toAccount.Version++
		resTo := tx.Model(&entity.Account{}).
			Where("id = ? AND version = ?", toAccount.ID, oldVersionTo).
			Updates(map[string]interface{}{"balance": toAccount.Balance, "version": toAccount.Version})

		if resTo.RowsAffected == 0 {
			return domainRepo.ErrConcurrency
		}

		return tx.Create(t).Error
	})
}
