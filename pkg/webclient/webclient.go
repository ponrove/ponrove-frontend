package webclient

import (
	"embed"
	"io/fs"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
)

// Embed the build directory from the frontend.
//
//go:embed build/*
var BuildFs embed.FS

func Register(router chi.Router) error {
	build, err := fs.Sub(BuildFs, "build")
	if err != nil {
		return err
	}

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// If the request is not for the root or a specific app path, serve index.html
		// to allow the frontend to handle routing.
		if r.RequestURI != "/" && !strings.HasPrefix(r.RequestURI, "/_app") {
			if _, err := build.Open(strings.TrimPrefix(r.RequestURI, "/")); err != nil {
				log.Info().Err(err).Msgf("Serving index.html for %s", r.RequestURI)
				http.ServeFileFS(w, r, build, "/index.html")
				return
			}
		}

		// Serve the requested file from the build directory
		http.FileServer(http.FS(build)).ServeHTTP(w, r)
	})

	return nil
}
