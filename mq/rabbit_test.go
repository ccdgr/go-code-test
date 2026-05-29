package mq

import (
	"log"
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func TestConn(t *testing.T) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("close connection failed: %s", err)
		}
	}()
}
