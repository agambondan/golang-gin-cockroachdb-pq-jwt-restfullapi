package controllers

import (
	"../auth"
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	//"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Login(c *gin.Context) {
	var user models.User
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}
	err = user.ValidateUser("login")
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	token, err := server.SignIn(user.Email, user.Username, user.Password)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{
			"message": "Username or Password is wrong ",
			"error":   err.Error(),
			"token":   token,
		})
		fmt.Println(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"username": user.Username,
		"password": user.Password,
		"email":    user.Email,
	})
}

func (server *Server) SignIn(email, username, password string) (string, error) {
	user := models.User{}
	err := server.DB.QueryRow("SELECT id, username, email, password FROM users WHERE username=$1 OR email=$2", username, email).Scan(&user.ID, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return err.Error(), err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println(err)
		return "", err
	}
	return auth.CreateToken(user.ID)
}
