package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"zametki/pkg/items"
)

type NoteHandler struct {
	Repo items.NoteRepoInterface
}

func NewNoteHandler(db *sql.DB) *NoteHandler {
	hand := new(NoteHandler)
	hand.Repo = items.NewNoteRepo(db)
	return hand
}

func (h *NoteHandler) NewNote(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi new notes")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("err with body")
		items.JSONError(w, 500, "some bad in read data")
		return
	}
	req := &items.Note{}

	err1 := json.Unmarshal(body, req)
	if err1 != nil {
		fmt.Println("err unmarshal body")
		items.JSONError(w, 500, "some bad in unmarshal data")
		return
	}
	err = h.Repo.NewNote(req.Text)
	if err != nil {
		fmt.Println("err json marshal")
		items.JSONError(w, 500, "some bad in newNote")
		return
	}
}

func (h *NoteHandler) GetNotes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi get notes")

	notes, err := h.Repo.GetNotes()
	if err != nil {
		fmt.Println("err get notes")
		items.JSONError(w, 500, "some bad in marshal response")
		return
	}

	response, err := json.Marshal(notes)
	if err != nil {
		fmt.Println("err json marshal")
		items.JSONError(w, 500, "some bad in marshal response")
		return
	}
	w.Write(response)
}

func (h *NoteHandler) DeleteNotes(w http.ResponseWriter, r *http.Request) {
	fmt.Println("hi delete notes")
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("err with body")
		items.JSONError(w, 500, "some bad in read data")
		return
	}
	req := &items.Note{}

	err1 := json.Unmarshal(body, req)
	if err1 != nil {
		fmt.Println("err unmarshal body")
		items.JSONError(w, 500, "some bad in unmarshal data")
		return
	}
	err = h.Repo.DeleteNote(req.ID)
	if err != nil {
		fmt.Println("err json marshal")
		items.JSONError(w, 500, "some bad in deleteNote")
		return
	}
}
