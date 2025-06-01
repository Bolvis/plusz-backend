package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type ScheduleRevision struct {
	Id      string   `json:"id"`
	Date    string   `json:"date"`
	Classes []*Class `json:"classes"`
}

func GetScheduleRevisionId(scheduleRevision *ScheduleRevision, scheduleId string) (*ScheduleRevision, bool, error) {
	db, err := Connect()
	defer db.Close()
	isNew := false

	if err != nil {
		return scheduleRevision, isNew, err
	}

	searchQuery := `
		SELECT id
		FROM schedule_revision 
		WHERE schedule_id = $1 AND date = $2
	`

	if err = db.QueryRow(searchQuery, scheduleId, scheduleRevision.Date).Scan(&scheduleRevision.Id); errors.Is(err, sql.ErrNoRows) {
		isNew = true
		fmt.Println(err)
		insertQuery := `INSERT INTO schedule_revision (schedule_id, date) VALUES ($1, $2) RETURNING id`
		stmt, err := db.Prepare(insertQuery)
		defer stmt.Close()
		if err != nil {
			return scheduleRevision, isNew, err
		}

		if err = stmt.QueryRow(scheduleId, scheduleRevision.Date).Scan(&scheduleRevision.Id); err != nil {
			return scheduleRevision, isNew, err
		}
	} else if err != nil {
		return scheduleRevision, isNew, err
	}

	return scheduleRevision, isNew, nil
}

func GetScheduleRevisions(scheduleId string) ([]*ScheduleRevision, error) {
	db, err := Connect()
	defer db.Close()

	var scheduleRevisions []*ScheduleRevision
	if err != nil {
		return scheduleRevisions, err
	}

	query := `
		SELECT 
			id,
			date
		FROM schedule_revision
		WHERE schedule_id = $1
	`

	rows, err := db.Query(query, scheduleId)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return scheduleRevisions, err
	}

	for rows.Next() {
		var scheduleRevision ScheduleRevision
		if err := rows.Scan(&scheduleRevision.Id, &scheduleRevision.Date); err != nil {
			fmt.Println(err)
			return scheduleRevisions, err
		}
		scheduleRevisions = append(scheduleRevisions, &scheduleRevision)
	}

	return scheduleRevisions, nil
}

func GetPreviousRevision(scheduleId string) (ScheduleRevision, error) {
	db, err := Connect()
	defer db.Close()

	var scheduleRevision ScheduleRevision
	if err != nil {
		return scheduleRevision, err
	}

	queryRevision := `
		SELECT 
			id,
			date
		FROM schedule_revision
		WHERE schedule_id = $1
		ORDER BY date DESC
		LIMIT 2
	`

	rows, err := db.Query(queryRevision, scheduleId)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return scheduleRevision, err
	}

	var revisions []ScheduleRevision
	for rows.Next() {
		var revision ScheduleRevision
		if err = rows.Scan(
			&revision.Id,
			&revision.Date,
		); err != nil {
			fmt.Println(err)
			return scheduleRevision, err
		}
		revisions = append(revisions, revision)
	}

	if len(revisions) == 1 {
		return revisions[0], nil
	}

	scheduleRevision = revisions[1]

	queryClasses := `
		SELECT
		    id,
			date,
			start_hour,
			end_hour,
			name,
			lecturer,
			group_number,
			class_number,
			changed
		FROM class
		WHERE schedule_revision_id = $1
	`

	rows, err = db.Query(queryClasses, scheduleRevision.Id)
	defer rows.Close()
	if err != nil {
		fmt.Println(err)
		return scheduleRevision, err
	}

	for rows.Next() {
		var class Class
		if err = rows.Scan(
			&class.Id,
			&class.Date,
			&class.StartHour,
			&class.EndHour,
			&class.Name,
			&class.Lecturer,
			&class.Group,
			&class.ClassNumber,
			&class.Changed,
		); err != nil {
			fmt.Println(err)
			return scheduleRevision, err
		}
		scheduleRevision.Classes = append(scheduleRevision.Classes, &class)
	}

	return scheduleRevision, nil
}
