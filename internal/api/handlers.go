package api

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"strings"

	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/internal/store"
	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/pkg/version"
)

type Handler struct {
	Store *store.Store
}

func (h *Handler) Health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func (h *Handler) Version(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"app":     "compliance-viewer",
		"version": version.Version,
		"commit":  version.Commit,
	})
}

// Checks supports optional ?tag= to narrow by framework tag (e.g. CIS).
func (h *Handler) Checks(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	writeJSON(w, http.StatusOK, h.Store.WithTag(tag))
}

func (h *Handler) CheckByID(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	c, ok := h.Store.ByID(id)
	if !ok {
		http.NotFound(w, r)
		return
	}
	writeJSON(w, http.StatusOK, c)
}

func (h *Handler) Drift(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.Store.Drift())
}

func (h *Handler) Summary(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, h.Store.Summary())
}

type frameworkRow struct {
	Tag   string `json:"tag"`
	Count int    `json:"count"`
}

func (h *Handler) Frameworks(w http.ResponseWriter, _ *http.Request) {
	hist := h.Store.TagHistogram()
	out := make([]frameworkRow, 0, len(hist))
	for t, n := range hist {
		out = append(out, frameworkRow{Tag: t, Count: n})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Tag < out[j].Tag })
	writeJSON(w, http.StatusOK, out)
}

func (h *Handler) ExportChecksCSV(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	rows := h.Store.WithTag(tag)
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", `attachment; filename="checks.csv"`)
	cw := csv.NewWriter(w)
	if err := cw.Write([]string{"id", "title", "category", "expected", "observed", "resource", "framework_tags"}); err != nil {
		log.Printf("export csv header: %v", err)
		http.Error(w, "failed to write CSV", http.StatusInternalServerError)
		return
	}
	for _, c := range rows {
		tags := strings.Join(c.FrameworkTags, "|")
		if err := cw.Write([]string{c.ID, c.Title, c.Category, c.Expected, c.Observed, c.Resource, tags}); err != nil {
			log.Printf("export csv row: %v", err)
			return
		}
	}
	cw.Flush()
	if err := cw.Error(); err != nil {
		log.Printf("export csv flush: %v", err)
	}
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	var buf bytes.Buffer
	enc := json.NewEncoder(&buf)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if _, err := w.Write(buf.Bytes()); err != nil {
		log.Printf("write json response: %v", err)
	}
}
