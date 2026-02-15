package service

import (
	"fmt"
	"net/http"
	"plusz-backend/db"

	"github.com/gin-gonic/gin"
)

type noteRequest struct {
	NoteBody string `json:"noteBody"`
	ClassId  string `json:"classId"`
}

func AddNote(c *gin.Context) {
	userId := c.MustGet("UserId").(string)

	var request noteRequest
	if err := c.BindJSON(&request); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	note := db.Note{ClassId: request.ClassId, NoteBody: request.NoteBody, AuthorId: userId}
	var err error
	note, err = db.InsertNote(note)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": note.Id})
}

func GetNote(c *gin.Context) {
	userId := c.MustGet("UserId").(string)
	classId := c.Param("classId")

	note, err := db.ReadNote(userId, classId)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, note)
}
