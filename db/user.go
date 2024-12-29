package db

import (
	"fmt"
)

type User struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func AuthUser(user *User) error {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		return err
	}

	if err = db.QueryRow(`SELECT id FROM "user" WHERE login = $1 AND password = $2`, user.Login, user.Password).Scan(&user.Id); err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func AssignUserSchedule(user User, schedule Schedule) error {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		return err
	}

	id := ""
	if err = db.QueryRow(`SELECT id FROM user_schedule_relation WHERE user_id = $1 AND schedule_id = $2`, user.Id, schedule.Id).Scan(&id); err != nil {
		if err = db.QueryRow(`INSERT INTO user_schedule_relation (user_id, schedule_id) VALUES ($1, $2) RETURNING id`, user.Id, schedule.Id).Err(); err != nil {
			return err
		}

		fmt.Println("Assign schedule(", schedule.Id, ") to user(", user.Id, ")")
	} else {
		fmt.Println("User(", user.Id, ") already assigned to schedule(", schedule.Id, ")")
	}

	return nil
}

func GetUserSchedules(user User) ([]Schedule, error) {
	db, err := Connect()
	defer db.Close()

	var schedules []Schedule
	if err != nil {
		return schedules, err
	}

	query := `
		SELECT
			s.id,
			s.year,
			s.field,
			s.academic_year,
			s.semester
		FROM schedule s
				 LEFT OUTER JOIN user_schedule_relation usr ON s.id = usr.schedule_id
		WHERE $1 = usr.user_id
 	`

	rows, err := db.Query(query, user.Id)
	defer rows.Close()
	if err != nil {
		return schedules, err
	}

	for rows.Next() {
		var schedule Schedule
		err := rows.Scan(
			&schedule.Id,
			&schedule.Year,
			&schedule.Field,
			&schedule.AcademicYear,
			&schedule.Semester,
		)
		if err != nil {
			return schedules, err
		}
		fmt.Println(schedule)
		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
