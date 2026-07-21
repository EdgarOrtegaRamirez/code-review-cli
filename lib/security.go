package codereview

import "strings"

// Security rules check for common security vulnerabilities.
type SecurityRule struct {
	Name       string
	Pattern    string
	Severity   string
	Message    string
	Recommendation string
}

var securityRules = []SecurityRule{
	{
		Name:       "hardcoded-password",
		Pattern:    `(?i)(password|passwd|pwd)\s*[=:]\s*["'][^"']+["']`,
		Severity:   "critical",
		Message:    "Hardcoded password detected",
		Recommendation: "Use environment variables or a secrets manager",
	},
	{
		Name:       "hardcoded-api-key",
		Pattern:    `(?i)(api[_-]?key|apikey|access[_-]?token|secret[_-]?key)\s*[=:]\s*["'][A-Za-z0-9_\-]{8,}["']`,
		Severity:   "critical",
		Message:    "Hardcoded API key or secret detected",
		Recommendation: "Use environment variables or a secrets manager",
	},
	{
		Name:       "sql-injection",
		Pattern:    `(?i)(SELECT|INSERT|UPDATE|DELETE|DROP)\s+.*(%s|%d|CONCAT\(|||)\s`,
		Severity:   "high",
		Message:    "Potential SQL injection vulnerability",
		Recommendation: "Use parameterized queries",
	},
	{
		Name:       "shell-injection",
		Pattern:    `(?i)(exec|system|popen|eval)\s*\(.*['"]\s*\+|['"]\s*\+\s*['"]`,
		Severity:   "high",
		Message:    "Potential shell injection vulnerability",
		Recommendation: "Use parameterized commands or avoid shell execution",
	},
	{
		Name:       "insecure-random",
		Pattern:    `(?i)rand\.Intn\s*\(|Math\.random\(\)|random\.random\(\)`,
		Severity:   "medium",
		Message:    "Insecure random number generation",
		Recommendation: "Use cryptographically secure random (crypto/rand, secrets module)",
	},
	{
		Name:       "weak-hash",
		Pattern:    `(?i)(md5|sha1)\s*\(|new\s+(MD5Digest|SHA1Digest)`,
		Severity:   "medium",
		Message:    "Weak cryptographic hash function",
		Recommendation: "Use SHA-256 or stronger hash function",
	},
	{
		Name:       "debug-mode",
		Pattern:    `(?i)(DEBUG|debug)\s*=\s*(true|True|1|on)`,
		Severity:   "high",
		Message:    "Debug mode enabled",
		Recommendation: "Disable debug mode in production",
	},
	{
		Name:       "hardcoded-credentials",
		Pattern:    `(?i)(aws_access_key|aws_secret|AZURE_KEY|GITHUB_TOKEN)\s*[=:]\s*["'][A-Za-z0-9_\-]{8,}["']`,
		Severity:   "critical",
		Message:    "Hardcoded cloud provider credentials detected",
		Recommendation: "Use IAM roles, instance profiles, or secrets manager",
	},
	{
		Name:       "eval-injection",
		Pattern:    `(?i)\beval\s*\(|exec\s*\(.*\+|Function\s*\(`,
		Severity:   "critical",
		Message:    "Dynamic code execution detected",
		Recommendation: "Avoid eval/exec; use safer alternatives",
	},
	{
		Name:       "path-traversal",
		Pattern:    `os\.path\.join\s*\(.*,\s*\w+\)|\.\./`,
		Severity:   "medium",
		Message:    "Potential path traversal vulnerability",
		Recommendation: "Validate and sanitize file paths",
	},
}

// CheckSecurity runs all security rules against a line of code.
func CheckSecurity(file string, lineNum int, line string) []Issue {
	var issues []Issue
	for _, rule := range securityRules {
		matched, _ := matchPattern(rule.Pattern, line)
		if matched {
			issues = append(issues, Issue{
				File:           file,
				Line:           lineNum,
				Severity:       rule.Severity,
				Category:       "security",
				Rule:           rule.Name,
				Message:        rule.Message,
				Recommendation: rule.Recommendation,
				Snippet:        strings.TrimSpace(line),
			})
		}
	}
	return issues
}