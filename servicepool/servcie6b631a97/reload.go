package servcie6b631a97

type Reload struct {
	service Service
	repo    Repository
	nats    NatsClient
}

type Option func(*Reload)

func WithService(s Service) Option {
	return func(reload *Reload) {
		reload.service = s
	}
}

func WithRepository(r Repository) Option {
	return func(reload *Reload) {
		reload.repo = r
	}
}

func NewReload(nats NatsClient, opts ...Option) *Reload {
	reload := &Reload{
		nats: nats,
	}

	for _, opt := range opts {
		opt(reload)
	}

	return reload
}

func (r *Reload) HandleConfigUpdate() {
	notify := make(chan []byte)
	go r.nats.SubscribeToKV(notify)

	for {
		select {
		case cfg := <-notify:
			_ = cfg
			// do something
		}
	}
}
