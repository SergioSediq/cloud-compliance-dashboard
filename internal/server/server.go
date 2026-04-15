package server

import (
	"io/fs"
	"net/http"

	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/internal/api"
)

// New wires HTTP API routes and static assets. API paths are registered before the file server.
func New(h *api.Handler, static fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /health", h.Health)
	mux.HandleFunc("GET /api/v1/version", h.Version)
	mux.HandleFunc("GET /api/v1/checks", h.Checks)
	mux.HandleFunc("GET /api/v1/checks/{id}", h.CheckByID)
	mux.HandleFunc("GET /api/v1/drift", h.Drift)
	mux.HandleFunc("GET /api/v1/summary", h.Summary)
	mux.HandleFunc("GET /api/v1/frameworks", h.Frameworks)
	mux.HandleFunc("GET /api/v1/export/checks.csv", h.ExportChecksCSV)
	mux.Handle("/", http.FileServer(http.FS(static)))
	return api.RequestLog(mux)
}
