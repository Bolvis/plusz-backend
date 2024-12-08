package db

import (
	"strings"

	"plusz-backend/util"
)

type Class struct {
	Id          string `json:"id"`
	Date        string `json:"date"`
	Hour        string `json:"hour"`
	Name        string `json:"name"`
	Lecturer    string `json:"lecturer"`
	Group       string `json:"group"`
	ClassNumber string `json:"classNumber"`
}

type ScheduleRevision struct {
	Id      string  `json:"id"`
	Date    string  `json:"date"`
	Classes []Class `json:"classes"`
}

type Schedule struct {
	Id                string             `json:"id"`
	ScheduleRevisions []ScheduleRevision `json:"schedule_revisions"`
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
		hours := strings.Split(strings.Replace(v.Hour, ".", ":", 2), "-")
		vals = append(vals, v.Date, hours[0], hours[1], v.Name, v.Lecturer, v.Group, v.ClassNumber)
	}
	query = strings.TrimSuffix(query, ",")
	query = util.ReplaceSQL(query, "?")
	stmt, _ := db.Prepare(query)
	if _, err := stmt.Exec(vals...); err != nil {
		return err
	}

	return nil
}

func InsertSchedules(schedules []Schedule) error {

	return nil
}

func InsertScheduleRevisions(scheduleRevisions []ScheduleRevision) error {

	return nil
}
