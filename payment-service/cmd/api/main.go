package main

import (
	"bmstu-dips-lab2/config"
	"bmstu-dips-lab2/payment-service/internal/server"
	"bmstu-dips-lab2/pkg/postgres"
	"log"
)

func main() {
	log.Println("Starting api server")

	cfgFile, err := config.LoadConfig("./config/config")
	if err != nil {
		log.Fatalf("LoadConfig: %v", err)
	}

	cfg, err := config.ParseConfig(cfgFile)
	if err != nil {
		log.Fatalf("ParseConfig: %v", err)
	}

	psqlDB, err := postgres.New(cfg)
	if err != nil {
		log.Fatalf("Postgresql init: %s", err)
	} else {
		log.Println("Connected to PostreSQL")
	}
	defer psqlDB.Close()

	s := server.New(cfg, psqlDB)
	if err = s.Run(); err != nil {
		log.Fatal(err)
	}
}
