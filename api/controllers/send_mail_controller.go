package controllers

import (
	"../auth"
	"fmt"
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"net/http"
	"os"
	"path/filepath"
)

const ConfigSmtpHost = "smtp.gmail.com"
const ConfigSmtpPort = 465
const ConfigEmail = "agam.pro234@gmail.com"
const ConfigPassword = "selamatagam0"

func (server *Server) SendEmail(c *gin.Context) {
	filenames, err := uploadFile(c)
	if len(filenames) > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Upload file can't greater than 5"})
	}
	to := c.PostForm("to")
	cc := c.PostForm("cc")
	ccName := c.PostForm("cc_name")
	subject := c.PostForm("subject")
	body := c.PostForm("body")
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", ConfigEmail)
	mailer.SetHeader("To", to)
	mailer.SetAddressHeader("Cc", cc, ccName)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)
	mailer.Attach(filenames[1])
	mailer.Attach(filenames[2])
	mailer.Attach(filenames[3])
	mailer.Attach(filenames[4])
	mailer.Attach(filenames[5])
	dialer := gomail.NewDialer(
		ConfigSmtpHost,
		ConfigSmtpPort,
		ConfigEmail,
		ConfigPassword,
	)
	err = dialer.DialAndSend(mailer)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"message": "Mail sent to" + "" + "!"})
}

func uploadFile(c *gin.Context) (filenames []string, err error) {
	err = auth.TokenValid(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized"})
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
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	files := form.File["files"]
	for _, file := range files {
		basename := filepath.Base(file.Filename)
		dir := filepath.Join("./assets/mail/", userToken.ID.String())
		if dir != "" {
			err := os.Mkdir("./assets/mail/"+userToken.ID.String(), os.ModePerm)
			if err != nil {
				fmt.Println(err.Error())
			}
		}
		regex := after(basename, ".")
		if regex != "png" && regex != "jpg" {
			c.JSON(http.StatusBadRequest, gin.H{"message": "file must be image png or jpg"})
			return
		}
		filename := filepath.Join("./assets/mail/", userToken.ID.String(), basename)
		err := c.SaveUploadedFile(file, filename)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error(), "test": 1})
			return filenames, err
		}
		for _, file := range files {
			filenames = append(filenames, file.Filename)
		}
	}
	return
}
