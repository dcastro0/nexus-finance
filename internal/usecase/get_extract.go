package usecase

import (
	"context"
	"time"

	"github.com/dcastro0/nexus-finance/internal/domain/repository"
)

type GetExtractInputDTO struct {
	AccountID string `json:"account_id"`
	Days      int    `json:"days"`
	Page      int    `json:"page"`
	Limit     int    `json:"limit"`
}

type ExtractTransactionDTO struct {
	ID            string    `json:"id"`
	FromAccountID string    `json:"from_account_id"`
	ToAccountID   string    `json:"to_account_id"`
	Amount        int64     `json:"amount"`
	Type          string    `json:"type"`
	CreatedAt     time.Time `json:"created_at"`
}

type GetExtractOutputDTO struct {
	Transactions []ExtractTransactionDTO `json:"transactions"`
}

type GetExtractUseCase struct {
	TransactionRepository repository.TransactionRepository
}

func NewGetExtractUseCase(repo repository.TransactionRepository) *GetExtractUseCase {
	return &GetExtractUseCase{TransactionRepository: repo}
}

func (uc *GetExtractUseCase) Execute(ctx context.Context, input GetExtractInputDTO) (*GetExtractOutputDTO, error) {
	if input.Days <= 0 {
		input.Days = 7
	}
	if input.Limit <= 0 {
		input.Limit = 10
	}
	if input.Page <= 0 {
		input.Page = 1
	}

	startDate := time.Now().AddDate(0, 0, -input.Days)
	offset := (input.Page - 1) * input.Limit

	transactions, err := uc.TransactionRepository.FindAllByAccount(ctx, input.AccountID, startDate, input.Limit, offset)
	if err != nil {
		return nil, err
	}

	var output []ExtractTransactionDTO
	for _, t := range transactions {
		txType := "DEBIT"
		if t.ToAccountID.String() == input.AccountID {
			txType = "CREDIT"
		}

		output = append(output, ExtractTransactionDTO{
			ID:            t.ID.String(),
			FromAccountID: t.FromAccountID.String(),
			ToAccountID:   t.ToAccountID.String(),
			Amount:        t.Amount,
			Type:          txType,
			CreatedAt:     t.CreatedAt,
		})
	}

	return &GetExtractOutputDTO{Transactions: output}, nil
}
