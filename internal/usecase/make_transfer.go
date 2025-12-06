package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dcastro0/nexus-finance/internal/domain/entity"
	"github.com/dcastro0/nexus-finance/internal/domain/repository"
	"github.com/google/uuid"
)

var ErrNightlyLimitExceeded = errors.New("transaction limit exceeded for nightly period (20h-06h)")

type MakeTransferInputDTO struct {
	FromAccountID string `json:"from_account_id"`
	ToAccountID   string `json:"to_account_id"`
	Amount        int64  `json:"amount"`
}

type MakeTransferOutputDTO struct {
	TransactionID string    `json:"transaction_id"`
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

type MakeTransferUseCase struct {
	TransactionRepository repository.TransactionRepository
	AccountRepository     repository.AccountRepository
	NightlyLimit          int64
}

func NewMakeTransferUseCase(tr repository.TransactionRepository, ar repository.AccountRepository, nightlyLimit int64) *MakeTransferUseCase {
	return &MakeTransferUseCase{
		TransactionRepository: tr,
		AccountRepository:     ar,
		NightlyLimit:          nightlyLimit,
	}
}

func (uc *MakeTransferUseCase) Execute(ctx context.Context, input MakeTransferInputDTO) (*MakeTransferOutputDTO, error) {
	if err := uc.validateNightlyLimit(input.Amount); err != nil {
		return nil, err
	}

	fromID, err := uuid.Parse(input.FromAccountID)
	if err != nil {
		return nil, err
	}
	toID, err := uuid.Parse(input.ToAccountID)
	if err != nil {
		return nil, err
	}

	t, err := entity.NewTransaction(fromID, toID, input.Amount)
	if err != nil {
		return nil, err
	}

	if err := uc.TransactionRepository.Transfer(ctx, t); err != nil {
		return nil, err
	}

	return &MakeTransferOutputDTO{TransactionID: t.ID.String(), Status: "COMPLETED", CreatedAt: t.CreatedAt}, nil
}

func (uc *MakeTransferUseCase) validateNightlyLimit(amount int64) error {
	now := time.Now()
	hour := now.Hour()

	isNightly := hour >= 20 || hour < 6

	if isNightly && amount > uc.NightlyLimit {
		return ErrNightlyLimitExceeded
	}

	return nil
}
