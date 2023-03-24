package main

import (
	"flag"
	"fmt"
	"github.com/dimau/go-hh-client"
	"net/url"
)

func main() {
	// Get parameters of launching for application
	appAccessToken := flag.String("APP_ACCESS_TOKEN", "", "Access token for application registered in hh.ru")
	flag.Parse()

	// Creating a new REST API client for hh.ru
	client := hh.NewClient(
		&url.URL{
			Scheme: "https",
			Host:   "api.hh.ru",
		},
		"dimau-app/1.0 (dimau777@gmail.com)",
		*appAccessToken)

	// Test information about app
	appInfo, err := client.Me()
	if err != nil {
		fmt.Printf("Error: %v", err)
	} else {
		fmt.Printf("App info: %+v", appInfo)
	}

	// Get vacancies
	var options = &hh.OptionsForGetVacancies{
		Text:         "react",
		SearchField:  "name",
		Period:       2,
		ItemsPerPage: 5,
		PageNumber:   0,
	}

	vacancies, err := client.GetVacancies(options)
	if err != nil {
		fmt.Printf("Error: %v", err)
	} else {
		fmt.Printf("Vacancies: %+v", vacancies)
	}
}
