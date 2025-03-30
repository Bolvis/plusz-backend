package service

import (
	"net/http"
	"strings"

	"plusz-backend/db"
	"plusz-backend/scrapper"

	"github.com/gin-gonic/gin"
)

type scheduleRequest struct {
	Year  string  `json:"year"`
	Field string  `json:"field"`
	User  db.User `json:"user"`
}

func AddScheduleRevision(c *gin.Context) {
	var request scheduleRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.AuthUser(&request.User); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	url := strings.Join([]string{
		"https://efz.usz.edu.pl/wp-content/include-me/plany_mick/zajecia_xml.php?kierunek=",
		request.Field,
		"&rok=",
		request.Year,
	}, "")

	schedule := db.Schedule{Field: request.Field, Year: request.Year}
	schedule, err := scrapper.ScrapUSZ(url, schedule)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if schedule, err = db.GetScheduleId(schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	isScheduleRevisionNew := false
	if schedule.ScheduleRevisions[0], isScheduleRevisionNew, err = db.GetScheduleRevisionId(schedule.ScheduleRevisions[0], schedule.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if isScheduleRevisionNew {
		if err = db.InsertClasses(schedule.ScheduleRevisions[0].Classes, schedule.ScheduleRevisions[0].Id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	if err := db.AssignUserSchedule(request.User, schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func GetUserSchedules(c *gin.Context) {
	login := c.Query("login")
	password := c.Query("password")

	user := db.User{Login: login, Password: password}

	if err := db.AuthUser(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	schedules, err := db.GetUserSchedules(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

func GetScheduleRevisions(c *gin.Context) {
	scheduleId := c.Query("scheduleId")
	schedule := db.Schedule{Id: scheduleId}
	var scheduleRevisions []*db.ScheduleRevision
	var err error
	if scheduleRevisions, err = db.GetScheduleRevisions(schedule.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	//for j, scheduleRevision := range scheduleRevisions {
	//	var classes []*db.Class
	//	if classes, err = db.GetScheduleRevisionClasses(scheduleRevision.Id); err != nil {
	//		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	//	}
	//	scheduleRevisions[j].Classes = classes
	//}
	schedule.ScheduleRevisions = scheduleRevisions

	c.JSON(http.StatusOK, schedule)
}

func GetRevisionClasses(c *gin.Context) {
	revisionId := c.Query("revisionId")

	revision := db.ScheduleRevision{Id: revisionId}
	var classes []*db.Class
	var err error
	if classes, err = db.GetScheduleRevisionClasses(revision.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	revision.Classes = classes

	c.JSON(http.StatusOK, revision)
}

func RemoveScheduleRevisionAssigment(c *gin.Context) {
	scheduleId := c.Query("scheduleId")
	user := db.User{Login: c.Query("login"), Password: c.Query("password")}
	if err := db.AuthUser(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	if err := db.RemoveUserScheduleAssigment(user.Id, scheduleId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
