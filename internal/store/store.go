package store

import (
	"encoding/json"
	"os"
	"sync"

	"github.com/SergioSediq/security-portfolio-projects/cloud-compliance-dashboard/internal/models"
)

type Store struct {
	mu     sync.RWMutex
	checks []models.Check
}

func Load(path string) (*Store, error) {
	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var checks []models.Check
	if err := json.Unmarshal(b, &checks); err != nil {
		return nil, err
	}
	return &Store{checks: checks}, nil
}

func (s *Store) All() []models.Check {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]models.Check, len(s.checks))
	copy(out, s.checks)
	return out
}

func (s *Store) Drift() []models.DriftItem {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var d []models.DriftItem
	for _, c := range s.checks {
		if c.Observed != c.Expected {
			d = append(d, models.DriftItem{Check: c, DriftType: "expected_vs_observed"})
		}
	}
	return d
}

type Summary struct {
	Total     int            `json:"total"`
	Pass      int            `json:"pass"`
	Fail      int            `json:"fail"`
	Drift     int            `json:"drift_count"`
	ByTag     map[string]int `json:"by_framework_tag"`
}

// ByID returns one check by stable identifier.
func (s *Store) ByID(id string) (models.Check, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.checks {
		if c.ID == id {
			return c, true
		}
	}
	return models.Check{}, false
}

// WithTag returns checks that include the framework tag (case-sensitive). Empty tag = all rows.
func (s *Store) WithTag(tag string) []models.Check {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if tag == "" {
		out := make([]models.Check, len(s.checks))
		copy(out, s.checks)
		return out
	}
	var out []models.Check
	for _, c := range s.checks {
		for _, t := range c.FrameworkTags {
			if t == tag {
				out = append(out, c)
				break
			}
		}
	}
	return out
}

// TagHistogram counts how many checks reference each framework tag.
func (s *Store) TagHistogram() map[string]int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	m := map[string]int{}
	for _, c := range s.checks {
		for _, t := range c.FrameworkTags {
			m[t]++
		}
	}
	return m
}

func (s *Store) Summary() Summary {
	s.mu.RLock()
	defer s.mu.RUnlock()
	su := Summary{ByTag: map[string]int{}}
	for _, c := range s.checks {
		su.Total++
		if c.Observed == "PASS" {
			su.Pass++
		} else {
			su.Fail++
		}
		if c.Observed != c.Expected {
			su.Drift++
		}
		for _, t := range c.FrameworkTags {
			su.ByTag[t]++
		}
	}
	return su
}
