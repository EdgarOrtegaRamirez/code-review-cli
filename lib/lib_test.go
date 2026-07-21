package codereview

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestCheckSecurity_Patterns(t *testing.T) {
	tests := []struct {
		name     string
		line     string
		file     string
		lineNum  int
		wantLen  int
	}{
		{"hardcoded password", "password = 'secret123'", "", 1, 1},
		{"hardcoded API key", "api_key = 'abcdef1234567890'", "", 2, 1},
		{"debug mode", "DEBUG = True", "", 3, 1},
		{"eval injection", "eval(user_input)", "", 4, 1},
		{"safe code", "x = 1 + 1", "", 5, 0},
		{"safe variable", "password_hash = compute_hash()", "", 6, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := CheckSecurity(tt.file, tt.lineNum, tt.line)
			if len(issues) != tt.wantLen {
				t.Errorf("expected %d issues, got %d: %v", tt.wantLen, len(issues), issues)
			}
		})
	}
}

func TestCheckStyle_LongLine(t *testing.T) {
	longLine := strings.Repeat("a", 121)
	issues := CheckStyle("test.go", 1, longLine, nil)
	found := false
	for _, issue := range issues {
		if issue.Rule == "line-length" {
			found = true
		}
	}
	if !found {
		t.Error("expected line-length warning for 121-char line")
	}
}

func TestCheckStyle_NoLongLine(t *testing.T) {
	issues := CheckStyle("test.go", 1, "short line", nil)
	for _, issue := range issues {
		if issue.Rule == "line-length" {
			t.Error("unexpected line-length warning")
		}
	}
}

func TestAnalyzeFile_NotFound(t *testing.T) {
	report := AnalyzeFile("/nonexistent/path/file.py")
	if len(report.Issues) == 0 {
		t.Error("expected error issue for nonexistent file")
	}
}

func TestAnalyzeFile_NewGoFile(t *testing.T) {
	dir := t.TempDir()
	testFile := filepath.Join(dir, "test.go")
	os.WriteFile(testFile, []byte("package main\n\nfunc main() {\n\tprint(\"hello\")\n}\n"), 0644)

	report := AnalyzeFile(testFile)
	if report.TotalFiles != 1 {
		t.Errorf("expected 1 file, got %d", report.TotalFiles)
	}
	if report.TotalLines != 6 {
		t.Errorf("expected 5 lines, got %d", report.TotalLines)
	}
}

func TestFormatJSON(t *testing.T) {
	report := Report{
		TotalLines: 10,
		Issues: []Issue{
			{File: "test.go", Line: 5, Severity: "high", Category: "security", Rule: "test-rule", Message: "test message"},
		},
		SeveritySummary: map[string]int{"high": 1},
	}

	jsonStr := FormatJSON(report)
	if !strings.Contains(jsonStr, `"total_lines"`) {
		t.Error("JSON missing total_lines")
	}
	if !strings.Contains(jsonStr, `"issues"`) {
		t.Error("JSON missing issues")
	}
	if !strings.Contains(jsonStr, `"test message"`) {
		t.Error("JSON missing message")
	}
}

func TestFormatMarkdown(t *testing.T) {
	report := Report{
		Issues: []Issue{
			{File: "test.go", Line: 5, Severity: "high", Category: "security", Rule: "test-rule", Message: "test issue"},
		},
		SeveritySummary: map[string]int{"high": 1},
	}

	mdStr := FormatMarkdown(report)
	if !strings.Contains(mdStr, "Code Review Report") {
		t.Error("markdown missing heading")
	}
	if !strings.Contains(mdStr, "high") {
		t.Error("markdown missing severity")
	}
}

func TestParseDiff(t *testing.T) {
	diff := `diff --git a/test.py b/test.py
index abc123..def456 100644
--- a/test.py
+++ b/test.py
@@ -1,3 +1,4 @@
 x = 1
+password = 'secret123'
 y = 2
 z = 3
`
	lines := ParseDiff(diff)
	addLines := 0
	for _, l := range lines {
		if l.Op == "add" {
			addLines++
		}
	}
	if addLines != 1 {
		t.Errorf("expected 1 add line, got %d", addLines)
	}
}

func TestCheckComplexity(t *testing.T) {
	lines := []string{
		"func main() {",
		"    if true {",
		"        if true {",
		"            if true {",
		"                if true {",
		"                    if true {",
		"                        x := 1",
		"                    }",
		"                }",
		"            }",
		"        }",
		"    }",
		"}",
	}
	// Line 7 (index 6, 1-based: 7) has depth > 5
	issues := CheckComplexity("test.go", 7, "x := 1", lines)
	found := false
	for _, issue := range issues {
		if issue.Rule == "deep-nesting" {
			found = true
		}
	}
	if !found {
		t.Error("expected deep-nesting warning")
	}
}

func TestWalkDir(t *testing.T) {
	dir := t.TempDir()
	os.WriteFile(filepath.Join(dir, "test.go"), []byte("package main"), 0644)
	os.WriteFile(filepath.Join(dir, "test.py"), []byte("print(1)"), 0644)
	os.WriteFile(filepath.Join(dir, "test.txt"), []byte("hello"), 0644)

	var found []string
	err := walkDir(dir, func(path string) { found = append(found, path) })
	if err != nil {
		t.Fatal(err)
	}

	if len(found) != 2 {
		t.Errorf("expected 2 files, got %d: %v", len(found), found)
	}
}