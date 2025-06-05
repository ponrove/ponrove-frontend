package app

import (
	"embed"
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Embed the build directory from the frontend.
//
//go:embed build/*
var BuildFs embed.FS

func Register(router chi.Router) {
	build, err := fs.Sub(BuildFs, "build")
	if err != nil {
		log.Fatal(err)
	}
	router.Handle("/*", http.FileServerFS(build))
}
