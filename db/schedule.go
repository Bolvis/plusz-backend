package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Schedule struct {
	Id                string              `json:"id"`
	Field             string              `json:"field"`
	Year              string              `json:"year"`
	AcademicYear      string              `json:"academic_year"`
	Semester          string              `json:"semester"`
	ScheduleType      string              `json:"schedule_type"`
	ScheduleRevisions []*ScheduleRevision `json:"schedule_revisions"`
}

func GetScheduleId(schedule Schedule) (Schedule, error) {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return schedule, err
	}

	searchQuery := `
		SELECT id
		FROM schedule 
		WHERE field = $1 AND year = $2 AND academic_year = $3 AND semester = $4
	`

	if err = db.QueryRow(searchQuery, schedule.Field, schedule.Year, schedule.AcademicYear, schedule.Semester).Scan(&schedule.Id); errors.Is(err, sql.ErrNoRows) {

		insertQuery := `INSERT INTO schedule (field, year, academic_year, semester, schedule_type) VALUES ($1, $2, $3, $4, $5) RETURNING id`
		stmt, err := db.Prepare(insertQuery)
		defer stmt.Close()
		if err != nil {
			fmt.Println(err)
			return schedule, err
		}

		if err = stmt.QueryRow(schedule.Field, schedule.Year, schedule.AcademicYear, schedule.Semester, schedule.ScheduleType).Scan(&schedule.Id); err != nil {
			fmt.Println(err)
			return schedule, err
		}
	}

	return schedule, nil
}

func GetAllSchedules() ([]Schedule, error) {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var schedules []Schedule

	result, err := db.Query(`SELECT id, field, year, academic_year, semester, schedule_type FROM schedule`)
	defer result.Close()
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	for result.Next() {
		var schedule Schedule
		err = result.Scan(&schedule.Id, &schedule.Field, &schedule.Year, &schedule.AcademicYear, &schedule.Semester, &schedule.ScheduleType)
		if err != nil {
			fmt.Println(err)
		}
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
