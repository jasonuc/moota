package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://moota.localhost",
			"http://localhost:5173",
			"http://127.0.0.1:5173",
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.authHandler.HandleRegisterRequest)
			r.Post("/login", app.authHandler.HandleLoginRequest)
			r.Post("/refresh", app.authHandler.HandleTokenRefresh)
			r.Post("/logout", app.authHandler.HandleLogout)

			r.Route("/u/{userID}", func(r chi.Router) {
				r.Use(app.authMiddleware.Authorise)
				r.Use(app.authMiddleware.ValidateUserAccess)

				r.Patch("/change-email", app.authHandler.HandleChangeEmail)
				r.Patch("/change-password", app.authHandler.HandleChangePassword)
				r.Patch("/change-username", app.authHandler.HandleChangeUsername)
			})
		})

		r.Route("/whoami", func(r chi.Router) {
			r.Use(app.authMiddleware.Authorise)
			r.Get("/", app.userHandler.HandleWhoAmI)
		})

		r.Group(func(r chi.Router) {
			r.Use(app.authMiddleware.Authorise)

			r.Route("/users", func(r chi.Router) {
				r.Get("/{username}/profile", app.userHandler.HandleGetUserProfile)

				r.Route("/u/{userID}", func(r chi.Router) {
					r.Use(app.authMiddleware.ValidateUserAccess)

					r.Get("/", app.userHandler.HandleGetUser)
				})
			})

			r.Route("/plants", func(r chi.Router) {
				r.Route("/u/{userID}", func(r chi.Router) {
					r.Use(app.authMiddleware.ValidateUserAccess)

					r.Get("/", app.plantHandler.HandleGetAllUserPlants)
					r.Get("/graveyard", app.plantHandler.HandleGetAllUserDeceasedPlants)
				})

				r.Route("/{plantID}", func(r chi.Router) {
					r.Get("/", app.plantHandler.HandleGetPlant)
					r.Patch("/", app.plantHandler.HandleChangePlantNickname)
					r.Post("/action", app.plantHandler.HandleActionOnPlant)
					r.Post("/kill", app.plantHandler.HandleKillPlant)
				})
			})

			r.Route("/seeds", func(r chi.Router) {
				r.Route("/u/{userID}", func(r chi.Router) {
					r.Use(app.authMiddleware.ValidateUserAccess)

					r.Get("/", app.seedHandler.HandleGetUserSeeds)
					r.Get("/request", app.seedHandler.HandleCheckWhenUserCanRequestSeed)
					r.Post("/request", app.seedHandler.HandleRequestForNewSeeds)
				})

				r.Post("/{seedID}", app.seedHandler.HandlePlantSeed)
			})
		})
	})

	return r
}

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	m := map[string]any{
		"system_status": "available",
		"system_info": map[string]any{
			"environment": app.cfg.env,
			"version":     version,
		},
	}

	js, err := json.Marshal(m)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	if _, err := w.Write(js); err != nil {
		app.logger.Println(err)
	}
}
