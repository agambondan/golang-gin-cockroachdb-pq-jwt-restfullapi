package controllers

import (
	"../auth"
	"../models"
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"strings"
)

func extractToken(c *gin.Context) (userToken *models.User) {
	var err error
	err = auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, " + err.Error()})
		return
	}
	userToken, err = auth.ExtractTokenUser(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
		return
	}
	return
}

func uploadFile(userToken *models.User, c *gin.Context) (filenames []string) {
	var err error
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	files := form.File["files"]
	for _, file := range files {
		basename := filepath.Base(file.Filename)
		filename := filepath.Join("./assets/mail/" + userToken.ID.String() + "-" + basename)
		err := c.SaveUploadedFile(file, filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
			return
		}
		filenames = append(filenames, filename)
	}
	return
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
