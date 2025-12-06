package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dcastro0/nexus-finance/internal/usecase"
	"github.com/go-chi/chi/v5"
)

type AccountHandler struct {
	CreateAccountUseCase *usecase.CreateAccountUseCase
	MakeDepositUseCase   *usecase.MakeDepositUseCase // Novo campo
}

func NewAccountHandler(createUC *usecase.CreateAccountUseCase, depositUC *usecase.MakeDepositUseCase) *AccountHandler {
	return &AccountHandler{
		CreateAccountUseCase: createUC,
		MakeDepositUseCase:   depositUC,
	}
}

// CreateAccount godoc
// @Summary      Criar nova conta
// @Description  Cria uma nova conta bancária com saldo inicial zero
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        request body usecase.CreateAccountInputDTO true "Dados da conta"
// @Success      201  {object}  usecase.CreateAccountOutputDTO
// @Failure      400  {object}  map[string]string
// @Router       /accounts [post]
func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
	// ... (código existente permanece igual)
	var input usecase.CreateAccountInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := h.CreateAccountUseCase.Execute(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}

// Deposit godoc
// @Summary      Realizar Depósito
// @Description  Adiciona fundos a uma conta existente
// @Tags         Accounts
// @Accept       json
// @Produce      json
// @Param        account_id path string true "ID da Conta"
// @Param        request body usecase.MakeDepositInputDTO true "Valor do depósito"
// @Success      200
// @Router       /accounts/{account_id}/deposit [post]
func (h *AccountHandler) Deposit(w http.ResponseWriter, r *http.Request) {
	accountID := chi.URLParam(r, "account_id")
	if accountID == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var input usecase.MakeDepositInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	input.AccountID = accountID

	output, err := h.MakeDepositUseCase.Execute(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
