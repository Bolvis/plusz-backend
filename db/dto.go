package db

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
