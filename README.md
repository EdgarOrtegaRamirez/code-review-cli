# Code Review CLI

A structured code review tool that analyzes code files or git diffs and generates reports with security checks, style analysis, and complexity warnings.

## Features

- **Security scanning**: Detects hardcoded passwords, API keys, SQL injection patterns, eval usage, and more
- **Style analysis**: Line length violations, TODO/FIXME markers, debug statements
- **Complexity warnings**: Deep nesting, overly long lines
- **Git diff analysis**: Analyze added lines from diffs with context awareness
- **Multiple output formats**: Table (human-readable), JSON (CI/CD integration), Markdown (reports)

## Install

```bash
go install github.com/EdgarOrtegaRamirez/code-review-cli/cmd/codereview@latest
```

Or download from [releases](https://github.com/EdgarOrtegaRamirez/code-review-cli/releases).

## Usage

### Analyze a file

```bash
codereview path/to/file.go
```

### Analyze a directory

```bash
codereview path/to/directory/
```

### Analyze a git diff

```bash
git diff HEAD~1 | codereview --diff -
codereview --diff diff.patch
```

### Output formats

```bash
# Table (default)
codereview file.go

# JSON (for CI/CD)
codereview file.go -f json

# Markdown report
codereview file.go -f markdown
```

## Security Rules

| Rule | Severity | Description |
|------|----------|-------------|
| hardcoded-password | critical | Detects passwords in source code |
| hardcoded-api-key | critical | Detects API keys and secrets |
| sql-injection | high | Detects SQL injection patterns |
| shell-injection | high | Detects shell execution with variables |
| eval-injection | critical | Detects dynamic code execution |
| hardcoded-credentials | critical | Cloud provider credentials |
| debug-mode | high | Debug mode enabled |
| insecure-random | medium | Weak random number generation |
| weak-hash | medium | MD5/SHA1 usage |
| path-traversal | medium | File path traversal risk |

## CI/CD Integration

Use JSON output for automated scanning:

```yaml
- name: Code Review Scan
  run: codereview ./... -f json > codereview-report.json

- name: Fail on critical issues
  run: |
    critical=$(jq '[.issues[] | select(.severity == "critical")] | length' codereview-report.json)
    if [ "$critical" -gt 0 ]; then echo "Critical issues found!"; exit 1; fi
```

## License

MIT