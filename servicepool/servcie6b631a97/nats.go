package servcie6b631a97

import (
	"fmt"
	"log"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	serviceID string
	bucket    string
	Conn      *nats.Conn
}

func NewNatsClient(srv, bucket, natsUrl string) (NatsClient, error) {
	conn, err := nats.Connect(natsUrl)
	if err != nil {
		return NatsClient{}, fmt.Errorf("nats connection failed: %v", err)
	}

	return NatsClient{
		serviceID: srv,
		bucket:    bucket,
		Conn:      conn,
	}, nil
}

func (nc NatsClient) SubscribeToKV(notify chan []byte) {
	jsc, err := nc.Conn.JetStream()
	if err != nil {
		log.Fatalf("nats jet stream connection failed: %v", err)
	}

	kv, err := jsc.KeyValue(nc.bucket)
	if err != nil {
		log.Fatalf("unable to subscribe to KV storage: %v", err)
	}
	// Start a watcher for the specified key
	watch, err := kv.Watch(nc.serviceID)
	if err != nil {
		log.Fatal(err)
	}
	defer watch.Stop()

	for {
		update := <-watch.Updates()
		if update != nil {
			notify <- update.Value()
			log.Printf("config update was received and notified")
		}
	}
}
