package handler

import (
	"encoding/json"
	"net/http"

	"github.com/dcastro0/nexus-finance/internal/usecase"
)

type AuthHandler struct {
	LoginUseCase *usecase.LoginUseCase
}

func NewAuthHandler(loginUseCase *usecase.LoginUseCase) *AuthHandler {
	return &AuthHandler{LoginUseCase: loginUseCase}
}

// Login godoc
// @Summary      Autenticação de Usuário
// @Description  Autentica o usuário via CPF e Senha e retorna um token JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        request body usecase.LoginInputDTO true "Credenciais de acesso"
// @Success      200  {object}  usecase.LoginOutputDTO
// @Failure      401  {object}  map[string]string
// @Router       /login [post]
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var input usecase.LoginInputDTO
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	output, err := h.LoginUseCase.Execute(r.Context(), input)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(output)
}
