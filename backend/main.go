// @title PersonAPI
// @version 1.0
// @description Testovoye zadanie
// @host localhost:8000
// @BasePath /

package main

import (
	"api-fio/app"
	"log"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("[FATAL] %s", err.Error())
	}
}

var (
	isContainer = true
	timeToWait  = time.Second * 5
)

func main() {
	if isContainer {
		time.Sleep(timeToWait)
	}
	app.Start()
}
