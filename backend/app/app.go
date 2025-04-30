package app

import (
	"api-fio/db"
	"api-fio/logging"
	"api-fio/routes"
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()

	db.InitDatabase()
	defer db.CloseDB()

	logger := logging.NewLogger()
	defer logger.Sync()

	routes.InitRoutes(r, logger, db.GetDB())

	r.Run(fmt.Sprintf("%s:%s",
		os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT")))
}
