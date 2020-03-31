package controllers

import (
	"../auth"
	"../models"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func (server *Server) UploadFile(c *gin.Context) {
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
	userToken, err := auth.ExtractTokenUser(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	if userToken.ID != uId {
		c.JSON(http.StatusUnauthorized, gin.H{"message": errors.New("your user id not the same like post author id " + uId.String() + " " + userToken.ID.String())})
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
		dir := filepath.Join("./assets/images/", userToken.ID.String())
		if dir != "" {
			err := os.Mkdir("./assets/images/"+userToken.ID.String(), os.ModePerm)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		regex := after(basename, ".")
		if regex != "png" && regex != "jpg" {
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

func after(value string, a string) string {
	// Get substring after a string.
	pos := strings.LastIndex(value, a)
	if pos == -1 {
		return ""
	}
	adjustedPos := pos + len(a)
	if adjustedPos >= len(value) {
		return ""
	}
	return value[adjustedPos:len(value)]
}
