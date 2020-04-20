package controllers

import (
	"../models"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"log"
	"net/http"
)

func (server *Server) GetConsumeMessageBroker(c *gin.Context) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(c, http.StatusInternalServerError, err)
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(c, http.StatusInternalServerError, err)
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"booking",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(c, http.StatusInternalServerError, err)
	msg, err := ch.Consume(
		q.Name, // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(c, http.StatusInternalServerError, err)
	forever := make(chan bool)
	go func() {
		for d := range msg {
			msg := string(d.Body)
			booking := models.Booking{}
			err := json.Unmarshal([]byte(msg), &booking)
			failOnError(c, http.StatusInternalServerError, err)
			fmt.Println("Sending notification to user : " + booking.Username + " with booking code : " + booking.Code)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
