package app

import (
	"log"
	"os"
	"testovoye/internal/database"
	"testovoye/internal/delivery/http"
	"testovoye/internal/repository"
	"testovoye/internal/service"
	"testovoye/internal/utils"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("[FATAL] %s", err.Error())
	}
}

func Start() {
	// dependencies
	db, err := database.InitPostgreSQL()
	if err != nil {
		log.Fatalf("[FATAL] %s", err.Error())
	}
	defer database.CloseDB()

	newFetcher := utils.NewFetcherAPIByName(time.Second * 5)

	// r&s&h
	repos := repository.NewRepositories(db)
	services := service.NewService(service.Deps{
		Repo:          repos,
		FetcherExtAPI: *newFetcher,
		Domain:        os.Getenv("SERVER_DOMAIN"),
	})
	handlers := http.NewHandlerAPI(services)

	// server
	server := handlers.Init()

	server.Run(os.Getenv("SERVER_ADDR"))
}
