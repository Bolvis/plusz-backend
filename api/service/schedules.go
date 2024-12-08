package service

import (
	"fmt"
	"plusz-backend/scrapper"
	"strings"

	"github.com/gin-gonic/gin"
)

type scheduleRequest struct {
	year  string
	field string
}

func GetSchedule(c *gin.Context) {
	var request scheduleRequest

	if err := c.BindJSON(&request); err != nil {
		return
	}

	url := strings.Join([]string{
		"https://efz.usz.edu.pl/wp-content/include-me/plany_mick/zajecia_xml.php?kierunek=",
		request.field,
		"&rok=",
		request.year,
	}, "")

	scheduleRevision, err := scrapper.Scrap(url)
	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(scheduleRevision.Date)
	for _, class := range scheduleRevision.Classes {
		fmt.Println(class)
	}
}
