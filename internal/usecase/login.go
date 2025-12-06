package usecase

import (
	"context"
	"errors"

	"github.com/dcastro0/nexus-finance/internal/domain/repository"
	"github.com/dcastro0/nexus-finance/pkg/security"
)

type LoginInputDTO struct {
	CPF      string `json:"cpf"`
	Password string `json:"password"`
}

type LoginOutputDTO struct {
	AccessToken string `json:"access_token"`
}

type LoginUseCase struct {
	AccountRepository repository.AccountRepository
	TokenService      *security.TokenService
}

func NewLoginUseCase(repo repository.AccountRepository, tokenService *security.TokenService) *LoginUseCase {
	return &LoginUseCase{
		AccountRepository: repo,
		TokenService:      tokenService,
	}
}

func (uc *LoginUseCase) Execute(ctx context.Context, input LoginInputDTO) (*LoginOutputDTO, error) {
	account, err := uc.AccountRepository.FindByCPF(ctx, input.CPF)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !security.CheckPasswordHash(input.Password, account.Secret) {
		return nil, errors.New("invalid credentials")
	}

	token, err := uc.TokenService.GenerateToken(account.ID.String())
	if err != nil {
		return nil, err
	}

	return &LoginOutputDTO{AccessToken: token}, nil
}
