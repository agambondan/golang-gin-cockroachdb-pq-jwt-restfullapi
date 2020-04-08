package controllers

import (
	"github.com/gin-gonic/gin"
	"gopkg.in/gomail.v2"
	"net/http"
)

const ConfigSmtpHost = "smtp.gmail.com"
const ConfigSmtpPort = 465
const ConfigEmail = "agam.pro234@gmail.com"
const ConfigPassword = "selamatagam0"

func (server *Server) SendEmail(c *gin.Context) {
	var err error
	userToken := extractToken(c)
	if userToken.Role.Name != "admin" {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "UnAuthorized, Your not admin, you role is " + userToken.Role.Name})
		return
	}
	filenames := uploadFile(userToken, c)
	if len(filenames) > 5 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Upload file can't greater than 5 file"})
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
	for i := 0; i < len(filenames); i++ {
		mailer.Attach(filenames[i])
	}
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
	c.JSON(http.StatusOK, gin.H{"message": "Mail sent to" + to + "!"})
}

