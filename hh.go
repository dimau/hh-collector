package main

import (
	"fmt"
	"github.com/dimau/hh-api-client-go"
	"log"
	"net/url"
	"time"
)

func getVacancies(client *hh.Client, vacanciesType string, page int) *[]hh.Vacancy {
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
	vacancyPublishTime, err := convertTimeStringHHToGoTime(vacancy.PublishedAt)
	failOnError(err, "Fail to convert vacancy publish time to time.Time in Go")
	if vacancyPublishTime.Compare(*lastHandledVacancyPublishTime) != 1 {
		return true
	} else {
		return false
	}
}

func getLastTime(lastHandledVacancyPublishTime *time.Time, filteredVacancies *[]hh.Vacancy) (*time.Time, error) {
	currentLastTime := lastHandledVacancyPublishTime
	for _, vacancy := range *filteredVacancies {
		vacancyTime, err := convertTimeStringHHToGoTime(vacancy.PublishedAt)
		if err != nil {
			return currentLastTime, err
		}

		if currentLastTime.Compare(*vacancyTime) == -1 {
			currentLastTime = vacancyTime
		}
	}

	return currentLastTime, nil
}

// Publish time in vacancies from HH has format: "2023-04-03T15:16:31+0300"
// RFC3339 required format: "2006-01-02T15:04:05+03:00"
func convertTimeStringHHToGoTime(source string) (*time.Time, error) {
	timeStringRFC3339 := fmt.Sprintf("%v:%v", source[0:22], source[22:24])
	goTime, err := time.Parse(time.RFC3339, timeStringRFC3339)
	return &goTime, err
}
