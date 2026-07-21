package codereview

import "strings"

// FormatJSON renders a report as JSON string.
func FormatJSON(r Report) string {
	return toJson(r)
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
		"critical": "🔴 Critical",
		"high":     "🟠 High",
		"medium":   "🟡 Medium",
		"low":      "🔵 Low",
		"info":     "ℹ️ Info",
	}
	total := 0
	for _, s := range severities {
		count := r.SeveritySummary[s]
		total += count
		sb.WriteString("| ")
		sb.WriteString(labels[s])
		sb.WriteString(" | ")
		sb.WriteString(strings.Repeat("⚠", count))
		sb.WriteString(" |\n")
	}
	sb.WriteString("| **Total** | **")
	sb.WriteString(string(rune(total+'0')))
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
		sb.WriteString("`")
		sb.WriteString(" - Line: ")
		sb.WriteString(string(rune(issue.Line+'0')))
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
	sb.WriteString(strings.Repeat("-", 100) + "\n")
	for _, issue := range r.Issues {
		sb.WriteString(issue.File)
		sb.WriteString("\t")
		sb.WriteString(string(rune(issue.Line+'0')))
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

func toJson(r Report) string {
	var sb strings.Builder
	sb.WriteString("{\n")
	sb.WriteString("  \"total_lines\": " + string(rune(r.TotalLines+'0')) + ",\n")
	sb.WriteString("  \"issues_count\": " + string(rune(len(r.Issues)+'0')) + ",\n")
	sb.WriteString("  \"issues\": [\n")
	for i, issue := range r.Issues {
		sb.WriteString("    {\n")
		sb.WriteString("      \"file\": \"" + issue.File + "\",\n")
		sb.WriteString("      \"line\": " + string(rune(issue.Line+'0')) + ",\n")
		sb.WriteString("      \"severity\": \"" + issue.Severity + "\",\n")
		sb.WriteString("      \"category\": \"" + issue.Category + "\",\n")
		sb.WriteString("      \"rule\": \"" + issue.Rule + "\",\n")
		sb.WriteString("      \"message\": \"" + issue.Message + "\"")
		if issue.Recommendation != "" {
			sb.WriteString(", \"recommendation\": \"" + issue.Recommendation + "\"")
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