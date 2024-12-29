package scrapper

import (
	"plusz-backend/db"
	"plusz-backend/util"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func ScrapUSZ(url string) (db.ScheduleRevision, error) {
	c := initColly()

	var scheduleRevision db.ScheduleRevision
	c.OnHTML("tr", func(e *colly.HTMLElement) {
		var line []string
		e.ForEach("td", func(i int, rowElement *colly.HTMLElement) {
			line = append(line, util.StandardizeSpaces(rowElement.DOM.Text()))
		})

		var class db.Class
		if len(line) == 6 && line[0] != "Data" {
			class.Date = line[0]
			hours := strings.Split(strings.Replace(line[1], ".", ":", 2), "-")
			class.StartHour = hours[0]
			class.EndHour = hours[1]
			class.Name = line[2]
			class.Lecturer = line[3]
			class.Group = line[4]
			class.ClassNumber = line[5]

			scheduleRevision.Classes = append(scheduleRevision.Classes, class)
		}
	})
	if err := c.Visit(url); err != nil {
		return scheduleRevision, err
	}

	scheduleRevision.Date = time.Now().Format(`2006-01-02`)

	return scheduleRevision, nil
}

func initColly() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("efz.usz.edu.pl"),
	)

	c.OnRequest(func(r *colly.Request) {
		// without headers usz domain sometimes blocks calls - perhaps as prevention from DDoS
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64)")
	})

	return c
}
