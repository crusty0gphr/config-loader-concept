package main

import (
	"github.com/config-loader-concept/configloader"
	"github.com/config-loader-concept/configloader/config"
)

func main() {
	cfg := config.Load()

	updates := make(chan string, 4)
	go configloader.RunFileWatcher(cfg.Configs, updates)
}
