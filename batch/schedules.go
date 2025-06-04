package batch

import (
	"fmt"
	"plusz-backend/api/service"
	"plusz-backend/db"
	"strings"
	"time"
)

func CheckForNewSchedules() {
	fmt.Println("Batch created")
	for {
		time.Sleep(60 * time.Second)
		schedules, _ := db.GetAllSchedules()
		for _, schedule := range schedules {
			go func() {
				var queryField1 string
				queryField2 := schedule.Year

				if schedule.ScheduleType == "USZ" {
					queryField1 = schedule.Field
				} else if schedule.ScheduleType == "USZLecturer" {
					queryField1 = strings.ReplaceAll(schedule.Field, " ", "%20")
				} else if schedule.ScheduleType == "USZRoom" {
					queryField1 = strings.ReplaceAll(schedule.Field, " ", "_")
				}

				_, err := service.ScrapSchedule(queryField1, queryField2, schedule.ScheduleType)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Updated schedule: ", schedule.Id)
			}()
		}
	}

}
