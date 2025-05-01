package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", app.authHandler.HandleRegisterRequest)
			r.Post("/login", app.authHandler.HandleLoginRequest)
			r.Post("/refresh", app.authHandler.HandleTokenRefresh)
		})

		r.Group(func(r chi.Router) {
			r.Use(app.authMiddleware.Authorise)

			r.Route("/plants", func(r chi.Router) {
				r.Get("/u/{userID}", app.plantHandler.HandleGetAllUserPlants)
				r.Get("/u/{userID}/graveyard", app.plantHandler.HandleGetAllUserDeceasedPlants)

				r.Get("/{plantID}", app.plantHandler.HandleGetPlant)

				r.Post("/{plantID}/action", app.plantHandler.HandleActionOnPlant)
				r.Post("/{plantID}/activate", app.plantHandler.HandleActivatePlant)
				r.Post("/{plantID}/kill", app.plantHandler.HandleKillPlant)
			})

			r.Route("/seeds", func(r chi.Router) {
				r.Get("/u/{userID}", app.seedHandler.HandleGetUserSeeds)

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
