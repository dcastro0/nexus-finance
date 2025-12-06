package repository

import (
	"context"

	"github.com/dcastro0/nexus-finance/internal/domain/entity"
	domainRepo "github.com/dcastro0/nexus-finance/internal/domain/repository"
	"gorm.io/gorm"
)

type AccountRepositoryPostgres struct {
	DB *gorm.DB
}

func NewAccountRepositoryPostgres(db *gorm.DB) *AccountRepositoryPostgres {
	return &AccountRepositoryPostgres{DB: db}
}

func (r *AccountRepositoryPostgres) Create(ctx context.Context, account *entity.Account) error {
	return r.DB.WithContext(ctx).Create(account).Error
}

func (r *AccountRepositoryPostgres) FindByCPF(ctx context.Context, cpf string) (*entity.Account, error) {
	var account entity.Account
	err := r.DB.WithContext(ctx).First(&account, "cpf = ?", cpf).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// CORREÇÃO: Retornando nil para a conta e o erro correto
			return nil, domainRepo.ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepositoryPostgres) FindByID(ctx context.Context, id string) (*entity.Account, error) {
	var account entity.Account
	err := r.DB.WithContext(ctx).Select("id, balance, version").First(&account, "id = ?", id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, domainRepo.ErrAccountNotFound
		}
		return nil, err
	}
	return &account, nil
}

func (r *AccountRepositoryPostgres) UpdateBalance(ctx context.Context, account *entity.Account) error {
	oldVersion := account.Version
	account.Version++ // Incrementa versão para Optimistic Locking

	result := r.DB.WithContext(ctx).Model(&entity.Account{}).
		Where("id = ? AND version = ?", account.ID, oldVersion).
		Updates(map[string]interface{}{
			"balance": account.Balance,
			"version": account.Version,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return domainRepo.ErrConcurrency
	}

	return nil
}
