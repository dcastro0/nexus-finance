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
	customMiddleware "github.com/dcastro0/nexus-finance/internal/infra/http/middleware"
	"github.com/dcastro0/nexus-finance/internal/infra/repository"
	"github.com/dcastro0/nexus-finance/internal/usecase"
	"github.com/dcastro0/nexus-finance/pkg/security"
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

	accountRepo := repository.NewAccountRepositoryPostgres(db)
	transactionRepo := repository.NewTransactionRepositoryPostgres(db)

	tokenService := security.NewTokenService(conf.JWTSecret, "nexus-finance-api")

	createAccountUseCase := usecase.NewCreateAccountUseCase(accountRepo)
	loginUseCase := usecase.NewLoginUseCase(accountRepo, tokenService)
	makeDepositUseCase := usecase.NewMakeDepositUseCase(accountRepo)
	makeTransferUseCase := usecase.NewMakeTransferUseCase(transactionRepo, accountRepo, conf.NightlyLimit)
	getExtractUseCase := usecase.NewGetExtractUseCase(transactionRepo)

	accountHandler := handler.NewAccountHandler(createAccountUseCase, makeDepositUseCase)
	authHandler := handler.NewAuthHandler(loginUseCase)
	transactionHandler := handler.NewTransactionHandler(makeTransferUseCase, getExtractUseCase)

	authMiddleware := customMiddleware.NewAuthMiddleware(tokenService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/accounts", accountHandler.CreateAccount)
	r.Post("/login", authHandler.Login)
	r.Post("/accounts/{account_id}/deposit", accountHandler.Deposit)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handle)
		r.Post("/transactions", transactionHandler.MakeTransfer)
		r.Get("/transactions", transactionHandler.GetExtract)
	})

	fmt.Printf("Nexus Finance API running on port %s\n", conf.WebServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", conf.WebServerPort), r)
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate(&entity.Account{}, &entity.Transaction{})
}
