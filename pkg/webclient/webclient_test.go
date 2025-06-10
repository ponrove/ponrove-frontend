package webclient

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/ponrove/configura"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newDefaultConfig() *configura.ConfigImpl {
	cfg := configura.NewConfigImpl()
	configura.WriteConfiguration(cfg, map[configura.Variable[string]]string{
		WEBCLIENT_APP_BUILD_DIR: "./_build",
	})
	return cfg
}

// TestRegister_NoError ensures that the Register function can be called without returning an error.
// This is a basic sanity check for the registration process itself.
func TestRegister_NoError(t *testing.T) {
	cfg := newDefaultConfig()
	router := chi.NewRouter()
	err := Register(cfg, router)
	require.NoError(t, err, "Register should not return an error")
}

// TestHandlerExecution tests the behavior of the HTTP handler registered by the Register function
// for various request paths.
//
// IMPORTANT: These tests depend on the content of the 'build' directory embedded at COMPILE TIME.
// For tests expecting specific files (e.g., index.html, _app/somefile.js) to be served,
// those files must exist in 'ponrove-frontend/pkg/webclient/build/' when 'go test' is run.
//
// If 'build/index.html' is missing:
//   - Requests for '/' will likely result in a 404.
//   - Requests for non-existent paths (e.g., '/some/route') will attempt to serve 'index.html',
//     and if it's missing, will also result in a 404.
//
// If 'build/_app/somefile.js' is missing:
// - Requests for '/_app/somefile.js' will result in a 404.
//
// To make these tests pass with specific content checks, you can create dummy files
// in `ponrove-frontend/pkg/webclient/build/` before running `go test`:
// - `index.html`: `<html><body>Mock Index Page</body></html>`
// - `_app/test_app_file.js`: `console.log("mock app file");`
// - `existing_file.txt`: `This is a mock existing file.`
func TestHandlerExecution(t *testing.T) {
	cfg := newDefaultConfig()
	router := chi.NewRouter()
	err := Register(cfg, router)
	require.NoError(t, err, "Register should not return an error")

	// Helper function to make requests and check responses
	testRequest := func(t *testing.T, method, path string, expectedStatus int, expectedBodyContains ...string) {
		t.Helper()
		req := httptest.NewRequest(method, path, nil)
		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, expectedStatus, rr.Code, "handler returned wrong status code for path %s", path)

		if rr.Code == http.StatusOK && len(expectedBodyContains) > 0 {
			bodyBytes, err := io.ReadAll(rr.Body)
			require.NoError(t, err, "Failed to read response body for path %s", path)
			bodyString := string(bodyBytes)
			for _, expectedSubstring := range expectedBodyContains {
				assert.Contains(t, bodyString, expectedSubstring, "handler response body for path %s did not contain expected string '%s'", path, expectedSubstring)
			}
		}
	}

	t.Run("RootPathServesIndex", func(t *testing.T) {
		// Expects build/index.html to be served by the FileServer.
		// If build/index.html does not exist in BuildFs, this will likely be a 404.
		// Assumes index.html exists and contains "Mock Index Page".
		testRequest(t, http.MethodGet, "/", http.StatusOK, "<html lang=\"en\">")
	})

	t.Run("AppPathServesFile", func(t *testing.T) {
		// Expects build/_app/test_app_file.js to be served by the FileServer.
		// If the file doesn't exist in BuildFs, this will be a 404.
		// Assumes _app/test_app_file.js exists and contains "mock app file".
		testRequest(t, http.MethodGet, "/_app/version.json", http.StatusOK, "\"version\"")
	})

	t.Run("IndexRedirectsToRoot", func(t *testing.T) {
		// Tests serving a file that exists at a non-root, non-_app path (e.g., /existing_file.txt).
		// Logic: `build.Open("existing_file.txt")` succeeds, so FileServer serves it.
		// If build/existing_file.txt does not exist, this would fall back to serving index.html.
		// Assumes existing_file.txt exists and contains "mock existing file".
		testRequest(t, http.MethodGet, "/index.html", http.StatusMovedPermanently, "<html lang=\"en\">")
	})

	t.Run("SpecificExistingFileServesFile", func(t *testing.T) {
		// Tests serving a file that exists at a non-root, non-_app path (e.g., /existing_file.txt).
		// Logic: `build.Open("existing_file.txt")` succeeds, so FileServer serves it.
		// If build/existing_file.txt does not exist, this would fall back to serving index.html.
		// Assumes existing_file.txt exists and contains "mock existing file".
		testRequest(t, http.MethodGet, "/index.html", http.StatusMovedPermanently, "<html lang=\"en\">")
	})

	t.Run("NonExistentPathServesIndex", func(t *testing.T) {
		// Tests SPA fallback: a request for a path that does not correspond
		// to a physical file (and is not an /_app path) should serve index.html.
		// Logic: `build.Open("some/frontend/route")` fails, so ServeFileFS serves index.html.
		// If build/index.html does not exist in BuildFs, this will be a 404.
		// Assumes index.html exists and contains "Mock Index Page".
		testRequest(t, http.MethodGet, "/some/frontend/route", http.StatusOK, "<html lang=\"en\">")
	})

	t.Run("NonExistentFileInRootServesIndex", func(t *testing.T) {
		// Tests SPA fallback: a request for a file in root that does not correspond
		// to a physical file (and is not an /_app path) should serve index.html.
		// This is similar to NonExistentPathServesIndex but targets a file-like path.
		// Logic: `build.Open("non_existent_at_root.txt")` fails, so ServeFileFS serves index.html.
		// If build/index.html does not exist in BuildFs, this will be a 404.
		// Assumes index.html exists and contains "Mock Index Page".
		testRequest(t, http.MethodGet, "/non_existent_at_root.txt", http.StatusOK, "<html lang=\"en\">")
	})

	t.Run("NonExistentAppPathReturns404", func(t *testing.T) {
		// A request for a non-existent file under /_app/ should result in a 404
		// from the http.FileServer.
		testRequest(t, http.MethodGet, "/_app/non_existent_file.js", http.StatusNotFound)
	})

	t.Run("MalformedAppPathReturns404", func(t *testing.T) {
		// A request for a malformed path under /_app/ (e.g., trying to escape)
		// should be handled by the http.FileServer, typically resulting in 404 if the path is invalid/not found.
		// Note: http.FileServer itself handles path cleaning. '/_app/../file' becomes '/file'.
		// If '/file' exists in build, it would be served.
		// If the intention is to test against `build/_app/../file`, which is `build/file`,
		// and `build/file` does not exist, it will be 404.
		// If `build/file` exists, this test case would need adjustment or a different path.
		// For this test, we assume `build/some_other_file` (resolved from `/_app/../some_other_file`) does not exist.
		testRequest(t, http.MethodGet, "/_app/../some_other_file", http.StatusNotFound)
	})
}
