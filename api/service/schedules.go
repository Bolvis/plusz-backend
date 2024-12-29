package service

import (
	"net/http"
	"plusz-backend/db"
	"strings"

	"plusz-backend/scrapper"

	"github.com/gin-gonic/gin"
)

type scheduleRequest struct {
	Year  string `json:"year"`
	Field string `json:"field"`
}

func GetSchedule(c *gin.Context) {
	var request scheduleRequest

	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	url := strings.Join([]string{
		"https://efz.usz.edu.pl/wp-content/include-me/plany_mick/zajecia_xml.php?kierunek=",
		request.Field,
		"&rok=",
		request.Year,
	}, "")

	_, err := scrapper.ScrapUSZ(url)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	schedule := db.Schedule{Field: request.Field, Year: request.Year, AcademicYear: "2024/2027"}
	if schedule, err = db.GetScheduleId(schedule); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, schedule)
}
