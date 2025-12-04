// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package report

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/aguamala/baseline-init/pkg/checker"
	"github.com/fatih/color"
	"gopkg.in/yaml.v3"
)

// Reporter handles formatting and output of compliance results
type Reporter struct {
	format string
}

// NewReporter creates a new Reporter instance
func NewReporter(format string) *Reporter {
	return &Reporter{
		format: format,
	}
}

// OutputCheckResult outputs the compliance check result
func (r *Reporter) OutputCheckResult(result *checker.CheckResult) error {
	switch r.format {
	case "json":
		return r.outputJSON(result)
	case "yaml":
		return r.outputYAML(result)
	case "text":
		return r.outputText(result)
	default:
		return fmt.Errorf("unsupported format: %s", r.format)
	}
}

// outputJSON outputs results as JSON
func (r *Reporter) outputJSON(result *checker.CheckResult) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}

// outputYAML outputs results as YAML
func (r *Reporter) outputYAML(result *checker.CheckResult) error {
	encoder := yaml.NewEncoder(os.Stdout)
	defer encoder.Close()
	return encoder.Encode(result)
}

// outputText outputs results as human-readable text
func (r *Reporter) outputText(result *checker.CheckResult) error {
	// Colors
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	cyan := color.New(color.FgCyan).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	// Header
	fmt.Println(bold("OpenSSF Baseline Compliance Check"))
	fmt.Println(strings.Repeat("=", 50))
	fmt.Printf("Repository: %s\n\n", result.Path)

	// Overall status
	if result.IsCompliant {
		fmt.Printf("Status: %s\n\n", green("✓ COMPLIANT"))
	} else {
		fmt.Printf("Status: %s\n\n", red("✗ NOT COMPLIANT"))
	}

	// File checks
	fmt.Println(bold("File Checks:"))
	for _, file := range result.Files {
		if file.Exists {
			fmt.Printf("  %s %s\n", green("✓"), file.Name)
			if file.Path != "" {
				fmt.Printf("    Location: %s\n", cyan(file.Path))
			}
			if len(file.Warnings) > 0 {
				for _, warning := range file.Warnings {
					fmt.Printf("    %s %s\n", yellow("⚠"), warning)
				}
			}
		} else {
			fmt.Printf("  %s %s\n", red("✗"), file.Name)
		}
	}
	fmt.Println()

	// Missing files
	if len(result.MissingFiles) > 0 {
		fmt.Println(bold("Missing Files:"))
		for _, missing := range result.MissingFiles {
			fmt.Printf("  %s %s\n", red("✗"), missing)
		}
		fmt.Println()
	}

	// Recommendations
	if len(result.Recommendations) > 0 {
		fmt.Println(bold("Recommendations:"))

		// Group by priority
		priorities := []string{"critical", "high", "medium", "low"}
		for _, priority := range priorities {
			var recs []checker.Recommendation
			for _, rec := range result.Recommendations {
				if rec.Priority == priority {
					recs = append(recs, rec)
				}
			}

			if len(recs) == 0 {
				continue
			}

			priorityColor := color.New(color.FgWhite).SprintFunc()
			switch priority {
			case "critical":
				priorityColor = color.New(color.FgRed, color.Bold).SprintFunc()
			case "high":
				priorityColor = color.New(color.FgRed).SprintFunc()
			case "medium":
				priorityColor = color.New(color.FgYellow).SprintFunc()
			case "low":
				priorityColor = color.New(color.FgCyan).SprintFunc()
			}

			for _, rec := range recs {
				fmt.Printf("\n  [%s] %s\n", priorityColor(strings.ToUpper(priority)), bold(rec.Description))
				fmt.Printf("  Category: %s\n", rec.Category)
				fmt.Printf("  Action: %s\n", cyan(rec.Action))
			}
		}
		fmt.Println()
	}

	// Summary
	if !result.IsCompliant {
		fmt.Println(bold("Next Steps:"))
		fmt.Println("  1. Run 'baseline-init setup --auto' to auto-generate missing files")
		fmt.Println("  2. Or run 'baseline-init setup --interactive' for guided setup")
		fmt.Println("  3. Review and customize generated files")
		fmt.Println("  4. Run 'baseline-init check' again to verify")
	}

	return nil
}
