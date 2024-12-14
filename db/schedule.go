package db

type Schedule struct {
	Id                string             `json:"id"`
	ScheduleRevisions []ScheduleRevision `json:"schedule_revisions"`
}

func InsertSchedules(schedules []Schedule) error {

	return nil
}
