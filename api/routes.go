package api

import (
	"plusz-backend/api/authorization"
	"plusz-backend/api/service"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()
	defaultPrefix := "/api/v1"

	openRoutes := router.Group(defaultPrefix)

	openRoutes.POST("/User/authenticate", service.AuthenticateUser)

	openRoutes.POST("/User/register", service.RegisterUser)

	protectedRoutes := router.Group(defaultPrefix)

	protectedRoutes.Use(authorization.AuthMiddleware)

	protectedRoutes.GET("/ScheduleVersion/:revisionId/Classes", service.GetRevisionClasses)

	protectedRoutes.POST("/ScheduleVersion/USZ/Student/add", service.AddScheduleRevision)

	protectedRoutes.POST("/ScheduleVersion/USZ/Room/add", service.AddRoomScheduleRevision)

	protectedRoutes.POST("/ScheduleVersion/USZ/Lecturer/add", service.AddLecturerScheduleRevision)

	protectedRoutes.POST("/Note/add", service.AddNote)

	protectedRoutes.GET("/Class/:classId/Note", service.GetNote)

	protectedRoutes.GET("/CurrentUser/Schedule", service.GetUserSchedules)

	protectedRoutes.GET("/CurrentUser/Schedule/:scheduleId/ScheduleVersions", service.GetScheduleRevisions)

	protectedRoutes.DELETE("/CurrentUser/Schedule/:scheduleId/removeAssignment", service.RemoveScheduleRevisionAssignment)

	if err := router.Run(); err != nil {
		panic(err)
	}
}
