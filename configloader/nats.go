package configloader

import (
	"log"

	"github.com/nats-io/nats.go"
)

func NatsConnect(natsURL string) *nats.Conn {
	conn, err := nats.Connect(natsURL)
	defer conn.Close()

	if err != nil {
		log.Fatalf("nats connection failed: %v", err)
	}

	return conn
}

func JetStreamConnect(nats *nats.Conn) *nats.JetStreamContext {
	js, err := nats.JetStream()
	if err != nil {
		log.Fatalf("nats jet stream connection failed: %v", err)
	}

	return &js
}
