package main

import (
	"flag"
	"log"
	"os"
	"time"
)

type config struct {
	env    string
	server struct {
		port int
	}
	db struct {
		dsn             string
		maxOpenConns    int
		maxIdleConns    int
		connMaxIdleTime struct {
			input string
			value time.Duration
		}
	}
}

type application struct {
	cfg    config
	logger *log.Logger
}

func main() {
	var cfg config

	flag.StringVar(&cfg.env, "env", "development", "developement|staging|production")

	flag.IntVar(&cfg.server.port, "port", 8080, "port the http server would run on")

	flag.StringVar(&cfg.db.dsn, "db-dsn", "", "database connection string")
	flag.IntVar(&cfg.db.maxOpenConns, "db-maxOpenConns", 15, "max open conns in the database connection pool (in-use/idle)")
	flag.IntVar(&cfg.db.maxIdleConns, "db-maxIdleConns", 10, "max idle conns in the database connection pool")
	flag.StringVar(&cfg.db.connMaxIdleTime.input, "db-connMaxIdleTime", "1h", "max time an idle connection would live in the connection pool")

	flag.Parse()

	// Parse config values into usable values
	// MUST BE DONE AFTER `flag.Parse()` otherwise would work with default values
	parseDurationForConfig(cfg.db.connMaxIdleTime.input, &cfg.db.connMaxIdleTime.value)

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
