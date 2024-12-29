package db

import (
	"fmt"
)

type Schedule struct {
	Id                string             `json:"id"`
	Field             string             `json:"field"`
	Year              string             `json:"year"`
	AcademicYear      string             `json:"academic_year"`
	ScheduleRevisions []ScheduleRevision `json:"schedule_revisions"`
}

func GetScheduleId(schedule Schedule) (Schedule, error) {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		return schedule, err
	}

	searchQuery := `
		SELECT id
		FROM schedule 
		WHERE field = $1 AND year = $2 AND academic_year = $3
	`

	if err = db.QueryRow(searchQuery, schedule.Field, schedule.Year, schedule.AcademicYear).Scan(&schedule.Id); err != nil {
		fmt.Println(err)
		fmt.Println("inserting a new schedule...")
		insertQuery := `INSERT INTO schedule (field, year, academic_year) VALUES ($1, $2, $3) RETURNING id`
		stmt, err := db.Prepare(insertQuery)
		if err != nil {
			return schedule, err
		}
		defer stmt.Close()

		if err = stmt.QueryRow(schedule.Field, schedule.Year, schedule.AcademicYear).Scan(&schedule.Id); err != nil {
			return schedule, err
		}
		fmt.Println("inserted successfully")
	}
	fmt.Println(schedule)
	return schedule, nil
}
