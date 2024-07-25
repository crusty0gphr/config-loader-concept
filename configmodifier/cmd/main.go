package main

import (
	"log"
	"math/rand/v2"
	"time"

	"github.com/config-loader-concept/configmodifier"
)

var configs = []string{
	"configs/cfg.e079e8d7-9153-42a0-8b4c-d4e0022c2e6b.yaml",
	"configs/cfg.6b631a97-b02f-4d77-8366-4a75f99b3d26.yaml",
	"configs/cfg.9fd2041b-88bc-4ac1-8157-0188d37d0f75.yaml",
}

func main() {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	log.Printf("Service up")
	for {
		i := rand.IntN(len(configs))

		select {
		case <-ticker.C:
			configmodifier.ModifyConfigFile(configs[i])
		}
	}
}
