package main

import (
	"database/sql"
	"log"

	srv "github.com/config-loader-concept/servicepool/servcie6b631a97"
	"github.com/config-loader-concept/servicepool/servcie6b631a97/config"
)

func main() {
	baseConfig := config.LoadBase()
	externalConfig := config.LoadExternal(baseConfig.GetExternalConfigPath())

	service := srv.NewService(baseConfig)
	_ = service

	pgHost := externalConfig.Database.Host
	db, err := sql.Open("postgres", pgHost)
	if err != nil {
		log.Fatalf("unable to connect to postgres: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Fatalf("failed closing db connection: %v", err)
		}
	}()

	repository := srv.NewRepository(db)
	if err = repository.Ping(); err != nil {
		log.Fatal(err)
	}

	nats, err := srv.NewNatsClient(baseConfig.App.Name, "configs", externalConfig.BuildNatsUrl())
	if err != nil {
		log.Fatal(err)
	}
	defer nats.Conn.Close()

	reload := srv.NewReload(
		nats,
		srv.WithService(service),
		srv.WithRepository(repository),
	)
	reload.HandleConfigUpdate()
}
