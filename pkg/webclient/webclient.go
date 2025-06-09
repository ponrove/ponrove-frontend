package webclient

import (
	"net/http"
	"os" // Added os import
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/ponrove/configura"
	"github.com/rs/zerolog/log"
)

const (
	WEBCLIENT_APP_BUILD_DIR configura.Variable[string] = "WEBCLIENT_APP_BUILD_DIR"
)

// Register sets up the web client routes for serving static files from the build directory. It's a static file server
// and does not require to be a part of the OpenAPI specification, as it serves the frontend application.
func Register(cfg configura.Config, router chi.Router) error {
	// Use os.DirFS to serve files from the "build" directory on the filesystem.
	// "build" should be a directory relative to where the application is run.
	build := os.DirFS(cfg.String(WEBCLIENT_APP_BUILD_DIR))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		// If the request is not for the root or a specific app path, try to serve index.html
		// to allow the frontend to handle routing. This is a common pattern for SPAs.
		// We check if the requested file exists; if not, serve index.html.
		requestedPath := strings.TrimPrefix(r.RequestURI, "/")
		if requestedPath == "" { // Handle root path explicitly
			requestedPath = "index.html"
		}

		// Attempt to open the requested file.
		f, err := build.Open(requestedPath)
		if err != nil {
			// If the file doesn't exist (os.ErrNotExist), and it's not an _app asset,
			// serve index.html.
			// This allows client-side routing for paths like /users/123.
			if os.IsNotExist(err) && !strings.HasPrefix(r.URL.Path, "/_app") {
				log.Info().Msgf("Serving index.html for %s as %s not found", r.RequestURI, requestedPath)
				http.ServeFileFS(w, r, build, "index.html")
				return
			}
			// For other errors, or if it's an _app asset that's not found,
			// let the FileServer handle it (which will likely result in a 404).
		}
		if f != nil {
			f.Close() // Close the file if we opened it for the check
		}

		// Serve the requested file from the build directory or index.html via the FileServer.
		// The FileServer will correctly serve index.html for directory requests.
		http.FileServer(http.FS(build)).ServeHTTP(w, r)
	})

	return nil
}
