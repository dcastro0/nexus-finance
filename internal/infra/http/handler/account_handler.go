package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dcastro0/nexus-finance/internal/usecase"
)

type AccountHandler struct {
	CreateAccountUseCase *usecase.CreateAccountUseCase
}

func NewAccountHandler(createAccountUseCase *usecase.CreateAccountUseCase) *AccountHandler {
	return &AccountHandler{
		CreateAccountUseCase: createAccountUseCase,
	}
}

func (h *AccountHandler) CreateAccount(w http.ResponseWriter, r *http.Request) {
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
