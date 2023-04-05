package main

import (
	"database/sql"
	"fmt"
	_ "github.com/jackc/pgx/v5/stdlib"
	"time"
)

func initializePostgresConnection(postgresUserName, postgresPasswd, postgresServerName, postgresPort, postgresDB string) *sql.DB {
	connectionURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v", postgresUserName, postgresPasswd, postgresServerName, postgresPort, postgresDB)
	db, err := sql.Open("pgx", connectionURL)
	failOnError(err, "Failed to create a connection to Postgres")
	return db
}

func saveLastHandledVacancyPublishTime(db *sql.DB, lastHandledPublishTime *time.Time) {
	stmt, err := db.Prepare(`UPDATE public."LastHandledEntry" SET last_handled_entry = $1 WHERE source = 'HeadHunter';`)
	failOnError(err, "Fail on preparing statement for insert new last handled vacancy publish time to DB")

	lastHandledEntry, err := lastHandledPublishTime.MarshalText()
	failOnError(err, "Fail on marshalling lastHandledPublishTime to []byte")

	_, err = stmt.Exec(string(lastHandledEntry))
	failOnError(err, "Fail on updating last handled entry for source HeadHunter")
}

func getLastHandledVacancyPublishTime(db *sql.DB) *time.Time {
	stmt, err := db.Prepare(`SELECT last_handled_entry from public."LastHandledEntry" where source = $1`)
	failOnError(err, "Fail on preparing statement for getting last handled entry for HeadHunter")
	defer stmt.Close()

	var lastHandledEntry string
	err = stmt.QueryRow("HeadHunter").Scan(&lastHandledEntry)
	failOnError(err, "Fail on getting last handled entry for HeadHunter")

	lastHandledVacancyPublishTime, err := time.Parse(time.RFC3339, lastHandledEntry)
	failOnError(err, "Fail to convert last_handled_entry for HeadHunter to time.Time in Go")

	return &lastHandledVacancyPublishTime
}
