package main

import (
	"github.com/streadway/amqp"
	"log"
	"strconv"
)

var conn *amqp.Connection
var ch *amqp.Channel

func init() {
	var err error
	conn, err = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	ch, err = conn.Channel()
	failOnError(err, "Failed to open a channel")
}

func main() {
	defer conn.Close()
	defer ch.Close()
	multipleConsumers()
	log.Print("Done!")
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// If there are multiple consumers subscribed to the same queue,
// only one of the consumers receives a message at a time.
// This is regardless of whether autoAck is set to true or false.
func multipleConsumers() {

	autoAck := false

	// Queue
	q, err := ch.QueueDeclare("multipleConsumers", false, false, false, false, nil)
	failOnError(err, "Failed to declare 'multipleConsumers' queue")

	// Declare two consumers
	consumer1, err := ch.Consume(q.Name, "", autoAck, false, false, false, nil)
	failOnError(err, "Failed to register consumer1")
	go func() {
		for d := range consumer1 {
			log.Printf("Consumer1: Received a message: %s", d.Body)
			if string(d.Body) == "5" && !d.Redelivered {
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	consumer2, err := ch.Consume(q.Name, "", autoAck, false, false, false, nil)
	failOnError(err, "Failed to register consumer2")
	go func() {
		for d := range consumer2 {
			log.Printf("Consumer2: Received a message: %s", string(d.Body))
			if string(d.Body) == "5" && !d.Redelivered {
				d.Nack(false, true)
			} else {
				d.Ack(false)
			}
		}
	}()

	// Publish
	for i := 0; i < 10; i++ {
		err = ch.Publish("", q.Name, false, false, amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(strconv.Itoa(i)),
		})
		failOnError(err, "Failed to publish a message to "+q.Name)
	}

	forever := make(chan bool)
	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
