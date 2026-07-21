package codereview

// Issue represents a single finding from code review analysis.
type Issue struct {
	File     string `json:"file"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Severity string `json:"severity"` // "critical", "high", "medium", "low", "info"
	Category string `json:"category"`
	Rule     string `json:"rule"`
	Message  string `json:"message"`
	Recommendation string `json:"recommendation,omitempty"`
	Snippet  string `json:"snippet,omitempty"`
}

// Report is the full code review report.
type Report struct {
	TotalFiles int    `json:"total_files"`
	TotalLines int    `json:"total_lines"`
	Issues     []Issue `json:"issues"`
	SeveritySummary map[string]int `json:"severity_summary"`
}