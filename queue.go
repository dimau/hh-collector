package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dimau/hh-api-client-go"
	amqp "github.com/rabbitmq/amqp091-go"
	"time"
)

func initializeRabbitMQConnection(rabbitMQUser, rabbitMQPass, rabbitMQServerName, rabbitMQPort string) (*amqp.Connection, *amqp.Channel, amqp.Queue) {
	// Create connection to RabbitMQ
	connectionURL := fmt.Sprintf("amqp://%v:%v@%v:%v/", rabbitMQUser, rabbitMQPass, rabbitMQServerName, rabbitMQPort)
	conn, err := amqp.Dial(connectionURL)
	failOnError(err, "Failed to connect to RabbitMQ")

	// Create channel to RabbitMQ
	rabbitChannel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	// Create HeadHunter queue (if it's already exist just get this channel)
	args := make(amqp.Table)
	args["x-message-ttl"] = int32(86400000)
	q, err := rabbitChannel.QueueDeclare(
		"HeadHunter", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		args,         // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return conn, rabbitChannel, q
}

// Publish all vacancies to RabbitMQ queue with name HeadHunter
func publishVacanciesToRabbitMQ(rabbitChannel *amqp.Channel, routingKey string, vacancies *[]hh.Vacancy) {
	for _, vacancy := range *vacancies {
		vacancyMarshaled, err := json.Marshal(vacancy)
		failOnError(err, "Error when marshal vacancy to json string")

		// Create a context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Publish received vacancy to the queue HeadHunter in RabbitMQ
		err = rabbitChannel.PublishWithContext(ctx,
			"",         // exchange
			routingKey, // routing key
			false,      // mandatory
			false,      // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        vacancyMarshaled,
			},
		)
		failOnError(err, "Failed to publish a message")
	}
}
