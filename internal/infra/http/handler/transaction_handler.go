package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dcastro0/nexus-finance/internal/usecase"
)

type TransactionHandler struct {
	MakeTransferUseCase *usecase.MakeTransferUseCase
}

func NewTransactionHandler(uc *usecase.MakeTransferUseCase) *TransactionHandler {
	return &TransactionHandler{MakeTransferUseCase: uc}
}

func (h *TransactionHandler) MakeTransfer(w http.ResponseWriter, r *http.Request) {
	var input usecase.MakeTransferInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	output, err := h.MakeTransferUseCase.Execute(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(output)
}
