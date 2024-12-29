package db

import "fmt"

type ScheduleRevision struct {
	Id      string   `json:"id"`
	Date    string   `json:"date"`
	Classes []*Class `json:"classes"`
}

func GetScheduleRevisionId(scheduleRevision *ScheduleRevision, scheduleId string) (*ScheduleRevision, error) {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		return scheduleRevision, err
	}

	searchQuery := `
		SELECT id
		FROM schedule_revision 
		WHERE schedule_id = $1 AND date = $2
	`

	if err = db.QueryRow(searchQuery, scheduleId, scheduleRevision.Date).Scan(&scheduleRevision.Id); err != nil {
		fmt.Println(err)
		fmt.Println("inserting a new schedule revision...")
		insertQuery := `INSERT INTO schedule_revision (schedule_id, date) VALUES ($1, $2) RETURNING id`
		stmt, err := db.Prepare(insertQuery)
		defer stmt.Close()
		if err != nil {
			return scheduleRevision, err
		}

		if err = stmt.QueryRow(scheduleId, scheduleRevision.Date).Scan(&scheduleRevision.Id); err != nil {
			return scheduleRevision, err
		}
		fmt.Println("inserted successfully")
	}

	return scheduleRevision, nil
}
