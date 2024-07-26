package srvone

import (
	"fmt"
	"log"

	"github.com/config-loader-concept/pkg/db"
	"github.com/config-loader-concept/services/srvtwo/config"
)

type details struct {
	name    string
	version string
	host    string
	port    int
}

func (s details) String() string {
	return fmt.Sprintf(
		"Base config loaded: %s:%s | host: %s - port %d",
		s.name, s.version, s.host, s.port,
	)
}

type features struct {
	logLevel       string
	logFormat      string
	enableFeatureX bool
	enableFeatureY bool
}

func (s features) String() string {
	return fmt.Sprintf(
		"External config loaded: log_level: %s | log_format: %s | feature_x: %v | feature_y: %v",
		s.logLevel, s.logFormat, s.enableFeatureX, s.enableFeatureY,
	)
}

type Service struct {
	details  details
	features features
	repo     *Repo
	reload   *Reload
}

type Option func(*Service)

func WithRepo(r *Repo) Option {
	return func(srv *Service) {
		srv.repo = r
	}
}

func WithReload(r *Reload) Option {
	return func(srv *Service) {
		srv.reload = r
	}
}

func NewService(base config.Base, external config.External, opts ...Option) *Service {
	srv := Service{
		details: details{
			name:    base.App.Name,
			version: base.App.Version,
			host:    base.Server.Host,
			port:    base.Server.Port,
		},
		features: features{
			logLevel:       external.Logging.Level,
			logFormat:      external.Logging.Format,
			enableFeatureX: external.Features.EnableFeatureX,
			enableFeatureY: external.Features.EnableFeatureY,
		},
	}

	for _, opt := range opts {
		opt(&srv)
	}
	return &srv
}

func (s *Service) Start() {
	log.Print(s.details)
	log.Print(s.features)
}

func (s *Service) Ping() {
	s.repo.Ping()
}

func (s *Service) ReloadHandler() {
	updates, errChan := s.reload.HandleConfigUpdate(s.details.name)

	for {
		select {
		case contents := <-updates:
			log.Print("Config update received")

			cfg := config.ParseExternalConfig(contents)
			s.updateConfig(cfg)
		case err := <-errChan:
			log.Fatalf("Key Value subsriber error: %v", err)
		}
	}
}

func (s *Service) updateConfig(external config.External) {
	s.repo.CloseDBConn()

	newDB := db.Connect(external.Database.Host)
	repo := NewRepo(newDB)
	repo.Ping()

	s.repo = repo
	s.features = features{
		logLevel:       external.Logging.Level,
		logFormat:      external.Logging.Format,
		enableFeatureX: external.Features.EnableFeatureX,
		enableFeatureY: external.Features.EnableFeatureY,
	}
	log.Print("Service updated")
	s.Start()
}
