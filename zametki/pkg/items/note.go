package items

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

type Note struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type NoteRepoInterface interface {
	NewNote(text string) error
	GetNotes() ([]Note, error)
	DeleteNote(id int) error
}

type NoteRepo struct {
	db  *sql.DB
	Mut sync.Mutex
}

func NewNoteRepo(db *sql.DB) NoteRepoInterface {
	repo := new(NoteRepo)
	repo.db = db
	return repo
}

func (r *NoteRepo) NewNote(text string) error {
	result, err := r.db.Exec(
		"INSERT INTO notes (`text`) VALUES (?)",
		text,
	)
	if err != nil {
		fmt.Println("err new note: ", err)
		return err
	}

	affected, err := result.RowsAffected()
	if err != nil {
		fmt.Println("err new note RowsAffected: ", err)
		return err

	}

	lastID, err := result.LastInsertId()
	if err != nil {
		fmt.Println("errnew note LastInsertId: ", err)
		return err
	}

	fmt.Println("Insert - RowsAffected", affected, "LastInsertId: ", lastID)

	return nil
}

func (r *NoteRepo) GetNotes() ([]Note, error) {
	rows, err := r.db.Query("SELECT * FROM notes")
	if err != nil {
		fmt.Println("get note error: ", err)
		return nil, err
	}

	var notes []Note

	for rows.Next() {

		var note = Note{}
		err = rows.Scan(&note.ID, &note.Text)
		if err != nil {
			fmt.Println("get events error = ", err)
			return nil, err
		}
		notes = append(notes, note)
	}
	// надо закрывать соединение, иначе будет течь
	rows.Close()
	return notes, nil
}

func (r *NoteRepo) DeleteNote(id int) error {
	_, err := r.db.Query("DELETE FROM notes WHERE id = " + strconv.Itoa(id))
	return err
}

func JSONError(w http.ResponseWriter, status int, msg string) {
	resp, err := json.Marshal(map[string]interface{}{
		"error": msg,
	})
	w.WriteHeader(status)
	if err != nil {
		fmt.Println("error in JSONError ")
		return
	}
	_, err2 := w.Write(resp)
	if err2 != nil {
		fmt.Println("some bad in JSONError write response")
	}
}
