package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jasonuc/moota/internal/utils"
)

func (app *application) serveStaticFiles(r chi.Router) {
	staticPath := app.cfg.static.path

	if _, err := os.Stat(staticPath); os.IsNotExist(err) {
		app.logger.Panicf("static files directory not found: %s", staticPath)
		return
	}

	r.Handle("/assets/*", http.StripPrefix("/assets/", http.FileServer(http.Dir(filepath.Join(staticPath, "assets")))))

	r.Handle("/favicon.ico", http.FileServer(http.Dir(staticPath)))
	r.Handle("/moota.png", http.FileServer(http.Dir(staticPath)))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			utils.NotFoundResponse(w)
			return
		}

		indexPath := filepath.Join(staticPath, "index.html")
		http.ServeFile(w, r, indexPath)
	})
}
