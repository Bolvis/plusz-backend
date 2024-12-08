package db

import (
	"plusz-backend/util"
	"strings"
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
	if err != nil {
		return err
	}

	query := "INSERT INTO public.class (id) VALUES "
	vals := []interface{}{}
	//INSERT INTO public.class (date, start_hour, end_hour, name, lecturer, "group", class_number)
	//VALUES ('2024-12-13', '15:16:29', '18:16:00', 'test', 'test', 'test',
	//	'test')

	for _, v := range classes {
		query += "(?::date,?::time,?::time,?::varchar,?::varchar,?::varchar,?::varchar),"
		hours := strings.Split(v.Hour, "-")
		vals = append(vals, v.Date, hours[0], hours[1], v.Name, v.Lecturer, v.Group, v.ClassNumber)
	}
	query = strings.TrimSuffix(query, ",")
	query = util.ReplaceSQL(query, "?")
	stmt, _ := db.Prepare(query)
	if _, err := stmt.Exec(vals...); err != nil {
		return err
	}

	defer db.Close()

	return nil
}

func InsertSchedules(schedules []Schedule) error {

	return nil
}

func InsertScheduleRevisions(scheduleRevisions []ScheduleRevision) error {

	return nil
}
