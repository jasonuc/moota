package main

import (
	"net/http"
	"os"
	"strconv"
	"time"
)

type config struct {
	env    string
	server struct {
		port int
	}
	static struct {
		path string
	}
	db struct {
		dsn             string
		maxOpenConns    int
		maxIdleConns    int
		connMaxIdleTime time.Duration
	}
	auth struct {
		accessTokenSecret  string
		accessTokenTTL     time.Duration
		refreshTokenTTL    time.Duration
		issuer             string
		cookieDomain       string
		cookieSameSiteMode int
	}
}

func parseConfig() config {
	var cfg config

	cfg.env = getStringEnv("ENV", "development")

	cfg.server.port = getIntEnv("PORT", 8080)

	cfg.static.path = getStringEnv("STATIC_FILES_PATH", "./web/dist")

	cfg.db.dsn = getStringEnv("DB_DSN", "postgresql://postgres:postgres@localhost:5432/moota?sslmode=disable&connect_timeout=10")
	cfg.db.maxOpenConns = getIntEnv("DB_MAX_OPEN_CONNS", 15)
	cfg.db.maxIdleConns = getIntEnv("DB_MAX_IDLE_CONNS", 10)
	cfg.db.connMaxIdleTime = getTimeDurationEnv("DB_CONN_MAX_IDLE_TIME", 1*time.Hour)

	cfg.auth.accessTokenSecret = getStringEnv("AUTH_ACCESS_TOKEN_SECRET", "moo_goes_the_cow")
	cfg.auth.accessTokenTTL = getTimeDurationEnv("AUTH_ACCESS_TOKEN_TTL", 24*time.Hour)
	cfg.auth.refreshTokenTTL = getTimeDurationEnv("AUTH_REFRESH_TOKEN_TTL", 7*24*time.Hour)
	cfg.auth.issuer = getStringEnv("AUTH_ISSUER", "moota")
	cfg.auth.cookieDomain = getStringEnv("AUTH_COOKIE_DOMAIN", "")
	cfg.auth.cookieSameSiteMode = getIntEnv("AUTH_COOKIE_SAME_SITE_MODE", int(http.SameSiteStrictMode))

	return cfg
}

func getStringEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	return val
}

func getIntEnv(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	intVal, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return intVal
}

func getTimeDurationEnv(key string, fallback time.Duration) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		return fallback
	}

	return duration
}
