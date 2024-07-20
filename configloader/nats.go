package configloader

import (
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
)

type NatsClient struct {
	Conn *nats.Conn
}

func NewNatsClient(natsUrl string) (NatsClient, error) {
	conn, err := nats.Connect(natsUrl)
	if err != nil {
		return NatsClient{}, fmt.Errorf("nats connection failed: %v", err)
	}

	return NatsClient{
		Conn: conn,
	}, nil
}

func (nc NatsClient) RunKVStore(configsBucketName string) (nats.KeyValue, error) {
	jsc, err := nc.Conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("nats jet stream connection failed: %v", err)
	}

	// Attempt to retrieve the Key-Value (KV) store with the specified bucket name
	kvs, err := jsc.KeyValue(configsBucketName)
	if err != nil {
		// If exists or cannot be retrieved -> create a new one
		kvs, err = jsc.CreateKeyValue(&nats.KeyValueConfig{
			Bucket:      configsBucketName,
			Description: "Service configuration bucket",
			// Set the ttl for key-value pairs in the store
			TTL: 5 * time.Minute,
		})
		if err != nil {
			return nil, fmt.Errorf("key-value: unable to bind key value storage bucket: %v", err)
		}
	}

	return kvs, nil
}
