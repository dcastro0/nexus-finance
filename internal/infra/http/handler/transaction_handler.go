package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dcastro0/nexus-finance/internal/usecase"
)

type TransactionHandler struct {
	MakeTransferUseCase *usecase.MakeTransferUseCase
	GetExtractUseCase   *usecase.GetExtractUseCase
}

func NewTransactionHandler(makeTransferUC *usecase.MakeTransferUseCase, getExtractUC *usecase.GetExtractUseCase) *TransactionHandler {
	return &TransactionHandler{
		MakeTransferUseCase: makeTransferUC,
		GetExtractUseCase:   getExtractUC,
	}
}

func (h *TransactionHandler) MakeTransfer(w http.ResponseWriter, r *http.Request) {
	var input usecase.MakeTransferInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	accountID, ok := r.Context().Value("account_id").(string)
	if !ok || accountID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	input.FromAccountID = accountID

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

func (h *TransactionHandler) GetExtract(w http.ResponseWriter, r *http.Request) {
	accountID, ok := r.Context().Value("account_id").(string)
	if !ok || accountID == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	days, _ := strconv.Atoi(r.URL.Query().Get("days"))
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	input := usecase.GetExtractInputDTO{
		AccountID: accountID,
		Days:      days,
		Page:      page,
		Limit:     limit,
	}

	output, err := h.GetExtractUseCase.Execute(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
