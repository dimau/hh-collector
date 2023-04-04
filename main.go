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

	// Initialize RabbitMQ connection
	rabbitConn, rabbitChannel, rabbitHHQueue := initializeRabbitMQConnection(rabbitMQUser, rabbitMQPass, rabbitMQServerName, rabbitMQPort)
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	// Initialize HeadHunter connection
	client := initializeHHClient(appAccessToken)

	// Get vacancies page by page
	for page := 0; page < 3; page++ {
		vacancies := getVacancies(client, "react", page)
		publishVacanciesToRabbitMQ(rabbitChannel, rabbitHHQueue.Name, vacancies)
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func getVacancies(client *hh.Client, vacanciesType string, page int) *hh.Vacancies {
	var text string

	switch vacanciesType {
	case "react":
		text = "react"
	case "angular":
		text = "angular"
	case "vue":
		text = "vue"
	default:
		text = ""
	}

	// Options for getting React vacancies
	var options = &hh.OptionsForGetVacancies{
		Text:         text,
		SearchField:  "name",
		Period:       2,
		ItemsPerPage: 6,
		PageNumber:   page,
	}

	// Get vacancies
	vacancies, err := client.GetVacancies(options)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return vacancies
}

// Publish all vacancies to RabbitMQ queue with name HeadHunter
func publishVacanciesToRabbitMQ(rabbitChannel *amqp.Channel, routingKey string, vacancies *hh.Vacancies) {
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

// Creating a new REST API client for hh.ru
func initializeHHClient(appAccessToken string) *hh.Client {
	client := hh.NewClient(
		&url.URL{
			Scheme: "https",
			Host:   "api.hh.ru",
		},
		"dimau-app/1.0 (dimau777@gmail.com)",
		appAccessToken)

	return client
}
