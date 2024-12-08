package service

import (
	"net/http"
	"plusz-backend/scrapper"
	"strings"

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

	scheduleRevision, err := scrapper.Scrap(url)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, scheduleRevision)
}
