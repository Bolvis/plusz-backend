package api

import (
	"plusz-backend/api/service"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()

	router.POST("/addScheduleVersion", service.AddScheduleRevision)

	router.GET("/getUserSchedules", service.GetUserSchedules)

	if err := router.Run(":2013"); err != nil {
		panic(err)
	}
}
