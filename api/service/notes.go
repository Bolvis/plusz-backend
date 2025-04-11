package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"plusz-backend/api/authorization"
	"plusz-backend/db"
)

type noteRequest struct {
	NoteBody string `json:"noteBody"`
	ClassId  string `json:"classId"`
}

func AddNote(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	token, err := authorization.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var request noteRequest
	if err := c.BindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := db.Note{ClassId: request.ClassId, NoteBody: request.NoteBody, AuthorId: token.UserId}

	note, err = db.InsertNote(note)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": note.Id})
}

func GetNote(c *gin.Context) {
	tokenString := c.Request.Header.Get("Authorization")

	token, err := authorization.VerifyToken(tokenString)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	classId := c.Query("classId")

	var note db.Note
	note, err = db.ReadNote(token.UserId, classId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}
