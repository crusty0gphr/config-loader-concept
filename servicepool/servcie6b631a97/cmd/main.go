package main

import (
	"log"

	natsClient "github.com/config-loader-concept/pkg/nats"
	srv "github.com/config-loader-concept/servicepool/servcie6b631a97"
	"github.com/config-loader-concept/servicepool/servcie6b631a97/config"
)

func main() {
	var err error
	baseConfig := config.LoadBase()

	nats := natsClient.NewClient(
		baseConfig.BuildNatsUrl(),
		baseConfig.App.KeyValueBucket,
	)
	nats, err = nats.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer nats.Close()

	service := srv.NewService(baseConfig)

	reload := srv.NewReload(
		nats,
		srv.WithService(service),
	)
	reload.HandleConfigUpdate()
	log.Printf("Service up")
}
