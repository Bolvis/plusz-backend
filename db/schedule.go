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
	ScheduleRevisions []*ScheduleRevision `json:"schedule_revisions"`
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
		WHERE field = $1 AND year = $2 AND academic_year = $3 AND semester = $4
	`

	if err = db.QueryRow(searchQuery, schedule.Field, schedule.Year, schedule.AcademicYear, schedule.Semester).Scan(&schedule.Id); errors.Is(err, sql.ErrNoRows) {

		insertQuery := `INSERT INTO schedule (field, year, academic_year, semester) VALUES ($1, $2, $3, $4) RETURNING id`
		stmt, err := db.Prepare(insertQuery)
		defer stmt.Close()
		if err != nil {
			return schedule, err
		}

		if err = stmt.QueryRow(schedule.Field, schedule.Year, schedule.AcademicYear, schedule.Semester).Scan(&schedule.Id); err != nil {
			return schedule, err
		}
	}

	if err != nil {
		fmt.Println(err)
		return schedule, err
	}

	return schedule, nil
}
