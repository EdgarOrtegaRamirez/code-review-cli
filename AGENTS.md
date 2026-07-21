# Code Review CLI — AI Agent Guide

## Overview

A CLI tool that analyzes code files or git diffs and generates structured review reports with security checks, style analysis, and complexity warnings.

## Building

```bash
go build -o codereview ./cmd/codereview
```

## Testing

```bash
go test ./...
```

## Running

```bash
# Analyze a file
codereview path/to/file.go

# Analyze a directory
codereview path/to/directory/

# Analyze git diff from stdin
git diff | codereview --diff -

# JSON output for CI
codereview file.go -f json

# Markdown report
codereview file.go -f markdown
```

## Architecture

- `cmd/codereview/main.go` — CLI entry point with Cobra
- `lib/types.go` — Issue and Report types
- `lib/security.go` — Security rule engine (10 rules)
- `lib/style.go` — Style and complexity analysis
- `lib/parser.go` — Git diff parser and orchestrator
- `lib/output.go` — JSON, Markdown, Table formatters
- `lib/analyzer.go` — File and directory analyzers
- `lib/walk.go` — File system traversal

## Adding a Security Rule

1. Add a struct to the `securityRules` slice in `lib/security.go`:
```go
{
    Name: "rule-name",
    Pattern: "regex-pattern",
    Severity: "medium",
    Message: "Description",
    Recommendation: "How to fix",
},
```
2. Run tests: `go test ./...`
3. Update README.md security rules table

## CI

The CI workflow runs:
- `go vet ./...`
- `go build ./...`
- `go test -race ./...`