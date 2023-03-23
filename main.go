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

	client := hh.NewClient(
		&url.URL{
			Scheme: "https",
			Host:   "api.hh.ru",
		},
		"dimau-app/1.0 (dimau777@gmail.com)",
		*appAccessToken)

	appInfo, err := client.Me()
	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	fmt.Printf("App info: %+v", appInfo)
}
