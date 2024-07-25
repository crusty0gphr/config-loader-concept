package servcie6b631a97

import (
	"log"

	natsClient "github.com/config-loader-concept/pkg/nats"
)

type Reload struct {
	service Service
	nats    *natsClient.Client
}

type Option func(*Reload)

func WithService(s Service) Option {
	return func(reload *Reload) {
		reload.service = s
	}
}

func NewReload(nats *natsClient.Client, opts ...Option) *Reload {
	reload := &Reload{
		nats: nats,
	}

	for _, opt := range opts {
		opt(reload)
	}

	return reload
}

func (r *Reload) HandleConfigUpdate() {
	updates := make(chan []byte)
	errChan := make(chan error)
	go r.nats.SubscribeForUpdates(r.service.Name, updates, errChan)

	for {
		select {
		case config := <-updates:
			_ = config
			log.Print("config update received")
		case err := <-errChan:
			log.Fatalf("key value subsriber error: %v", err)
		}
	}
}
