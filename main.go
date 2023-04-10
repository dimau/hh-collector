package main

import (
	"github.com/ijustfool/docker-secrets"
	"log"
	"os"
	"time"
)

func main() {
	// Get environment vars for the application
	rabbitMQServerName := os.Getenv("RABBIT_MQ_SERVER_NAME")
	rabbitMQPort := os.Getenv("RABBIT_MQ_PORT")
	postgresServerName := os.Getenv("POSTGRES_SERVER_NAME")
	postgresPort := os.Getenv("POSTGRES_PORT")
	postgresDB := os.Getenv("POSTGRES_DB")

	// Get Docker secrets
	dockerSecrets, _ := secrets.NewDockerSecrets("")
	rabbitMQUser, _ := dockerSecrets.Get("rabbitmq_user")
	rabbitMQPass, _ := dockerSecrets.Get("rabbitmq_passwd")
	postgresUser, _ := dockerSecrets.Get("postgres_user")
	postgresPass, _ := dockerSecrets.Get("postgres_passwd")
	appAccessToken, _ := dockerSecrets.Get("hh_api_token")

	// Initialize RabbitMQ connection
	rabbitConn, rabbitChannel, rabbitHHQueue := initializeRabbitMQConnection(rabbitMQUser, rabbitMQPass, rabbitMQServerName, rabbitMQPort)
	defer rabbitConn.Close()
	defer rabbitChannel.Close()

	// Initialize Postgres connection
	db := initializePostgresConnection(postgresUser, postgresPass, postgresServerName, postgresPort, postgresDB)
	defer db.Close()

	lastHandledVacancyPublishTime := getLastHandledVacancyPublishTime(db)

	// Initialize HeadHunter connection
	client := initializeHHClient(appAccessToken)

	// Get vacancies page by page
	lastPublishVacancyTime := lastHandledVacancyPublishTime
	var err error

	for page := 0; page < 3; page++ {
		vacancies := getVacancies(client, "react", page, lastHandledVacancyPublishTime)
		filteredVacancies := filterAlreadyHandledVacancies(vacancies, lastHandledVacancyPublishTime)
		lastPublishVacancyTime, err = getLastTime(lastPublishVacancyTime, filteredVacancies)
		failOnError(err, "Fail on getting last time for filtered vacancies")
		publishVacanciesToRabbitMQ(rabbitChannel, rabbitHHQueue.Name, vacancies)

		// If we find at least one already handled vacancy, then stop getting new vacancies from API
		if len(*vacancies) != len(*filteredVacancies) {
			break
		}

		// Delay for 1 second to prevent overwhelming HeadHunter API
		t := time.NewTimer(1 * time.Second)
		<-t.C
	}

	// Save the last publish vacancy time in DB
	saveLastHandledVacancyPublishTime(db, lastPublishVacancyTime)
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
