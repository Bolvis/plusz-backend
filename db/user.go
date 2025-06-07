package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type User struct {
	Id       string `json:"id"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func GetUserByLogin(login string) (User, error) {
	var user User
	db, err := Connect()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return user, err
	}

	if err = db.QueryRow(`SELECT id, login, password FROM "user" WHERE login = $1`, login).Scan(&user.Id, &user.Login, &user.Password); err != nil {
		return user, err
	}

	return user, nil
}

func InsertUser(user User) (string, error) {
	db, err := Connect()
	defer db.Close()
	if err != nil {
		return "", err
	}

	if err = db.QueryRow(`SELECT id FROM "user" WHERE login = $1`, user.Login).Scan(); errors.Is(err, sql.ErrNoRows) {
		var id string
		if err = db.QueryRow(`INSERT INTO "user" (login, password) VALUES ($1, $2) RETURNING id`, user.Login, user.Password).Scan(&id); err != nil {
			return "", err
		}
		return id, nil
	}

	if err != nil {
		return "", err
	}

	return "", errors.New("user with the same login already exists")
}

func AssignUserSchedule(userId string, scheduleId string) error {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return err
	}

	id := ""
	if err = db.QueryRow(`SELECT id FROM user_schedule_relation WHERE user_id = $1 AND schedule_id = $2`, userId, scheduleId).Scan(&id); errors.Is(err, sql.ErrNoRows) {
		if err = db.QueryRow(`INSERT INTO user_schedule_relation (user_id, schedule_id) VALUES ($1, $2)`, userId, scheduleId).Err(); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	return nil
}

func RemoveUserScheduleAssignment(userId string, scheduleId string) error {
	db, err := Connect()
	if err != nil {
		return err
	}
	defer db.Close()

	if err = db.QueryRow(`DELETE FROM user_schedule_relation WHERE user_id = $1 AND schedule_id = $2`, userId, scheduleId).Err(); err != nil {
		return err
	}

	return nil
}

func GetUserSchedules(userId string) ([]Schedule, error) {
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

	rows, err := db.Query(query, userId)
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

		schedules = append(schedules, schedule)
	}

	return schedules, nil
}
