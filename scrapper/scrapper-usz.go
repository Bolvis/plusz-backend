package scrapper

import (
	"fmt"
	"plusz-backend/db"
	"plusz-backend/util"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

func ScrapUSZ(url string, schedule db.Schedule) (db.Schedule, error) {
	tableCollector := initColly()
	metadataCollector := tableCollector.Clone()
	metadataCollector.AllowURLRevisit = true

	var scheduleRevision db.ScheduleRevision
	tableCollector.OnHTML("tr", func(e *colly.HTMLElement) {
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
	if err := tableCollector.Visit(url); err != nil {
		fmt.Println("Failed to scrape usz table")
		return schedule, err
	}

	metadataCollector.OnHTML("b", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Rok akademicki:") {
			schedule.AcademicYear = strings.TrimSpace(strings.Split(e.Text, ":")[1])
		}
		if strings.Contains(e.Text, "Semestr:") {
			schedule.Semester = strings.TrimSpace(strings.Split(e.Text, ":")[1])
		}
	})

	metadataCollector.OnHTML("span", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "Data aktualizacji:") {
			dateArray := strings.Split(strings.TrimSpace(strings.Split(e.Text, ":")[1]), ".")
			year, err := strconv.Atoi(dateArray[2])
			if err != nil {
				fmt.Println(err)
			}
			month, err := strconv.Atoi(dateArray[1])
			if err != nil {
				fmt.Println(err)
			}
			day, err := strconv.Atoi(dateArray[0])
			if err != nil {
				fmt.Println(err)
			}
			schedule.LastUpdateDate = time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC).Format("2006-01-02")
		}
	})

	if err := metadataCollector.Visit(url); err != nil {
		fmt.Println("Failed to scrap USZ metadata")
		return schedule, err
	}

	scheduleRevision.Date = time.Now().Format(`2006-01-02`)
	schedule.ScheduleRevisions = []db.ScheduleRevision{scheduleRevision}

	return schedule, nil
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
