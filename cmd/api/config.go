package main

import (
	"flag"
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

func parseConfig() config {
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

	return cfg
}

func parseDurationForConfig(input string, ptr *time.Duration) {
	duration, err := time.ParseDuration(input)
	if err != nil {
		return
	}

	(*ptr) = duration
}
