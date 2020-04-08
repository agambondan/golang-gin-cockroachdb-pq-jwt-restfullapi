package controllers

import (
	"../models"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"path/filepath"
)

func (server *Server) UploadFile(c *gin.Context) {
	userToken := extractToken(c)
	role := models.Role{}
	err := c.BindJSON(&role)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": err.Error()})
		return
	}
	if userToken.Role.Name != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, Your not admin, you role is " + userToken.Role.Name})
		return
	}
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	files := form.File["files"]
	for _, file := range files {
		basename := filepath.Base(file.Filename)
		regex := after(basename, ".")
		if regex == "png" || regex == "jpg" {
			dir := filepath.Join("./assets/images/", userToken.ID.String())
			if dir != "" {
				err := os.Mkdir("./assets/images/"+userToken.ID.String(), os.ModePerm)
				if err != nil {
					fmt.Println(err.Error())
				}
			}
			c.JSON(http.StatusBadRequest, gin.H{"message": "file must be image png or jpg"})
			return
		}
		filename := filepath.Join("./assets/images/", userToken.ID.String(), basename)
		err := c.SaveUploadedFile(file, filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "test": 1})
			return
		}
	}
	var filenames []string
	for _, file := range files {
		filenames = append(filenames, file.Filename)
	}
	c.JSON(http.StatusCreated, gin.H{"code": http.StatusAccepted, "message": "upload ok!", "data": gin.H{"files": filenames}})
}
