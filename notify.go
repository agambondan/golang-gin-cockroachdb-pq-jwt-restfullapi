package main

import (
	"./api/models"
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func main() {
	ReceiveMessage()
}

func ReceiveMessage() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	failOnError(err, "Failed to connect AMQP")
	defer conn.Close()
	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	q, err := ch.QueueDeclare(
		"booking",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to declare a queue")
	msg, err := ch.Consume(
		q.Name, // queue
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	failOnError(err, "Failed to register a consumer")
	forever := make(chan bool)
	go func() {
		for d := range msg {
			msg := string(d.Body)
			booking := models.Booking{}
			err := json.Unmarshal([]byte(msg), &booking)
			failOnError(err, fmt.Sprintln(err))
			fmt.Println("Sending notification to user : " + booking.Username + " with booking code : " + booking.Code)
		}
	}()
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

//msgCount :=0
//go func() {
//	for d := range msg {
//
//		msgCount++
//
//		fmt.Printf("\nMessage Count: %d, Message Body: %s\n", msgCount, d.Body)
//
//	}
//}()
//
//select {
//case <-time.After(time.Second * 10):
//fmt.Printf("Total Messages Fetched: %d\n",msgCount)
//fmt.Println("No more messages in queue. Timing out...")
//
//}
