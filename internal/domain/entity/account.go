package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Account struct {
	ID        uuid.UUID      `gorm:"type:uuid;primary_key;" json:"id"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	CPF       string         `gorm:"type:varchar(14);uniqueIndex;not null" json:"cpf"`
	Secret    string         `gorm:"type:varchar(255);not null" json:"-"`
	Balance   int64          `gorm:"type:bigint;not null;default:0" json:"balance"`
	Version   int            `gorm:"default:1" json:"version"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func NewAccount(name, cpf, secret string) *Account {
	return &Account{
		ID:      uuid.New(),
		Name:    name,
		CPF:     cpf,
		Secret:  secret,
		Balance: 0,
		Version: 1,
	}
}

func (a *Account) BeforeCreate(tx *gorm.DB) (err error) {
	if a.ID == uuid.Nil {
		a.ID = uuid.New()
	}
	return
}
