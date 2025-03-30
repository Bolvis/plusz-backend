package api

import (
	"plusz-backend/api/service"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()

	router.POST("/addScheduleVersion", service.AddScheduleRevision)

	router.POST("/addNote", service.AddNote)

	router.GET("/getNote", service.GetNote)

	router.GET("/getUserSchedules", service.GetUserSchedules)

	router.GET("/getScheduleRevisions", service.GetScheduleRevisions)

	router.GET("/getRevisionClasses", service.GetRevisionClasses)

	router.DELETE("/usersAssignedSchedule", service.RemoveScheduleRevisionAssigment)

	if err := router.Run(":2013"); err != nil {
		panic(err)
	}
}
