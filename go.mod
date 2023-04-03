module github.com/dimau/hh-collector

go 1.20

require (
	github.com/dimau/hh-api-client-go v0.0.0-20230403180624-ce9f733b7a47
	github.com/ijustfool/docker-secrets v0.0.0-20191021062307-b25ea5007562
	github.com/rabbitmq/amqp091-go v1.8.0
)

require github.com/mitchellh/mapstructure v1.1.2 // indirect

//replace github.com/dimau/hh-api-client-go => ../hh-api-client-go
