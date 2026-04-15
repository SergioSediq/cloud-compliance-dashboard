package models

// Check represents one benchmark-style row (mock or collected).
type Check struct {
	ID            string   `json:"id"`
	Title         string   `json:"title"`
	Category      string   `json:"category"`
	FrameworkTags []string `json:"framework_tags"`
	Expected      string   `json:"expected"`
	Observed      string   `json:"observed"`
	Resource      string   `json:"resource"`
}

// DriftItem is a check where observed != expected.
type DriftItem struct {
	Check
	DriftType string `json:"drift_type"`
}
