package batch

import (
	"fmt"
	"plusz-backend/api/service"
	"plusz-backend/db"
	"time"
)

func CheckForNewSchedules() {
	fmt.Println("Batch created")
	for {
		time.Sleep(60 * time.Second)
		schedules, _ := db.GetAllSchedules()
		for _, schedule := range schedules {
			go func() {
				_, err := service.ScrapSchedule(schedule.Field, schedule.Year)
				if err != nil {
					fmt.Println(err)
				}
				fmt.Println("Updated schedule: ", schedule.Id)
			}()
		}
	}

}
