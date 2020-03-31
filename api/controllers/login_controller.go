package controllers

import (
	"../auth"
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
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
	err = server.DB.QueryRow("SELECT id, full_name, username, email, password, role_id FROM users WHERE username=$1 OR email=$2", user.Username, user.Email).
		Scan(&user.ID, &user.FullName, &user.Username, &user.Email, &user.Password, &user.RoleId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{
		"token":    token,
		"user":     user,
	})
}

func (server *Server) SignIn(email, username, password string) (string, error) {
	user := models.User{}
	err := server.DB.QueryRow("SELECT id, full_name, username, email, password, role_id FROM users WHERE username=$1 OR email=$2", username, email).
		Scan(&user.ID, &user.FullName, &user.Username, &user.Email, &user.Password, &user.RoleId)
	if err != nil {
		return err.Error(), err
	}
	err = models.VerifyPassword(user.Password, password)
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword {
		fmt.Println(err)
		return "", err
	}
	role := models.Role{}
	roleById, err := role.FindRoleById(server.DB, user.RoleId)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	user.Role = *roleById
	return auth.CreateToken(user)
}
