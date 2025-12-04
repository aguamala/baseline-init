// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/aguamala/baseline-init/pkg/checker"
	"github.com/aguamala/baseline-init/pkg/report"
	"github.com/spf13/cobra"
)

var (
	checkOutputFormat string
	checkPath         string
)

var checkCmd = &cobra.Command{
	Use:   "check [path]",
	Short: "Check repository for OpenSSF baseline compliance",
	Long: `Scan a repository and identify missing or invalid OpenSSF baseline
compliance requirements.

The check command will:
- Look for required compliance files (SECURITY-INSIGHTS.yml, etc.)
- Validate existing files against schemas
- Report what's missing or needs fixing
- Provide actionable recommendations

Example:
  baseline-init check
  baseline-init check /path/to/repo
  baseline-init check --format json
  baseline-init check --format yaml`,
	Args: cobra.MaximumNArgs(1),
	RunE: runCheck,
}

func init() {
	rootCmd.AddCommand(checkCmd)

	checkCmd.Flags().StringVarP(&checkOutputFormat, "format", "f", "text", "Output format (text, json, yaml)")
	checkCmd.Flags().StringVarP(&checkPath, "path", "p", ".", "Path to repository")
}

func runCheck(cmd *cobra.Command, args []string) error {
	// Determine repository path
	repoPath := checkPath
	if len(args) > 0 {
		repoPath = args[0]
	}

	// Verify path exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", repoPath)
	}

	// Run compliance check
	c := checker.New(repoPath)
	result, err := c.Check()
	if err != nil {
		return fmt.Errorf("compliance check failed: %w", err)
	}

	// Format and output results
	reporter := report.NewReporter(checkOutputFormat)
	if err := reporter.OutputCheckResult(result); err != nil {
		return fmt.Errorf("failed to output results: %w", err)
	}

	// Exit with error code if not compliant
	if !result.IsCompliant {
		os.Exit(1)
	}

	return nil
}
