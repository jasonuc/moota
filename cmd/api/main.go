package main

import (
	"log"
	"os"
)

var version = "0.0.1"

type application struct {
	cfg    config
	logger *log.Logger
}

func main() {
	cfg := parseConfig()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	_, err := openDB(cfg)
	if err != nil {
		logger.Panicf("error: %v\n", err)
	}

	app := application{
		cfg:    cfg,
		logger: logger,
	}

	if err := app.serve(); err != nil {
		log.Fatalf("Server shutting down forcefully: %v\n", err)
	} else {
		log.Print("server shutting down gracefully\n")
	}
}
