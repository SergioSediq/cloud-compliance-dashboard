package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"

	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/internal/store"
)

func testStore(t *testing.T) *store.Store {
	t.Helper()
	p := filepath.Join("..", "..", "data", "checks.json")
	st, err := store.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	return st
}

func TestHealth(t *testing.T) {
	h := &Handler{Store: testStore(t)}
	rr := httptest.NewRecorder()
	h.Health(rr, httptest.NewRequest(http.MethodGet, "/health", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Fatalf("got %#v", body)
	}
}

func TestChecks_JSON(t *testing.T) {
	h := &Handler{Store: testStore(t)}
	rr := httptest.NewRecorder()
	h.Checks(rr, httptest.NewRequest(http.MethodGet, "/api/v1/checks", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	var raw []json.RawMessage
	if err := json.Unmarshal(rr.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if len(raw) < 1 {
		t.Fatal("expected non-empty array")
	}
}

func TestSummary_JSON(t *testing.T) {
	h := &Handler{Store: testStore(t)}
	rr := httptest.NewRecorder()
	h.Summary(rr, httptest.NewRequest(http.MethodGet, "/api/v1/summary", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	var body struct {
		Total int `json:"total"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body.Total < 1 {
		t.Fatalf("total %d", body.Total)
	}
}

func TestVersion_JSON(t *testing.T) {
	h := &Handler{Store: testStore(t)}
	rr := httptest.NewRecorder()
	h.Version(rr, httptest.NewRequest(http.MethodGet, "/api/v1/version", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if body["app"] != "compliance-viewer" || body["version"] == "" {
		t.Fatalf("got %#v", body)
	}
}

func TestCheckByID_OK(t *testing.T) {
	h := &Handler{Store: testStore(t)}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/checks/cis-1.1", nil)
	req.SetPathValue("id", "cis-1.1")
	rr := httptest.NewRecorder()
	h.CheckByID(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	var c map[string]any
	if err := json.Unmarshal(rr.Body.Bytes(), &c); err != nil {
		t.Fatal(err)
	}
	if c["id"] != "cis-1.1" {
		t.Fatalf("got %#v", c)
	}
}

func TestChecks_TagFilter(t *testing.T) {
	h := &Handler{Store: testStore(t)}
	rr := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/checks?tag=CIS", nil)
	h.Checks(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	var raw []json.RawMessage
	if err := json.Unmarshal(rr.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if len(raw) < 1 {
		t.Fatal("expected CIS-tagged rows")
	}
}

func TestDrift_JSON(t *testing.T) {
	h := &Handler{Store: testStore(t)}
	rr := httptest.NewRecorder()
	h.Drift(rr, httptest.NewRequest(http.MethodGet, "/api/v1/drift", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	var raw []json.RawMessage
	if err := json.Unmarshal(rr.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if len(raw) < 1 {
		t.Fatal("expected drift rows")
	}
}

func TestFrameworks_JSON(t *testing.T) {
	h := &Handler{Store: testStore(t)}
	rr := httptest.NewRecorder()
	h.Frameworks(rr, httptest.NewRequest(http.MethodGet, "/api/v1/frameworks", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	var raw []json.RawMessage
	if err := json.Unmarshal(rr.Body.Bytes(), &raw); err != nil {
		t.Fatal(err)
	}
	if len(raw) < 1 {
		t.Fatal("expected framework rows")
	}
}

func TestCheckByID_NotFound(t *testing.T) {
	h := &Handler{Store: testStore(t)}
	req := httptest.NewRequest(http.MethodGet, "/api/v1/checks/does-not-exist", nil)
	req.SetPathValue("id", "does-not-exist")
	rr := httptest.NewRecorder()
	h.CheckByID(rr, req)
	if rr.Code != http.StatusNotFound {
		t.Fatalf("status %d", rr.Code)
	}
}

func TestExportChecksCSV(t *testing.T) {
	h := &Handler{Store: testStore(t)}
	rr := httptest.NewRecorder()
	h.ExportChecksCSV(rr, httptest.NewRequest(http.MethodGet, "/api/v1/export/checks.csv", nil))
	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	ct := rr.Header().Get("Content-Type")
	if ct != "text/csv; charset=utf-8" {
		t.Fatalf("content-type %q", ct)
	}
	body := rr.Body.String()
	if len(body) < 20 || body[:2] != "id" {
		n := 40
		if len(body) < n {
			n = len(body)
		}
		t.Fatalf("unexpected csv prefix: %q", body[:n])
	}
}
