package main

import (
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/jasonuc/moota/internal/utils"
	"github.com/jasonuc/moota/web"
)

func (app *application) serveStaticFiles(r chi.Router) {
	staticPath := app.cfg.static.path

	staticFs, err := fs.Sub(web.Files, staticPath)
	if err != nil {
		app.logger.Panicf("%s does not exist on the embedded file system: %v", staticPath, err)
	}

	r.Handle("/assets/*", http.FileServerFS(staticFs))

	r.Handle("/favicon.ico", http.FileServerFS(staticFs))
	r.Handle("/moota.png", http.FileServerFS(staticFs))

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/api/") {
			utils.NotFoundResponse(w)
			return
		}

		http.ServeFileFS(w, r, staticFs, "index.html")
	})
}
