package entity

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrAmountInvalid       = errors.New("amount must be greater than zero")
	ErrSameAccountTransfer = errors.New("cannot transfer to the same account")
)

type Transaction struct {
	ID            uuid.UUID `gorm:"type:uuid;primary_key" json:"id"`
	FromAccountID uuid.UUID `gorm:"type:uuid;not null;index" json:"from_account_id"`
	ToAccountID   uuid.UUID `gorm:"type:uuid;not null;index" json:"to_account_id"`
	Amount        int64     `gorm:"type:bigint;not null" json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
}

func NewTransaction(fromID, toID uuid.UUID, amount int64) (*Transaction, error) {
	if amount <= 0 {
		return nil, ErrAmountInvalid
	}
	if fromID == toID {
		return nil, ErrSameAccountTransfer
	}

	return &Transaction{
		ID:            uuid.New(),
		FromAccountID: fromID,
		ToAccountID:   toID,
		Amount:        amount,
		CreatedAt:     time.Now(),
	}, nil
}
