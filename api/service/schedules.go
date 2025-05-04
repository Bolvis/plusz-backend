package service

import (
	"fmt"
	"net/http"
	"plusz-backend/api/authorization"
	"strings"

	"plusz-backend/db"
	"plusz-backend/scrapper"

	"github.com/gin-gonic/gin"
)

type scheduleRequest struct {
	Year  string `json:"year"`
	Field string `json:"field"`
}

func AddScheduleRevision(c *gin.Context) {
	var request scheduleRequest
	tokenString := c.Request.Header.Get("Authorization")

	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := authorization.VerifyToken(tokenString)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	schedule, err := ScrapSchedule(request.Field, request.Year)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := db.AssignUserSchedule(token.UserId, schedule.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func ScrapSchedule(field string, year string) (db.Schedule, error) {
	url := strings.Join([]string{
		"https://efz.usz.edu.pl/wp-content/include-me/plany_mick/zajecia_xml.php?kierunek=",
		field,
		"&rok=",
		year,
	}, "")

	schedule := db.Schedule{Field: field, Year: year}
	schedule, err := scrapper.ScrapUSZ(url, schedule)
	if err != nil {
		return schedule, err
	}

	if schedule, err = db.GetScheduleId(schedule); err != nil {
		return schedule, err
	}

	isScheduleRevisionNew := false
	if schedule.ScheduleRevisions[0], isScheduleRevisionNew, err = db.GetScheduleRevisionId(schedule.ScheduleRevisions[0], schedule.Id); err != nil {
		fmt.Println(err)
		return schedule, err
	}

	if isScheduleRevisionNew {
		if err = db.InsertClasses(schedule.ScheduleRevisions[0].Classes, schedule.ScheduleRevisions[0].Id); err != nil {
			fmt.Println(err)
			return schedule, err
		}
	}

	return schedule, nil
}

func GetUserSchedules(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	token, err := authorization.VerifyToken(tokenString)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	schedules, err := db.GetUserSchedules(token.UserId)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedules)
}

func GetScheduleRevisions(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	_, err := authorization.VerifyToken(tokenString)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	scheduleId := c.Query("scheduleId")
	schedule := db.Schedule{Id: scheduleId}
	var scheduleRevisions []*db.ScheduleRevision
	if scheduleRevisions, err = db.GetScheduleRevisions(schedule.Id); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	schedule.ScheduleRevisions = scheduleRevisions

	c.JSON(http.StatusOK, schedule)
}

func GetRevisionClasses(c *gin.Context) {
	revisionId := c.Query("revisionId")
	tokenString := c.Request.Header.Get("Authorization")

	token, err := authorization.VerifyToken(tokenString)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	revision := db.ScheduleRevision{Id: revisionId}
	var classes []*db.Class
	if classes, err = db.GetScheduleRevisionClasses(token.UserId, revision.Id); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	revision.Classes = classes

	c.JSON(http.StatusOK, revision)
}

func RemoveScheduleRevisionAssigment(c *gin.Context) {
	scheduleId := c.Query("scheduleId")
	tokenString := c.Request.Header.Get("Authorization")

	token, err := authorization.VerifyToken(tokenString)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	if err := db.RemoveUserScheduleAssigment(token.UserId, scheduleId); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
