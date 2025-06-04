package scrapper

import (
	"fmt"
	"strings"
	"time"

	"plusz-backend/db"
	"plusz-backend/util"

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
			if len(hours) >= 2 {
				class.StartHour = hours[0]
				class.EndHour = hours[1]
			}
			class.Name = line[2]
			class.Lecturer = line[3]
			class.Group = line[4]
			class.ClassNumber = line[5]

			scheduleRevision.Classes = append(scheduleRevision.Classes, &class)
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
			scheduleRevision.Date =
				util.ConvertToDate(strings.TrimSpace(strings.Split(e.Text, ":")[1]), ".").Format(time.DateOnly)
		}
	})

	if err := metadataCollector.Visit(url); err != nil {
		fmt.Println("Failed to scrap USZ metadata")
		return schedule, err
	}

	schedule.ScheduleRevisions = []*db.ScheduleRevision{&scheduleRevision}

	return schedule, nil
}

func ScrapUSZRoom(url string, schedule db.Schedule, classNumber string) (db.Schedule, error) {
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
		if len(line) == 5 && line[0] != "Data" {
			class.Date = line[0]
			hours := strings.Split(strings.Replace(line[1], ".", ":", 2), "-")
			if len(hours) >= 2 {
				class.StartHour = hours[0]
				class.EndHour = hours[1]
			}
			class.Name = line[2]
			class.Lecturer = line[3]
			class.Group = line[4]
			class.ClassNumber = classNumber

			scheduleRevision.Classes = append(scheduleRevision.Classes, &class)
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
			scheduleRevision.Date =
				util.ConvertToDate(strings.TrimSpace(strings.Split(e.Text, ":")[1]), ".").Format(time.DateOnly)
		}
	})

	if err := metadataCollector.Visit(url); err != nil {
		fmt.Println("Failed to scrap USZ metadata")
		return schedule, err
	}

	schedule.ScheduleRevisions = []*db.ScheduleRevision{&scheduleRevision}

	return schedule, nil
}

func ScrapUSZLecturer(url string, schedule db.Schedule, lecturerName string) (db.Schedule, error) {
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
		if len(line) == 5 && line[0] != "Data" {
			class.Date = line[0]
			hours := strings.Split(strings.Replace(line[1], ".", ":", 2), "-")
			if len(hours) >= 2 {
				class.StartHour = hours[0]
				class.EndHour = hours[1]
			}
			class.Name = line[2]
			class.Lecturer = lecturerName
			class.Group = line[3]
			class.ClassNumber = line[4]

			scheduleRevision.Classes = append(scheduleRevision.Classes, &class)
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
			scheduleRevision.Date =
				util.ConvertToDate(strings.TrimSpace(strings.Split(e.Text, ":")[1]), ".").Format(time.DateOnly)
		}
	})

	if err := metadataCollector.Visit(url); err != nil {
		fmt.Println("Failed to scrap USZ metadata")
		return schedule, err
	}

	schedule.ScheduleRevisions = []*db.ScheduleRevision{&scheduleRevision}

	return schedule, nil
}

func initColly() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("efz.usz.edu.pl"),
	)

	c.OnRequest(func(r *colly.Request) {
		// without headers, usz domain sometimes blocks calls - perhaps as prevention from DDoS
		r.Headers.Set("User-Agent", "Mozilla/5.0 (Windows NT 6.1; Win64; x64)")
	})

	return c
}
