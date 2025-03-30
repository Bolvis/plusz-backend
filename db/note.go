package db

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
		return note, err
	}
	insertQuery := `INSERT INTO note (class_id, author_id, note_body) VALUES ($1, $2, $3) RETURNING id`
	stmt, err := db.Prepare(insertQuery)
	defer stmt.Close()
	if err != nil {
		return note, err
	}

	if err = stmt.QueryRow(note.ClassId, note.AuthorId, note.NoteBody).Scan(&note.Id); err != nil {
		return note, err
	}

	return note, nil
}
