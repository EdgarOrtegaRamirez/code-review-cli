package codereview

import (
	"regexp"
	"strings"
)

func matchPattern(pattern, input string) (bool, string) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return false, ""
	}
	matches := re.FindStringSubmatch(input)
	if len(matches) > 0 {
		return true, matches[0]
	}
	return false, ""
}

// DiffLine represents a line from a git diff.
type DiffLine struct {
	Op      string // "add", "remove", "context"
	NewLine int
	OldLine int
	Content string
}

// ParseDiff parses a git diff string into individual lines.
func ParseDiff(diff string) []DiffLine {
	var lines []DiffLine
	currentOld := 0
	currentNew := 0
	var filePath string

	linesList := strings.Split(diff, "\n")

	for _, rawLine := range linesList {
		// File header
		if strings.HasPrefix(rawLine, "diff --git") {
			re := regexp.MustCompile(`diff --git a/(.+) b/.+`)
			matches := re.FindStringSubmatch(rawLine)
			if len(matches) > 1 {
				filePath = matches[1]
			}
			continue
		}

		// Index line
		if strings.HasPrefix(rawLine, "index ") {
			continue
		}

		// --- / +++ lines
		if strings.HasPrefix(rawLine, "--- ") || strings.HasPrefix(rawLine, "+++ ") {
			continue
		}

		// Hunk header
		if strings.HasPrefix(rawLine, "@@") {
			re := regexp.MustCompile(`@@ -(\d+)(?:,(\d+))? \+(\d+)(?:,(\d+))? @@`)
			matches := re.FindStringSubmatch(rawLine)
			if len(matches) > 3 {
				currentOld = atoi(matches[1])
				currentNew = atoi(matches[3])
			}
			continue
		}

		// Content lines
		if len(rawLine) == 0 {
			continue
		}

		line := DiffLine{
			Content: rawLine,
		}

		switch rawLine[0] {
		case '+':
			line.Op = "add"
			line.NewLine = currentNew
			line.OldLine = currentOld
			currentNew++
		case '-':
			line.Op = "remove"
			line.NewLine = currentNew
			line.OldLine = currentOld
			currentOld++
		case ' ':
			line.Op = "context"
			line.NewLine = currentNew
			line.OldLine = currentOld
			currentNew++
			currentOld++
		case '\\':
			continue
		default:
			line.Op = "context"
			line.OldLine = currentOld
		}

		if filePath != "" {
			line.Content = strings.TrimPrefix(rawLine, rawLine[:1])
		}

		lines = append(lines, line)
	}

	return lines
}

func atoi(s string) int {
	var n int
	for _, c := range s {
		if c >= '0' && c <= '9' {
			n = n*10 + int(c-'0')
		}
	}
	return n
}

// AnalyzeDiff analyzes a git diff and returns review issues.
func AnalyzeDiff(diff string) Report {
	diffLines := ParseDiff(diff)
	var allIssues []Issue
	files := make(map[string]bool)
	var currentFile string

	for _, dl := range diffLines {
		if dl.Op == "context" {
			continue
		}

		// Track file from content
		filePath := currentFile
		if dl.Op == "add" {
			// CheckStyle expects []string, convert from DiffLine
					var contextLines []string
					for _, d := range diffLines {
						if d.Op == "context" {
							contextLines = append(contextLines, d.Content)
						}
					}
					securityIssues := CheckSecurity(filePath, dl.NewLine, dl.Content)
					styleIssues := CheckStyle(filePath, dl.NewLine, dl.Content, contextLines)
					complexityIssues := CheckComplexity(filePath, dl.NewLine, dl.Content, contextLines)

			allIssues = append(allIssues, securityIssues...)
			allIssues = append(allIssues, styleIssues...)
			allIssues = append(allIssues, complexityIssues...)
		}

		_ = files
		_ = dl.Content
	}

	severitySummary := map[string]int{
		"critical": 0,
		"high":     0,
		"medium":   0,
		"low":      0,
		"info":     0,
	}
	for _, issue := range allIssues {
		severitySummary[issue.Severity]++
	}

	return Report{
		Issues:          allIssues,
		SeveritySummary: severitySummary,
		TotalLines:      len(diffLines),
	}
}