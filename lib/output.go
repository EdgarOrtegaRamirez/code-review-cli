package codereview

import (
	"strconv"
	"strings"
)

// FormatJSON renders a report as JSON string.
func FormatJSON(r Report) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString("  \"total_lines\": " + strconv.Itoa(r.TotalLines) + ",\n")
	sb.WriteString("  \"total_files\": " + strconv.Itoa(r.TotalFiles) + ",\n")
	sb.WriteString("  \"issues_count\": " + strconv.Itoa(len(r.Issues)) + ",\n")
	sb.WriteString("  \"severity_summary\": {\n")
	sevs := []string{"critical", "high", "medium", "low", "info"}
	for i, s := range sevs {
		count := r.SeveritySummary[s]
		sb.WriteString("    \"")
		sb.WriteString(s)
		sb.WriteString("\": ")
		sb.WriteString(strconv.Itoa(count))
		if i < len(sevs)-1 {
			sb.WriteString(",\n")
		} else {
			sb.WriteString("\n")
		}
	}
	sb.WriteString("  },\n")
	sb.WriteString("  \"issues\": [\n")
	for i, issue := range r.Issues {
		sb.WriteString("    {\n")
		sb.WriteString("      \"file\": \"")
		sb.WriteString(escapeJSON(issue.File))
		sb.WriteString("\",\n")
		sb.WriteString("      \"line\": " + strconv.Itoa(issue.Line) + ",\n")
		sb.WriteString("      \"severity\": \"")
		sb.WriteString(issue.Severity)
		sb.WriteString("\",\n")
		sb.WriteString("      \"category\": \"")
		sb.WriteString(issue.Category)
		sb.WriteString("\",\n")
		sb.WriteString("      \"rule\": \"")
		sb.WriteString(issue.Rule)
		sb.WriteString("\",\n")
		sb.WriteString("      \"message\": \"")
		sb.WriteString(escapeJSON(issue.Message))
		sb.WriteString("\"")
		if issue.Recommendation != "" {
			sb.WriteString(", \"recommendation\": \"")
			sb.WriteString(escapeJSON(issue.Recommendation))
			sb.WriteString("\"")
		}
		sb.WriteString("\n    }")
		if i < len(r.Issues)-1 {
			sb.WriteString(",")
		}
		sb.WriteString("\n")
	}
	sb.WriteString("  ]\n")
	sb.WriteString("}\n")
	return sb.String()
}

func escapeJSON(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

// FormatMarkdown renders a report as Markdown string.
func FormatMarkdown(r Report) string {
	var sb strings.Builder

	sb.WriteString("# Code Review Report\n\n")
	sb.WriteString("## Summary\n\n")
	sb.WriteString("| Severity | Count |\n")
	sb.WriteString("|----------|-------|\n")

	severities := []string{"critical", "high", "medium", "low", "info"}
	labels := map[string]string{
		"critical": "Critical",
		"high":     "High",
		"medium":   "Medium",
		"low":      "Low",
		"info":     "Info",
	}
	total := 0
	for _, s := range severities {
		count := r.SeveritySummary[s]
		total += count
		sb.WriteString("| ")
		sb.WriteString(labels[s])
		sb.WriteString(" | ")
		sb.WriteString(strings.Repeat("!", count))
		sb.WriteString(" |\n")
	}
	sb.WriteString("| **Total** | **")
	sb.WriteString(strconv.Itoa(total))
	sb.WriteString("** |\n\n")

	sb.WriteString("## Details\n\n")

	for _, issue := range r.Issues {
		sb.WriteString("### ")
		sb.WriteString(issue.Severity)
		sb.WriteString(": ")
		sb.WriteString(issue.Message)
		sb.WriteString("\n\n")
		sb.WriteString("- **File**: `")
		sb.WriteString(issue.File)
		sb.WriteString("` - Line: ")
		sb.WriteString(strconv.Itoa(issue.Line))
		sb.WriteString("\n")
		sb.WriteString("- **Rule**: `" + issue.Rule + "`\n")
		sb.WriteString("- **Snippet**: `" + strings.ReplaceAll(issue.Snippet, "`", "'") + "`\n")
		if issue.Recommendation != "" {
			sb.WriteString("- **Fix**: " + issue.Recommendation + "\n")
		}
		sb.WriteString("\n")
	}

	return sb.String()
}

// FormatTable renders a report as a text table.
func FormatTable(r Report) string {
	var sb strings.Builder
	sb.WriteString("File\tLine\tSeverity\tCategory\tRule\tMessage\n")
	sb.WriteString(strings.Repeat("-", 120) + "\n")
	for _, issue := range r.Issues {
		sb.WriteString(issue.File)
		sb.WriteString("\t")
		sb.WriteString(strconv.Itoa(issue.Line))
		sb.WriteString("\t")
		sb.WriteString(issue.Severity)
		sb.WriteString("\t")
		sb.WriteString(issue.Category)
		sb.WriteString("\t")
		sb.WriteString(issue.Rule)
		sb.WriteString("\t")
		sb.WriteString(issue.Message)
		sb.WriteString("\n")
	}
	return sb.String()
}