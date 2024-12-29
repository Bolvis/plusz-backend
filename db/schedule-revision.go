package db

type ScheduleRevision struct {
	Id      string  `json:"id"`
	Date    string  `json:"date"`
	Classes []Class `json:"classes"`
}

func GetScheduleRevisionId(scheduleRevisions []ScheduleRevision) error {
	//db, err := Connect()
	//defer db.Close()
	//
	//if err != nil {
	//	return err
	//}
	//
	//query := `INSERT INTO public.class (date, start_hour, end_hour, name, lecturer, group_number, class_number) VALUES `

	return nil
}
