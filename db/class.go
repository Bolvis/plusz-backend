package db

import (
	"plusz-backend/util"
	"strings"
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

	insertClasses := `INSERT INTO public.class (date, start_hour, end_hour, name, lecturer, group_number, class_number) VALUES `
	var classesVals []interface{}

	for _, v := range classes {
		insertClasses += "(?::date,?::time,?::time,?::varchar,?::varchar,?::varchar,?::varchar),"
		classesVals = append(classesVals, v.Date, v.StartHour, v.EndHour, v.Name, v.Lecturer, v.Group, v.ClassNumber)
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

	insertRelations := `INSERT INTO public.schedule_revision_class_relation (class_id, schedule_revision_id) VALUES `
	var relationsVals []interface{}
	for _, v := range classes {
		insertRelations += "(?::integer,?::integer),"
		relationsVals = append(relationsVals, v.Id, scheduleRevisionId)
	}
	insertRelations = strings.TrimSuffix(insertRelations, ",")
	insertRelations = util.ReplaceSQL(insertRelations, "?")
	if _, err := db.Exec(insertRelations, relationsVals...); err != nil {
		return err
	}

	return nil
}
