package db

import (
	"strings"

	"plusz-backend/util"
)

type Class struct {
	Id          string `json:"id"`
	Date        string `json:"date"`
	StartHour   string `json:"startHour"`
	EndHour     string `json:"endHour"`
	Name        string `json:"name"`
	Lecturer    string `json:"lecturer"`
	Group       string `json:"group"`
	ClassNumber string `json:"classNumber"`
}

func InsertClasses(classes []*Class, scheduleRevisionId string) error {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		return err
	}

	insertClasses := `
		INSERT INTO 
		    public.class (date, start_hour, end_hour, name, lecturer, group_number, class_number, schedule_revision_id) 
		VALUES `
	var classesVals []interface{}

	for _, v := range classes {
		insertClasses += "(?::date,?::time,?::time,?::varchar,?::varchar,?::varchar,?::varchar,?::integer),"
		classesVals = append(classesVals, v.Date, v.StartHour, v.EndHour, v.Name, v.Lecturer, v.Group, v.ClassNumber, scheduleRevisionId)
	}
	insertClasses = strings.TrimSuffix(insertClasses, ",")
	insertClasses += " RETURNING id"
	insertClasses = util.ReplaceSQL(insertClasses, "?")

	rows, err := db.Query(insertClasses, classesVals...)
	defer rows.Close()
	if err != nil {
		return err
	}

	classesIds := []string{}
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return err
		}
		classesIds = append(classesIds, id)
	}
	for i, class := range classes {
		class.Id = classesIds[i]
	}

	return nil
}

func GetScheduleRevisionClasses(scheduleRevisionId string) ([]*Class, error) {
	db, err := Connect()
	defer db.Close()

	var classes []*Class
	if err != nil {
		return classes, err
	}

	query := `
		SELECT 
		    id, 
		    date, 
		    start_hour, 
		    end_hour, 
		    name, 
		    lecturer, 
		    group_number, 
		    class_number 
		FROM class c
		WHERE schedule_revision_id = $1
	`

	rows, err := db.Query(query, scheduleRevisionId)
	defer rows.Close()

	if err != nil {
		return classes, err
	}
	for rows.Next() {
		var class Class
		err := rows.Scan(
			&class.Id,
			&class.Date,
			&class.StartHour,
			&class.EndHour,
			&class.Name,
			&class.Lecturer,
			&class.Group,
			&class.ClassNumber,
		)
		if err != nil {
			return nil, err
		}
		classes = append(classes, &class)
	}

	return classes, nil
}
