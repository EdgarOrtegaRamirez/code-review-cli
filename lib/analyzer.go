package codereview

import (
	"os"
	"strings"
)

// AnalyzeFile analyzes a file and returns review issues.
func AnalyzeFile(filePath string) Report {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return Report{Issues: []Issue{{
			File:     filePath,
			Severity: "critical",
			Category: "tool",
			Rule:     "file-read-error",
			Message:  "Failed to read file: " + err.Error(),
		}}}
	}

	lines := strings.Split(string(data), "\n")
	var issues []Issue

	for i, line := range lines {
		if line == "" {
			continue
		}
		lineNum := i + 1

		securityIssues := CheckSecurity(filePath, lineNum, line)
		styleIssues := CheckStyle(filePath, lineNum, line, lines)
		complexityIssues := CheckComplexity(filePath, lineNum, line, lines)

		issues = append(issues, securityIssues...)
		issues = append(issues, styleIssues...)
		issues = append(issues, complexityIssues...)
	}

	severitySummary := map[string]int{
		"critical": 0, "high": 0, "medium": 0, "low": 0, "info": 0,
	}
	for _, issue := range issues {
		severitySummary[issue.Severity]++
	}

	return Report{
		TotalFiles:      1,
		TotalLines:      len(lines),
		Issues:          issues,
		SeveritySummary: severitySummary,
	}
}

// AnalyzeDirectory walks a directory and analyzes all files.
func AnalyzeDirectory(dirPath string) Report {
	var report Report
	var allIssues []Issue

	err := walkDir(dirPath, func(filePath string) {
		fileReport := AnalyzeFile(filePath)
		allIssues = append(allIssues, fileReport.Issues...)
		report.TotalFiles++
		report.TotalLines += fileReport.TotalLines
	})

	if err != nil {
		return Report{
			Issues: []Issue{{
				Severity: "critical",
				Category: "tool",
				Rule:     "dir-error",
				Message:  "Failed to walk directory: " + err.Error(),
			}},
		}
	}

	severitySummary := map[string]int{
		"critical": 0, "high": 0, "medium": 0, "low": 0, "info": 0,
	}
	for _, issue := range allIssues {
		severitySummary[issue.Severity]++
	}

	report.Issues = allIssues
	report.SeveritySummary = severitySummary
	return report
}