package servcie6b631a97

import (
	"fmt"

	"github.com/config-loader-concept/servicepool/servcie6b631a97/config"
)

type Service struct {
	Name    string
	Version string
	Host    string
	Port    int
}

func (s Service) String() string {
	return fmt.Sprintf(
		"Service %s:%s loaded new config | host: %s - port %d",
		s.Name, s.Version, s.Host, s.Port,
	)
}

func NewService(base config.Base) Service {
	return Service{
		Name:    base.App.Name,
		Version: base.App.Name,
		Host:    base.Server.Host,
		Port:    base.Server.Port,
	}
}
