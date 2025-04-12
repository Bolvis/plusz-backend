package api

import (
	"plusz-backend/api/service"

	"strings"

	"github.com/gin-gonic/gin"
)

func Init() {
	router := gin.Default()
	defaultPrefix := "/api/v1"

	router.GET(strings.Join([]string{defaultPrefix, "/ScheduleVersion/:revisionId/Classes"}, ""), service.GetRevisionClasses)

	router.POST(strings.Join([]string{defaultPrefix, "/ScheduleVersion/USZ/Student/add"}, ""), service.AddScheduleRevision)

	router.POST(strings.Join([]string{defaultPrefix, "/ScheduleVersion/USZ/Room/add"}, ""), service.AddRoomScheduleRevision)

	router.POST(strings.Join([]string{defaultPrefix, "/ScheduleVersion/USZ/Lecturer/add"}, ""), service.AddLecturerScheduleRevision)

	router.POST(strings.Join([]string{defaultPrefix, "/Note/add"}, ""), service.AddNote)

	router.GET(strings.Join([]string{defaultPrefix, "/Class/:classId/Note"}, ""), service.GetNote)

	router.POST(strings.Join([]string{defaultPrefix, "/User/authenticate"}, ""), service.AuthenticateUser)

	router.POST(strings.Join([]string{defaultPrefix, "/User/register"}, ""), service.RegisterUser)

	router.GET(strings.Join([]string{defaultPrefix, "/CurrentUser/Schedule"}, ""), service.GetUserSchedules)

	router.GET(strings.Join([]string{defaultPrefix, "/CurrentUser/Schedule/:scheduleId/ScheduleVersions"}, ""), service.GetScheduleRevisions)

	router.DELETE(strings.Join([]string{defaultPrefix, "/CurrentUser/ScheduleVersion/removeAssigment"}, ""), service.RemoveScheduleRevisionAssigment)

	if err := router.Run(); err != nil {
		panic(err)
	}
}
