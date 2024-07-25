package main

import (
	"log"

	"github.com/config-loader-concept/configloader"
	"github.com/config-loader-concept/configloader/config"
	natsClient "github.com/config-loader-concept/pkg/nats"
)

func main() {
	var err error

	cfg := config.Load()

	nats := natsClient.NewClient(
		cfg.BuildNatsUrl(),
		cfg.ConfigBucketName,
	)
	nats, err = nats.Connect()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer nats.Close()

	keyValueStore, err := nats.InitKeyValueStore()
	if err != nil {
		log.Fatal(err)
		return
	}

	updates := make(chan string, 4)
	go configloader.RunFileWatcher(cfg.ConfigsList, updates)

	log.Printf("Service up")
	for {
		select {
		case service := <-updates:
			if err = configloader.FileChangeHandler(
				service,
				cfg.GetConfigsList()[service],
				keyValueStore,
			); err != nil {
				log.Fatal(err)
			}
		}
	}

}
