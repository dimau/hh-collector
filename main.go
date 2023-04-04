package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/dimau/hh-api-client-go"
	"github.com/ijustfool/docker-secrets"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
	"net/url"
	"os"
	"time"
)

func main() {
	// Get environment vars for the application
	rabbitMQServerName := os.Getenv("RABBIT_MQ_SERVER_NAME")
	rabbitMQPort := os.Getenv("RABBIT_MQ_PORT")

	// Get Docker secrets
	dockerSecrets, _ := secrets.NewDockerSecrets("")
	rabbitMQUser, _ := dockerSecrets.Get("rabbitmq_user")
	rabbitMQPass, _ := dockerSecrets.Get("rabbitmq_pass")
	appAccessToken, _ := dockerSecrets.Get("hh_api_token")

	// Create connection to RabbitMQ
	connectionURL := fmt.Sprintf("amqp://%v:%v@%v:%v/", rabbitMQUser, rabbitMQPass, rabbitMQServerName, rabbitMQPort)
	conn, err := amqp.Dial(connectionURL)
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	// Create channel to RabbitMQ
	rabbitChannel, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer rabbitChannel.Close()

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

	// Creating a new REST API client for hh.ru
	client := hh.NewClient(
		&url.URL{
			Scheme: "https",
			Host:   "api.hh.ru",
		},
		"dimau-app/1.0 (dimau777@gmail.com)",
		appAccessToken)

	// Options for getting React vacancies
	var options = &hh.OptionsForGetVacancies{
		Text:         "react",
		SearchField:  "name",
		Period:       2,
		ItemsPerPage: 5,
		PageNumber:   0,
	}

	// Get vacancies
	vacancies, err := client.GetVacancies(options)
	if err != nil {
		fmt.Printf("Error: %v", err)
	} else {
		fmt.Printf("Vacancies: %+v", vacancies)
	}

	// Publish all vacancies to RabbitMQ queue with name HeadHunter
	for _, vacancy := range vacancies.Items {
		vacancyMarshaled, err := json.Marshal(vacancy)
		if err != nil {
			log.Fatalln(err.Error())
		}

		// Create a context
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Publish received vacancy to the queue HeadHunter in RabbitMQ
		err = rabbitChannel.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "application/json",
				Body:        vacancyMarshaled,
			},
		)
		failOnError(err, "Failed to publish a message")
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
