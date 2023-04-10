module github.com/dimau/hh-collector

go 1.20

require (
	github.com/dimau/hh-api-client-go v0.0.0-20230407182029-8edbfd066ea9
	github.com/ijustfool/docker-secrets v0.0.0-20191021062307-b25ea5007562
	github.com/jackc/pgx/v5 v5.3.1
	github.com/rabbitmq/amqp091-go v1.8.0
)

require (
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	golang.org/x/crypto v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
)

//replace github.com/dimau/hh-api-client-go => ../hh-api-client-go
