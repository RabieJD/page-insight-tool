package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rabie/page-insight-tool/app/handlers"
)

// New returns a new router
func New() http.Handler {
	r := mux.NewRouter()

	// Serve static files
	fs := http.FileServer(http.Dir("app/static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Route handlers
	r.HandleFunc("/", handlers.IndexHandler).Methods("GET")
	r.HandleFunc("/analyze", handlers.AnalyzeHandler).Methods("POST")

	return r
}
