package controllers

import (
	"../auth"
	"../models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func (server *Server) CreateRole(c *gin.Context) {
	err := auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
		return
	}
	role := models.Role{}
	err = c.BindJSON(&role)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	userToken, err := auth.ExtractTokenUser(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, uid : " + userToken.ID.String() + err.Error()})
		return
	}
	if userToken.Role.Name != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, Your not admin, you role is " + userToken.Role.Name})
		return
	}
	err = role.Validate()
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	saveRole, err := role.SaveRole(server.DB)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": saveRole})
}

func (server *Server) GetAllRole(c *gin.Context) {
	role := models.Role{}
	findAllRole, err := role.FindAllRole(server.DB)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": findAllRole})
}

func (server *Server) GetRoleById(c *gin.Context) {
	id := c.Params.ByName("id")
	uId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	role := models.Role{}
	findRoleById, err := role.FindRoleById(server.DB, uId)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusFound, gin.H{"data": findRoleById})
}

func (server *Server) UpdateRoleById(c *gin.Context) {
	id := c.Params.ByName("id")
	uId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	role := models.Role{}
	err = c.BindJSON(&role)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	updateRoleById, err := role.UpdateRoleById(server.DB, uId)
	if err != nil {
		c.JSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updateRoleById})
}

func (server *Server) DeleteRoleById(c *gin.Context) {
	id := c.Params.ByName("id")
	uId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	role := models.Role{}
	softDeleteRoleById, err := role.SoftDeleteRoleById(server.DB, uId)
	if err != nil {
		c.JSON(http.StatusFound, gin.H{"message": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": softDeleteRoleById})

}
