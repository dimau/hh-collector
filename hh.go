package main

import (
	"github.com/dimau/hh-api-client-go"
	"log"
	"net/url"
	"time"
)

func getVacancies(client *hh.Client, vacanciesType string, page int, fromPublishTime *time.Time) *[]hh.Vacancy {
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
		Period:       0,
		ItemsPerPage: 6,
		PageNumber:   page,
		DateFrom:     fromPublishTime,
		OrderBy:      "publication_time",
	}

	// Get vacancies
	vacancies, err := client.GetVacancies(options)
	if err != nil {
		log.Fatalln(err.Error())
	}

	return &vacancies.Items
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

func filterAlreadyHandledVacancies(vacancies *[]hh.Vacancy, lastHandledVacancyPublishTime *time.Time) *[]hh.Vacancy {
	res := []hh.Vacancy{}

	for _, vacancy := range *vacancies {
		if !isAlreadyHandledVacancy(&vacancy, lastHandledVacancyPublishTime) {
			res = append(res, vacancy)
		}
	}

	return &res
}

func isAlreadyHandledVacancy(vacancy *hh.Vacancy, lastHandledVacancyPublishTime *time.Time) bool {
	vacancyPublishTime, err := convertTimeFromISO8601StringToGoTime(vacancy.PublishedAt)
	failOnError(err, "Fail to convert vacancy publish time to time.Time in Go")
	if vacancyPublishTime.Compare(*lastHandledVacancyPublishTime) != 1 {
		return true
	} else {
		return false
	}
}

func getLastTime(lastHandledVacancyPublishTime *time.Time, vacancies *[]hh.Vacancy) (*time.Time, error) {
	currentLastTime := lastHandledVacancyPublishTime
	for _, vacancy := range *vacancies {
		vacancyTime, err := convertTimeFromISO8601StringToGoTime(vacancy.PublishedAt)
		if err != nil {
			return nil, err
		}

		if currentLastTime.Compare(*vacancyTime) == -1 {
			currentLastTime = vacancyTime
		}
	}

	return currentLastTime, nil
}
