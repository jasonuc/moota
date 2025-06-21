package main

import (
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

type config struct {
	env  string
	cors struct {
		allowedOrigins   []string
		allowedMethods   []string
		allowedHeaders   []string
		exposedHeaders   []string
		allowCredentials bool
		maxAge           int
	}
	server struct {
		port int
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

	cfg.cors.allowedOrigins = getStringArrayEnv("CORS_ALLOWED_ORIGINS", []string{
		"https://moota.localhost",
		"http://localhost:5173",
	})
	cfg.cors.allowedMethods = getStringArrayEnv("CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	cfg.cors.allowedHeaders = getStringArrayEnv("CORS_ALLOWED_HEADERS", []string{"Accept", "Content-Type"})
	cfg.cors.exposedHeaders = getStringArrayEnv("CORS_EXPOSED_HEADERS", []string{"Link"})
	cfg.cors.allowCredentials = getBoolEnv("CORS_ALLOW_CREDENTIALS", true)
	cfg.cors.maxAge = getIntEnv("CORS_MAX_AGE", 300)

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

func getStringArrayEnv(key string, fallback []string) []string {
	origins, ok := os.LookupEnv(key)
	if !ok || origins == "" {
		return fallback
	}
	return strings.Split(origins, ",")
}

func getBoolEnv(key string, fallback bool) bool {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	boolVal, err := strconv.ParseBool(val)
	if err != nil {
		return fallback
	}

	return boolVal
}
