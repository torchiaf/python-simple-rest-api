module test

go 1.23.1

replace github.com/torchiaf/Sensors/adapter => ../adapter

require github.com/torchiaf/Sensors/adapter v0.0.0-00010101000000-000000000000

require (
	github.com/fatih/structs v1.1.0 // indirect
	github.com/itchyny/gojq v0.12.17 // indirect
	github.com/itchyny/timefmt-go v0.1.6 // indirect
	github.com/rabbitmq/amqp091-go v1.10.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)
