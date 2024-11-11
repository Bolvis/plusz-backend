package scrapper

import (
	"fmt"
	"github.com/gocolly/colly"
	"plusz-backend/db"
	"plusz-backend/util"
	"time"
)

var scheduleRevision db.ScheduleRevision

func Scrap(url string) db.ScheduleRevision {
	c := initColly()

	if err := c.Visit(url); err != nil {
		fmt.Println(err.Error())
	}

	scheduleRevision.Date = time.Now().Format(`2006-01-02`)

	return scheduleRevision
}

func initColly() *colly.Collector {
	c := colly.NewCollector(
		colly.AllowedDomains("efz.usz.edu.pl"),
	)

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting: ", r.URL)
	})

	// triggered when the scraper encounters an error
	c.OnError(func(_ *colly.Response, err error) {
		fmt.Println("Something went wrong: ", err)
	})

	// fired when the server responds
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Page visited: ", r.Request.URL)
	})

	// triggered when a CSS selector matches an element
	c.OnHTML("tr", onHTML)

	// triggered once scraping is done (e.g., write the data to a CSV file)
	c.OnScraped(func(r *colly.Response) {
		fmt.Println(r.Request.URL, " scraped!")
		fmt.Println(r.Request.Body)
	})

	return c
}

func onHTML(e *colly.HTMLElement) {
	var line []string
	e.ForEach("td", func(i int, rowElement *colly.HTMLElement) {
		line = append(line, util.StandardizeSpaces(rowElement.DOM.Text()))
	})

	var class db.Class
	if len(line) == 6 {
		class.Date = line[0]
		class.Hour = line[1]
		class.Name = line[2]
		class.Lecturer = line[3]
		class.Group = line[4]
		class.ClassNumber = line[5]
	}

	fmt.Print("[")
	for _, v := range line {
		fmt.Print(v, ", ")
	}
	fmt.Println("]")

	if len(line) > 1 {
		scheduleRevision.Classes = append(scheduleRevision.Classes, class)
	}
	//fmt.Println(strings.Fields(e.ChildText("td")))

}
