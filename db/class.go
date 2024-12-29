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

func InsertClasses(classes []Class) error {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		return err
	}

	query := `INSERT INTO public.class (date, start_hour, end_hour, name, lecturer, group_number, class_number) VALUES `
	var vals []interface{}

	for _, v := range classes {
		query += "(?::date,?::time,?::time,?::varchar,?::varchar,?::varchar,?::varchar),"
		vals = append(vals, v.Date, v.StartHour, v.EndHour, v.Name, v.Lecturer, v.Group, v.ClassNumber)
	}
	query = strings.TrimSuffix(query, ",")
	query = util.ReplaceSQL(query, "?")
	stmt, _ := db.Prepare(query)
	if _, err := stmt.Exec(vals...); err != nil {
		return err
	}

	return nil
}
