package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/dcastro0/nexus-finance/internal/domain/entity"
	"github.com/dcastro0/nexus-finance/internal/domain/repository"
	"github.com/dcastro0/nexus-finance/pkg/security"
	"github.com/dcastro0/nexus-finance/pkg/validation" // Importe o novo pacote
)

type CreateAccountInputDTO struct {
	Name   string `json:"name"`
	CPF    string `json:"cpf"`
	Secret string `json:"secret"`
}

type CreateAccountOutputDTO struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CPF       string    `json:"cpf"`
	Balance   int64     `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateAccountUseCase struct {
	AccountRepository repository.AccountRepository
}

func NewCreateAccountUseCase(accountRepository repository.AccountRepository) *CreateAccountUseCase {
	return &CreateAccountUseCase{
		AccountRepository: accountRepository,
	}
}

func (uc *CreateAccountUseCase) Execute(ctx context.Context, input CreateAccountInputDTO) (*CreateAccountOutputDTO, error) {
	// 1. Validação de Domínio (CPF e Senha)
	if !validation.IsCPFValid(input.CPF) {
		return nil, repository.ErrInvalidCPF
	}

	if !validation.IsPasswordStrong(input.Secret) {
		return nil, repository.ErrWeakPassword
	}

	// 2. Verifica existência (Regra de Unicidade)
	_, err := uc.AccountRepository.FindByCPF(ctx, input.CPF)
	if err == nil {
		return nil, errors.New("account already exists for this cpf")
	}
	if err != repository.ErrAccountNotFound {
		return nil, err
	}

	// 3. Hash da Senha e Criação da Entidade
	hashedSecret, err := security.HashPassword(input.Secret)
	if err != nil {
		return nil, err
	}

	account := entity.NewAccount(input.Name, input.CPF, hashedSecret)

	// 4. Persistência
	err = uc.AccountRepository.Create(ctx, account)
	if err != nil {
		return nil, err
	}

	return &CreateAccountOutputDTO{
		ID:        account.ID.String(),
		Name:      account.Name,
		CPF:       account.CPF,
		Balance:   account.Balance,
		CreatedAt: account.CreatedAt,
	}, nil
}
