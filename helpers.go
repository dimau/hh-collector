package main

import (
	"fmt"
	"time"
)

// Publish time in vacancies from HH has format ISO8601: "2023-04-03T15:16:31+0300"
// RFC3339 required format: "2006-01-02T15:04:05+03:00"
func convertTimeFromISO8601StringToGoTime(source string) (*time.Time, error) {
	timeStringRFC3339 := fmt.Sprintf("%v:%v", source[0:22], source[22:24])
	goTime, err := time.Parse(time.RFC3339, timeStringRFC3339)
	return &goTime, err
}

// Publish time in vacancies from HH has format ISO8601: "2023-04-03T15:16:31+0300"
func convertGoTimeToISO8601(t *time.Time) string {
	return t.Format("2006-01-02T15:04:05-0700")
}
