package store_test

import (
	"path/filepath"
	"testing"

	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/internal/store"
)

func TestLoadAndDrift(t *testing.T) {
	p := filepath.Join("..", "..", "data", "checks.json")
	s, err := store.Load(p)
	if err != nil {
		t.Fatal(err)
	}
	all := s.All()
	if len(all) < 1 {
		t.Fatalf("expected checks")
	}
	d := s.Drift()
	if len(d) < 1 {
		t.Fatalf("expected at least one drift row in fixture data")
	}
	su := s.Summary()
	if su.Total != len(all) {
		t.Errorf("summary total %d want %d", su.Total, len(all))
	}
	if _, ok := s.ByID("cis-1.1"); !ok {
		t.Fatal("expected ByID")
	}
	cis := s.WithTag("CIS")
	if len(cis) < 1 {
		t.Fatal("expected CIS tag filter")
	}
	if len(s.TagHistogram()) < 1 {
		t.Fatal("expected tag histogram")
	}
}
