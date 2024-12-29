package api

import (
	"github.com/gin-gonic/gin"
	"plusz-backend/api/service"
)

func Init() {
	router := gin.Default()

	router.POST("/addScheduleVersion", service.AddScheduleRevision)

	router.GET("/getUserSchedules", service.AddScheduleRevision)

	if err := router.Run(":2013"); err != nil {
		panic(err)
	}
}
