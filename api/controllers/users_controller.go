package controllers

import (
	"../auth"
	"../models"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (server *Server) CreateUser(c *gin.Context) {
	user := models.User{}
	err := c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	err = user.ValidateUser("")
	saveUser, err := user.SaveUser(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"data":    saveUser,
		"message": "Create successfully",
	})
}

func (server *Server) GetAllUser(c *gin.Context) {
	user := models.User{}
	findAllUsers, err := user.FindAllUser(server.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": findAllUsers})
}

func (server *Server) GetUserById(c *gin.Context) {
	id := c.Params.ByName("id")
	uId := uuid.MustParse(id)
	user := models.User{}
	userById, err := user.FindUserById(server.DB, uId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"message": `Data By Id ` + uId.String() + ` Is Found`,
		"data":    userById,
	})
}

func (server *Server) UpdateUserById(c *gin.Context) {
	id := c.Params.ByName("id")
	uId := uuid.MustParse(id)
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, " + err.Error()})
		return
	}
	user := models.User{}
	_, err = user.FindUserById(server.DB, uId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	err = c.BindJSON(&user)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	userToken, err := auth.ExtractTokenUser(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	if userToken.ID != uId {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errors.New("your user id not the same like post author id " + uId.String() + " " + userToken.ID.String())})
		return
	}
	userById, err := user.UpdateUserById(server.DB, uId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"message": `Data By Id ` + uId.String() + ` Is Found`,
		"data":    userById,
	})
}

func (server *Server) DeleteUserById(c *gin.Context) {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, " + err.Error()})
		return
	}
	id := c.Params.ByName("id")
	uId := uuid.MustParse(id)
	user := models.User{}
	userToken, err := auth.ExtractTokenUser(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	if userToken.ID != uId {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errors.New("your user id not the same like post author id " + uId.String() + " " + userToken.ID.String())})
		return
	}
	deleteUserById, err := user.SoftDeleteUserById(server.DB, uId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusFound, gin.H{
		"message": `Delete Data By Id ` + id + ` Successfully`,
		"data":    deleteUserById,
	})
}
