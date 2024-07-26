package srvone

import (
	natsClient "github.com/config-loader-concept/pkg/nats"
)

type Reload struct {
	nats *natsClient.Client
}

func NewReload(nats *natsClient.Client, opts ...Option) *Reload {
	reload := &Reload{
		nats: nats,
	}

	return reload
}

func (r *Reload) HandleConfigUpdate(service string) (chan []byte, chan error) {
	updates := make(chan []byte)
	errChan := make(chan error)
	go r.nats.SubscribeForUpdates(service, updates, errChan)

	return updates, errChan
}
