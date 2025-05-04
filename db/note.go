package db

import (
	"database/sql"
	"errors"
	"fmt"
)

type Note struct {
	Id       string `json:"id"`
	ClassId  string `json:"classId"`
	AuthorId string `json:"authorId"`
	NoteBody string `json:"noteBody"`
}

func InsertNote(note Note) (Note, error) {
	db, err := Connect()
	defer db.Close()

	if err != nil {
		fmt.Println(err)
		return note, err
	}

	err = db.QueryRow(
		`SELECT id FROM note WHERE author_id = $1 AND class_id = $2`,
		note.AuthorId,
		note.ClassId,
	).Scan(&note.Id)

	if errors.Is(err, sql.ErrNoRows) {
		insertQuery := `INSERT INTO note (class_id, author_id, note_body) VALUES ($1, $2, $3) RETURNING id`
		stmt, err := db.Prepare(insertQuery)
		defer stmt.Close()
		if err != nil {
			fmt.Println(err)
			return note, err
		}

		if err = stmt.QueryRow(note.ClassId, note.AuthorId, note.NoteBody).Scan(&note.Id); err != nil {
			return note, err
		}

		return note, nil
	} else if err != nil {
		return note, err
	}

	return UpdateNote(note)
}

func UpdateNote(note Note) (Note, error) {
	db, err := Connect()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return note, err
	}

	err = db.QueryRow(
		`UPDATE note SET note_body = $1 WHERE id = $2 RETURNING id`,
		note.NoteBody,
		note.Id,
	).Err()

	if err != nil {
		fmt.Println(err)
		return note, err
	}

	return note, nil
}

func ReadNote(authorId string, classId string) (Note, error) {
	db, err := Connect()
	defer db.Close()
	if err != nil {
		fmt.Println(err)
		return Note{}, err
	}

	var note Note
	err = db.QueryRow(
		`SELECT id, note_body FROM note WHERE author_id = $1 AND class_id = $2`,
		authorId,
		classId,
	).Scan(&note.Id, &note.NoteBody)

	if errors.Is(err, sql.ErrNoRows) {
		fmt.Println(err)
		return Note{}, nil
	} else if err != nil {
		fmt.Println(err)
		return Note{}, err
	}

	note.ClassId = classId
	note.AuthorId = authorId

	return note, nil
}
