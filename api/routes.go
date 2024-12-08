package api

import (
	"github.com/gin-gonic/gin"
	"plusz-backend/api/service"
)

func Init() {
	router := gin.Default()

	router.POST("/srcaped", service.GetSchedule)

	if err := router.Run("localhost:8080"); err != nil {
		panic(err)
	}
}
