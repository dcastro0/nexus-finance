package usecase

import (
	"context"

	"github.com/dcastro0/nexus-finance/internal/domain/repository"
)

type GetBalanceOutputDTO struct {
	Balance int64 `json:"balance"`
}

type GetBalanceUseCase struct {
	AccountRepository repository.AccountRepository
}

func NewGetBalanceUseCase(repo repository.AccountRepository) *GetBalanceUseCase {
	return &GetBalanceUseCase{AccountRepository: repo}
}

func (uc *GetBalanceUseCase) Execute(ctx context.Context, accountID string) (*GetBalanceOutputDTO, error) {
	account, err := uc.AccountRepository.FindByID(ctx, accountID)
	if err != nil {
		return nil, err
	}

	return &GetBalanceOutputDTO{
		Balance: account.Balance,
	}, nil
}
