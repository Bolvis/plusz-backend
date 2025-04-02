package service

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"plusz-backend/db"
)

func RegisterUser(c *gin.Context) {
	var user db.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id, err := db.InsertUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}
