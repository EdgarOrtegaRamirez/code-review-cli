package codereview

import (
	"strings"
)

// CheckStyle analyzes code for style issues.
func CheckStyle(file string, lineNum int, line string, lines []string) []Issue {
	var issues []Issue

	// Line too long
	if len(line) > 120 {
		issues = append(issues, Issue{
			File:     file,
			Line:     lineNum,
			Severity: "low",
			Category: "style",
			Rule:     "line-length",
			Message:  "Line exceeds 120 characters",
			Snippet:  strings.TrimSpace(line),
		})
	}

	// TODO/FIXME/HACK/WARN comment markers
	lower := strings.ToLower(line)
	if strings.Contains(lower, "todo") || strings.Contains(lower, "fixme") ||
		strings.Contains(lower, "hack") || strings.Contains(lower, "xxx") ||
		strings.Contains(lower, "deprecated") {
		issues = append(issues, Issue{
			File:     file,
			Line:     lineNum,
			Severity: "info",
			Category: "maintenance",
			Rule:     "todo-marker",
			Message:  "Todo marker found: " + strings.TrimSpace(line),
			Snippet:  strings.TrimSpace(line),
		})
	}

	// Check for print/debug statements
	if isDebugStatement(line) {
		issues = append(issues, Issue{
			File:     file,
			Line:     lineNum,
			Severity: "medium",
			Category: "quality",
			Rule:     "debug-statement",
			Message:  "Debug statement found",
			Recommendation: "Remove or use proper logging",
			Snippet:  strings.TrimSpace(line),
		})
	}

	return issues
}

func isDebugStatement(line string) bool {
	lower := strings.ToLower(strings.TrimSpace(line))
	debugPatterns := []string{
		"print(", "printf(", "println(", "console.log(", "log.debug",
		"log.Println", "log.Print(", "syslog.",
	}
	for _, p := range debugPatterns {
		if strings.Contains(lower, p) {
			return true
		}
	}
	return false
}

// CheckComplexity analyzes lines for complexity indicators.
func CheckComplexity(file string, lineNum int, line string, lines []string) []Issue {
	var issues []Issue

	// Deep nesting detection
	nesting := countNesting(lines[:lineNum])
	if nesting > 5 {
		issues = append(issues, Issue{
			File:     file,
			Line:     lineNum,
			Severity: "medium",
			Category: "complexity",
			Rule:     "deep-nesting",
			Message:  "Deep nesting detected (depth: " + string(rune(nesting+'0')) + ")",
			Recommendation: "Refactor to reduce nesting level",
			Snippet:  strings.TrimSpace(line),
		})
	}

	// Long line after whitespace trim
	trimmed := strings.TrimSpace(line)
	if len(trimmed) > 200 {
		issues = append(issues, Issue{
			File:     file,
			Line:     lineNum,
			Severity: "low",
			Category: "complexity",
			Rule:     "long-line",
			Message:  "Very long line (200+ chars)",
			Snippet:  trimmed,
		})
	}

	_ = strings.Map(func(r rune) rune {
		if r == ' ' || r == '	' || r == '\n' || r == '\r' {
			return 0
		}
		return r
	}, line)

	return issues
}

func countNesting(lines []string) int {
	depth := 0
	maxDepth := 0
	for _, line := range lines {
		for _, ch := range line {
			if ch == '{' {
				depth++
				if depth > maxDepth {
					maxDepth = depth
				}
			} else if ch == '}' {
				depth--
				if depth < 0 {
					depth = 0
				}
			}
		}
	}
	return maxDepth
}