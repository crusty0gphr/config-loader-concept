package nats

import (
	"fmt"
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

const valueTTL = 5 * time.Minute

type Client struct {
	kvBucket     string
	url          string
	natsConn     *nats.Conn
	jetStreamCtx nats.JetStreamContext
}

func NewClient(natsUrl, kvBucket string) *Client {
	return &Client{
		kvBucket: kvBucket,
		url:      natsUrl,
	}
}

func (c *Client) Connect() (*Client, error) {
	conn, err := nats.Connect(c.url)
	if err != nil {
		return nil, fmt.Errorf("nats connection failed: %v", err)
	}

	jetStreamCtx, err := conn.JetStream()
	if err != nil {
		return nil, fmt.Errorf("nats jet stream connection failed: %v", err)
	}

	c.natsConn = conn
	c.jetStreamCtx = jetStreamCtx
	return c, nil
}

func (c *Client) InitKeyValueStore() (nats.KeyValue, error) {
	// try to retrieve the KV store with the specified bucket name
	kv, err := c.jetStreamCtx.KeyValue(c.kvBucket)
	if err != nil {
		// create a new one if unable to retrieved
		kv, err = c.jetStreamCtx.CreateKeyValue(&nats.KeyValueConfig{
			Bucket: c.kvBucket,
			TTL:    valueTTL, // key-value pairs TTL
		})
		if err != nil {
			return nil, fmt.Errorf("key-value: unable to bind key value storage bucket: %v", err)

		}
	}
	return kv, nil
}

func (c *Client) SubscribeForUpdates(key string, update chan []byte, errChan chan error) {
	kv, err := c.jetStreamCtx.KeyValue(c.kvBucket)
	if err != nil {
		errChan <- fmt.Errorf("unable to subscribe to KV storage: %w", err)
		return
	}

	watch, err := kv.Watch(key)
	if err != nil {
		errChan <- fmt.Errorf("unable to watch kv update channel: %w", err)
		return
	}
	defer watch.Stop()

	for {
		entry := <-watch.Updates()
		if entry != nil {
			update <- entry.Value()
			log.Printf("update was received and notified")
		}
	}
}

func (c *Client) Close() {
	c.natsConn.Close()
}
