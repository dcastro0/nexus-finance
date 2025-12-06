package usecase

import (
	"context"

	"github.com/dcastro0/nexus-finance/internal/domain/repository"
)

type MakeDepositInputDTO struct {
	AccountID string `json:"account_id"`
	Amount    int64  `json:"amount"`
}

type MakeDepositOutputDTO struct {
	NewBalance int64 `json:"new_balance"`
}

type MakeDepositUseCase struct {
	AccountRepository repository.AccountRepository
}

func NewMakeDepositUseCase(accountRepo repository.AccountRepository) *MakeDepositUseCase {
	return &MakeDepositUseCase{
		AccountRepository: accountRepo,
	}
}

func (uc *MakeDepositUseCase) Execute(ctx context.Context, input MakeDepositInputDTO) (*MakeDepositOutputDTO, error) {
	account, err := uc.AccountRepository.FindByID(ctx, input.AccountID)
	if err != nil {
		return nil, err
	}

	// Regra simples: Apenas soma ao saldo
	account.Balance += input.Amount

	// O repositório já trata o Optimistic Locking (versionamento)
	err = uc.AccountRepository.UpdateBalance(ctx, account)
	if err != nil {
		return nil, err
	}

	return &MakeDepositOutputDTO{
		NewBalance: account.Balance,
	}, nil
}
