package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/torchiaf/Sensors/controller/config"
	"github.com/torchiaf/Sensors/controller/models"
	"github.com/torchiaf/Sensors/controller/utils"
)

func FailOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func randomString(l int) string {
	bytes := make([]byte, l)
	for i := 0; i < l; i++ {
		bytes[i] = byte(randInt(65, 90))
	}
	return string(bytes)
}

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())

	return min + rand.Intn(max-min)
}

func Exec(routingKey string, message models.Message) (res string, err error) {

	address := fmt.Sprintf("amqp://%s:%s@%s:%s/", config.Config.RabbitMQ.Username, config.Config.RabbitMQ.Password, config.Config.RabbitMQ.Host, config.Config.RabbitMQ.Port)

	conn, err := amqp.Dial(address)
	FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"",    // name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // noWait
		nil,   // arguments
	)
	FailOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	FailOnError(err, "Failed to register a consumer")

	corrId := randomString(32)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	log.Printf("Msg: %s:", utils.ToString(message))

	err = ch.PublishWithContext(ctx,
		"",         // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:   "text/plain",
			CorrelationId: corrId,
			ReplyTo:       q.Name,
			Body:          []byte(utils.ToString(message)),
		})
	FailOnError(err, "Failed to publish a message")

	for d := range msgs {
		if corrId == d.CorrelationId {
			res = string(d.Body)
			FailOnError(err, "Error msgs")
			break
		}
	}

	return
}
