package api

import (
	"github.com/gin-gonic/gin"

	"plusz-backend/api/service"
)

func Init() {
	router := gin.Default()

	router.POST("/addScheduleVersion", service.AddScheduleRevision)

	router.POST("/registerUser", service.RegisterUser)

	router.POST("/authenticateUser", service.AuthenticateUser)

	router.POST("/addNote", service.AddNote)

	router.GET("/getNote", service.GetNote)

	router.GET("/getUserSchedules", service.GetUserSchedules)

	router.GET("/getScheduleRevisions", service.GetScheduleRevisions)

	router.GET("/getRevisionClasses", service.GetRevisionClasses)

	router.DELETE("/usersAssignedSchedule", service.RemoveScheduleRevisionAssigment)

	if err := router.Run(); err != nil {
		panic(err)
	}
}
