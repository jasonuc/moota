package main

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (app *application) routes() http.Handler {
	r := chi.NewRouter()

	r.Route("/api", func(r chi.Router) {
		r.Get("/health", app.healthCheckHandler)

		r.Group(func(r chi.Router) {
			r.Use(app.authMiddleware.Authorise)
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
