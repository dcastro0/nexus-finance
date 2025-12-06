package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"

	"github.com/dcastro0/nexus-finance/configs"
	"github.com/dcastro0/nexus-finance/internal/domain/entity"
	"github.com/dcastro0/nexus-finance/internal/infra/database"
	"github.com/dcastro0/nexus-finance/internal/infra/http/handler"
	"github.com/dcastro0/nexus-finance/internal/infra/repository"
	"github.com/dcastro0/nexus-finance/internal/usecase"
)

func main() {
	conf, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	db, err := database.NewPostgresConnection(&database.Config{
		Host:     conf.DBHost,
		Port:     conf.DBPort,
		User:     conf.DBUser,
		Password: conf.DBPassword,
		DBName:   conf.DBName,
		SSLMode:  "disable",
	})
	if err != nil {
		panic(err)
	}

	runMigrations(db)

	// Repositories
	accountRepo := repository.NewAccountRepositoryPostgres(db)
	transactionRepo := repository.NewTransactionRepositoryPostgres(db)

	// Use Cases
	createAccountUseCase := usecase.NewCreateAccountUseCase(accountRepo)
	makeDepositUseCase := usecase.NewMakeDepositUseCase(accountRepo) // Novo
	makeTransferUseCase := usecase.NewMakeTransferUseCase(transactionRepo, accountRepo)

	// Handlers
	// Atualizado com o novo UseCase de Depósito
	accountHandler := handler.NewAccountHandler(createAccountUseCase, makeDepositUseCase)
	transactionHandler := handler.NewTransactionHandler(makeTransferUseCase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Rotas
	r.Post("/accounts", accountHandler.CreateAccount)
	r.Post("/accounts/{account_id}/deposit", accountHandler.Deposit) // Nova rota
	r.Post("/transactions", transactionHandler.MakeTransfer)

	fmt.Printf("Nexus Finance API running on port %s\n", conf.WebServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", conf.WebServerPort), r)
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate(&entity.Account{}, &entity.Transaction{})
}
