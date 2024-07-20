package main

import (
	"log"

	"github.com/config-loader-concept/configloader"
	"github.com/config-loader-concept/configloader/config"
)

func main() {
	cfg := config.Load()

	nats, err := configloader.NewNatsClient(cfg.BuildNatsUrl())
	if err != nil {
		log.Fatal(err)
		return
	}

	defer func() {
		log.Print("nat connection closed")
		nats.Conn.Close()
	}()

	kvs, err := nats.RunKVStore(cfg.ConfigBucketName)
	if err != nil {
		log.Fatal(err)
		return
	}

	updates := make(chan string, 4)
	go configloader.RunFileWatcher(cfg.Configs, updates)

	for {
		select {
		case service := <-updates:
			if err = configloader.FileChangeHandler(
				service,
				cfg.GetConfigsList()[service],
				kvs,
			); err != nil {
				log.Fatal(err)
			}
		}
	}
}
