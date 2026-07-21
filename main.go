package main

import (
	"fmt"
	"os"

	"github.com/EdgarOrtegaRamirez/code-review-cli"
	"github.com/spf13/cobra"
)

var (
	format   string
	output   string
	diffFile string
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "codereview [flags] PATH",
		Short: "Code Review CLI — structured code review reports",
		Long:  "A CLI tool that analyzes code files or git diffs and generates structured review reports with security checks, style analysis, and complexity warnings.",
		RunE: func(cmd *cobra.Command, args []string) error {
			if diffFile != "" {
				return runDiff(diffFile)
			}
			if len(args) > 0 {
				path := args[0]
				info, err := os.Stat(path)
				if err != nil {
					return fmt.Errorf("cannot access %s: %w", path, err)
				}
				if info.IsDir() {
					return runDir(path)
				}
				return runFile(path)
			}
			return cmd.Help()
		},
	}

	rootCmd.Flags().StringVarP(&format, "format", "f", "table", "Output format: table, json, markdown")
	rootCmd.Flags().StringVarP(&output, "output", "o", "", "Write output to file")
	rootCmd.Flags().StringVarP(&diffFile, "diff", "d", "", "Read git diff from file")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runFile(path string) error {
	report := codereview.AnalyzeFile(path)
	outputStr := formatReport(report)
	return writeOutput(outputStr)
}

func runDir(path string) error {
	report := codereview.AnalyzeDirectory(path)
	outputStr := formatReport(report)
	return writeOutput(outputStr)
}

func runDiff(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("cannot read diff file: %w", err)
	}
	report := codereview.AnalyzeDiff(string(data))
	outputStr := formatReport(report)
	return writeOutput(outputStr)
}

func formatReport(report codereview.Report) string {
	switch format {
	case "json":
		return codereview.FormatJSON(report)
	case "markdown":
		return codereview.FormatMarkdown(report)
	default:
		return codereview.FormatTable(report)
	}
}

func writeOutput(outputStr string) error {
	if output != "" {
		return os.WriteFile(output, []byte(outputStr), 0644)
	}
	fmt.Print(outputStr)
	return nil
}