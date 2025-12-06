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

	accountRepo := repository.NewAccountRepositoryPostgres(db)
	createAccountUseCase := usecase.NewCreateAccountUseCase(accountRepo)
	accountHandler := handler.NewAccountHandler(createAccountUseCase)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/accounts", accountHandler.CreateAccount)

	fmt.Printf("Nexus Finance API running on port %s\n", conf.WebServerPort)
	http.ListenAndServe(fmt.Sprintf(":%s", conf.WebServerPort), r)
}

func runMigrations(db *gorm.DB) {
	db.AutoMigrate(&entity.Account{})
}
