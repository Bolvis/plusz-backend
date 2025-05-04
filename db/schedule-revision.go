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
