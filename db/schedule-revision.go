package db

type ScheduleRevision struct {
	Id      string  `json:"id"`
	Date    string  `json:"date"`
	Classes []Class `json:"classes"`
}

func InsertScheduleRevisions(scheduleRevisions []ScheduleRevision) error {

	return nil
}
