package app

import (
	"bookstore-users-api/logger"
	"fmt"

	"github.com/gin-gonic/gin"
)

var (
	router = gin.Default()
)

func StartApplication() {
	mapUrls()
	logger.Info("about to start the application......")
	router.Run(":9191")
	fmt.Println("start application")
}
