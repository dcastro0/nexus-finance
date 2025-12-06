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

// MakeTransfer godoc
// @Summary      Transferência P2P
// @Description  Realiza transferência entre contas internas. Requer Token JWT.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request body usecase.MakeTransferInputDTO true "Dados da transferência (from_account_id é ignorado, usa-se o do token)"
// @Success      201  {object}  usecase.MakeTransferOutputDTO
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /transactions [post]
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

// GetExtract godoc
// @Summary      Extrato Bancário
// @Description  Retorna o histórico de transações com paginação. Requer Token JWT.
// @Tags         Transactions
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        days   query int false "Dias de histórico (default 7)"
// @Param        page   query int false "Número da página (default 1)"
// @Param        limit  query int false "Itens por página (default 10)"
// @Success      200  {object}  usecase.GetExtractOutputDTO
// @Router       /transactions [get]
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
