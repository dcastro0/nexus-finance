package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/gorm"

	"github.com/dcastro0/nexus-finance/configs"
	_ "github.com/dcastro0/nexus-finance/docs" // Import da documentação Swagger
	"github.com/dcastro0/nexus-finance/internal/domain/entity"
	"github.com/dcastro0/nexus-finance/internal/infra/database"
	"github.com/dcastro0/nexus-finance/internal/infra/http/handler"
	customMiddleware "github.com/dcastro0/nexus-finance/internal/infra/http/middleware"
	"github.com/dcastro0/nexus-finance/internal/infra/repository"
	"github.com/dcastro0/nexus-finance/internal/usecase"
	"github.com/dcastro0/nexus-finance/pkg/security"
)

// @title           Nexus Finance API
// @version         1.0
// @description     API de serviços bancários do Nexus Finance Ecosystem (Core Banking).
// @termsOfService  http://swagger.io/terms/

// @contact.name    Suporte Nexus
// @contact.email   suporte@nexus.finance

// @host            localhost:8080
// @BasePath        /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
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
	getBalanceUseCase := usecase.NewGetBalanceUseCase(accountRepo) // Novo UseCase

	accountHandler := handler.NewAccountHandler(createAccountUseCase, makeDepositUseCase, getBalanceUseCase)
	authHandler := handler.NewAuthHandler(loginUseCase)
	transactionHandler := handler.NewTransactionHandler(makeTransferUseCase, getExtractUseCase)

	authMiddleware := customMiddleware.NewAuthMiddleware(tokenService)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/accounts", accountHandler.CreateAccount)
	r.Post("/login", authHandler.Login)
	r.Post("/accounts/{account_id}/deposit", accountHandler.Deposit)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Handle)
		r.Post("/transactions", transactionHandler.MakeTransfer)
		r.Get("/transactions", transactionHandler.GetExtract)
		// Nova Rota Protegida
		r.Get("/accounts/{account_id}/balance", accountHandler.GetBalance)
	})

	fmt.Printf("Nexus Finance API running on port %s\n", conf.WebServerPort)
	fmt.Printf("Documentation available at http://localhost:%s/swagger/index.html\n", conf.WebServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", conf.WebServerPort), r)
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate(&entity.Account{}, &entity.Transaction{})
}
