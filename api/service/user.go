package service

import (
	"net/http"

	"plusz-backend/api/authorization"
	"plusz-backend/db"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func RegisterUser(c *gin.Context) {
	var user db.User

	if err := c.BindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	passwordByte, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = string(passwordByte)

	id, err := db.InsertUser(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"id": id})
}

func AuthenticateUser(c *gin.Context) {
	var inputUser db.User
	if err := c.BindJSON(&inputUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	dbUser, err := db.GetUserByLogin(inputUser.Login)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(inputUser.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	token, err := authorization.GenerateToken(dbUser.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"sessionToken": token})
}
