package main

import (
	"database/sql"
	"log"

	dbClient "github.com/config-loader-concept/pkg/db"
	natsClient "github.com/config-loader-concept/pkg/nats"
	srv "github.com/config-loader-concept/services/srvone"
	"github.com/config-loader-concept/services/srvone/config"
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

	externalConfig := config.LoadExternal(baseConfig.GetExternalConfigPath())
	db := dbClient.Connect(externalConfig.Database.Host)
	defer func(db *sql.DB) {
		err = db.Close()
		if err != nil {
			log.Fatalf("Failed closing DB connection")
		}
	}(db)

	repo := srv.NewRepo(db)
	repo.Ping()

	reload := srv.NewReload(nats)

	service := srv.NewService(
		baseConfig,
		externalConfig,
		srv.WithRepo(repo),
		srv.WithReload(reload),
	)
	service.Start()
	service.ReloadHandler()
	log.Printf("Service up")
}
