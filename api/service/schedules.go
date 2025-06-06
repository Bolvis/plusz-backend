package service

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"plusz-backend/util"
	"slices"
	"strings"

	"plusz-backend/api/authorization"
	"plusz-backend/db"
	"plusz-backend/scrapper"

	"github.com/gin-gonic/gin"
)

type scheduleRequest struct {
	Year  string `json:"year"`
	Field string `json:"field"`
}

type roomScheduleRequest struct {
	Name string `json:"name"`
}

type lecturerScheduleRequest struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
}

type FieldChanges struct {
	ChangeType string        `json:"changeType"`
	Changes    []FieldChange `json:"changes"`
}

type FieldChange struct {
	Key   string `json:"key"`
	Value string `json:"value"`
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

	schedule, err := ScrapSchedule(request.Field, request.Year, "USZ")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = db.AssignUserSchedule(token.UserId, schedule.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func AddRoomScheduleRevision(c *gin.Context) {
	var request roomScheduleRequest
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

	schedule, err := ScrapSchedule(request.Name, "", "USZRoom")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = db.AssignUserSchedule(token.UserId, schedule.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func AddLecturerScheduleRevision(c *gin.Context) {
	var request lecturerScheduleRequest
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

	schedule, err := ScrapSchedule(request.Name, request.Surname, "USZLecturer")
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = db.AssignUserSchedule(token.UserId, schedule.Id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}

func ScrapSchedule(queryField1 string, queryField2 string, scrapperType string) (db.Schedule, error) {
	var url string
	var err error
	var schedule db.Schedule

	if scrapperType == "USZ" {
		schedule = db.Schedule{Field: queryField1, Year: queryField2, ScheduleType: scrapperType}
		url = strings.Join([]string{
			"https://efz.usz.edu.pl/wp-content/include-me/plany_mick/zajecia_xml.php?kierunek=",
			queryField1,
			"&rok=",
			queryField2,
		}, "")
		schedule, err = scrapper.ScrapUSZ(url, schedule)
		if err != nil {
			return schedule, err
		}
	}

	if scrapperType == "USZRoom" {
		formattedName := strings.ReplaceAll(queryField1, "_", " ")
		schedule = db.Schedule{Field: formattedName, Year: queryField2, ScheduleType: scrapperType}
		url = strings.Join([]string{
			"https://efz.usz.edu.pl/wp-content/include-me/plany_mick/sala_zajecia_xml.php?nazw=",
			queryField1,
		}, "")
		schedule, err = scrapper.ScrapUSZRoom(url, schedule, strings.ReplaceAll(formattedName, "_", " "))
		if err != nil {
			return schedule, err
		}
	}

	if scrapperType == "USZLecturer" {
		formattedName := strings.ReplaceAll(queryField1, "%20", " ")
		schedule = db.Schedule{Field: formattedName, Year: queryField2, ScheduleType: scrapperType}
		url = strings.Join([]string{
			"https://efz.usz.edu.pl/wp-content/include-me/plany_mick/pracow_zajecia_xml.php?nazw=",
			queryField1,
			"&nazw1=",
			queryField2,
		}, "")
		schedule, err = scrapper.ScrapUSZLecturer(url, schedule, strings.Join(
			[]string{formattedName, queryField2}, " "))
		if err != nil {
			return schedule, err
		}
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
		if err = ProcessBeforeInsert(&schedule); err != nil {
			fmt.Println(err)
			return schedule, err
		}

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

	scheduleId := c.Param("scheduleId")
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
	revisionId := c.Param("revisionId")
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
	scheduleId := c.Param("scheduleId")
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

func ProcessBeforeInsert(newSchedule *db.Schedule) error {
	previousRevision, err := db.GetPreviousRevision(newSchedule.Id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil
	} else if err != nil {
		fmt.Println(err)
		return err
	}

	var addedClasses []*db.Class
	var foundedMatchesIds []string
	isNew := len(previousRevision.Classes) > 0
	for _, newClass := range newSchedule.ScheduleRevisions[0].Classes {
		changes := FieldChanges{}
		newClassDate := util.ConvertToDate(newClass.Date, "/")

		for _, previousClass := range previousRevision.Classes {
			foundedMatchesIds = append(foundedMatchesIds, previousClass.Id)
			previousClassDate := util.ConvertToDate(previousClass.Date, "-")

			if previousClassDate == newClassDate && previousClass.Name == newClass.Name {
				isNew = false
				previousClassStartHourFormatted := util.FormatTime(previousClass.StartHour)
				previousClassEndHourFormatted := util.FormatTime(previousClass.EndHour)

				if previousClassStartHourFormatted != newClass.StartHour {
					changes.Changes = append(
						changes.Changes,
						FieldChange{
							Key:   "StartHour",
							Value: previousClassStartHourFormatted,
						},
					)
				}

				if previousClassEndHourFormatted != newClass.EndHour {
					changes.Changes = append(
						changes.Changes,
						FieldChange{
							Key:   "EndHour",
							Value: previousClassEndHourFormatted,
						},
					)
				}

				if previousClass.ClassNumber != newClass.ClassNumber {
					changes.Changes = append(
						changes.Changes,
						FieldChange{
							Key:   "ClassNumber",
							Value: previousClass.ClassNumber,
						},
					)
				}

				if previousClass.Group != newClass.Group {
					changes.Changes = append(
						changes.Changes,
						FieldChange{
							Key:   "Group",
							Value: previousClass.Group,
						},
					)
				}

				if previousClass.Lecturer != newClass.Lecturer {
					changes.Changes = append(
						changes.Changes,
						FieldChange{
							Key:   "Lecturer",
							Value: previousClass.Lecturer,
						},
					)
				}

				break
			}
		}

		if len(changes.Changes) > 0 {
			changes.ChangeType = "changed"
		} else if isNew {
			addedClasses = append(addedClasses, newClass)
			changes.ChangeType = "added"
		} else {
			changes.ChangeType = "none"
		}

		changesBytes, err := json.Marshal(changes)
		if err != nil {
			fmt.Println(err)
			return err
		}
		newClass.Changed = string(changesBytes)
	}

	deleteChange := FieldChanges{
		ChangeType: "delete",
	}
	for _, previousClass := range previousRevision.Classes {
		if !slices.Contains(foundedMatchesIds, previousClass.Id) {
			deleteChangeBytes, err := json.Marshal(deleteChange)
			if err != nil {
				fmt.Println(err)
				return err
			}
			previousClass.Changed = string(deleteChangeBytes)
			previousClass.StartHour = util.FormatTime(previousClass.StartHour)
			previousClass.EndHour = util.FormatTime(previousClass.EndHour)
			newSchedule.ScheduleRevisions[0].Classes = append(newSchedule.ScheduleRevisions[0].Classes, previousClass)
		}
	}

	return nil
}
