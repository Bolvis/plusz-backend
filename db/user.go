package db

import "fmt"

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
