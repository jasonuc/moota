package main

import (
	"log"
	"os"

	"github.com/jasonuc/moota/internal/handlers"
	"github.com/jasonuc/moota/internal/middlewares"
	"github.com/jasonuc/moota/internal/services"
	"github.com/jasonuc/moota/internal/store"
	"github.com/joho/godotenv"
)

var version = "0.0.1"

type application struct {
	cfg    config
	logger *log.Logger

	store *store.Store

	plantService services.PlantService
	soilService  services.SoilService
	seedService  services.SeedService
	authService  services.AuthService
	userService  services.UserService

	authMiddleware middlewares.AuthMiddleware

	authHandler  *handlers.AuthHandler
	seedHandler  *handlers.SeedHandler
	plantHandler *handlers.PlantHandler
	userHandler  *handlers.UserHandler
}

func main() {
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	cfg := parseConfig()

	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	db, err := openDB(cfg)
	if err != nil {
		logger.Panicf("error: %v\n", err)
	}
	//nolint:errcheck
	defer db.Close()

	store := store.NewStore(db)

	plantService := services.NewPlantService(store)
	soilService := services.NewSoilSerivce(store)
	seedService := services.NewSeedService(store, soilService, plantService)
	authService := services.NewAuthService(store, []byte(cfg.auth.accessTokenSecret), cfg.auth.refreshTokenTTL, cfg.auth.accessTokenTTL, cfg.auth.issuer)
	userService := services.NewUserService(store)

	authMiddlware := middlewares.NewAuthMiddleware(authService)

	authHandler := handlers.NewAuthHandler(authService)
	seedHandler := handlers.NewSeedHandler(seedService)
	plantHandler := handlers.NewPlantService(plantService)
	userHandler := handlers.NewUserHandler(userService)

	app := application{
		cfg:    cfg,
		logger: logger,
		store:  store,

		plantService: plantService,
		soilService:  soilService,
		seedService:  seedService,
		authService:  authService,
		userService:  userService,

		authMiddleware: authMiddlware,

		authHandler:  authHandler,
		seedHandler:  seedHandler,
		plantHandler: plantHandler,
		userHandler:  userHandler,
	}

	if err := app.serve(); err != nil {
		log.Fatalf("Server shutting down forcefully: %v\n", err)
	} else {
		log.Print("server shutting down gracefully\n")
	}
}
