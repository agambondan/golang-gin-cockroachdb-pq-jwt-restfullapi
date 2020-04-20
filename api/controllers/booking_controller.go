package controllers

import (
	"../models"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"net/http"
)

func (server *Server) SendMessage(c *gin.Context) {
	var booking models.Booking
	booking.Code = uuid.New().String()
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(c, http.StatusInternalServerError, err)
	defer conn.Close()
	err = c.Bind(&booking)
	failOnError(c, http.StatusUnprocessableEntity, err)
	channel, err := conn.Channel()
	failOnError(c, http.StatusInternalServerError, err)
	queue, err := channel.QueueDeclare(
		"booking",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(c, http.StatusInternalServerError, err)
	bytes, err := json.Marshal(&booking)
	failOnError(c, http.StatusUnprocessableEntity, err)
	err = channel.Publish(
		"notifyExchange",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(string(bytes)),
		})
	failOnError(c, http.StatusUnprocessableEntity, err)
	c.JSON(http.StatusOK, gin.H{"data": booking})
}
