package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"plusz-backend/db"
)

type noteRequest struct {
	NoteBody string  `json:"noteBody"`
	ClassId  string  `json:"classId"`
	User     db.User `json:"user"`
}

func AddNote(c *gin.Context) {
	var request noteRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := db.AuthUser(&request.User); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	note := db.Note{ClassId: request.ClassId, NoteBody: request.NoteBody, AuthorId: request.User.Id}

	note, err := db.InsertNote(note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": note.Id})
}

func GetNote(c *gin.Context) {
	var err error
	classId := c.Query("classId")

	user := db.User{Login: c.Query("login"), Password: c.Query("password")}
	if err = db.AuthUser(&user); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	var note db.Note
	note, err = db.ReadNote(user.Id, classId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}
