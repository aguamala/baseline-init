// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/aguamala/baseline-init/pkg/validator"
	"github.com/spf13/cobra"
)

var validateCmd = &cobra.Command{
	Use:   "validate <file>",
	Short: "Validate a compliance file against its schema",
	Long: `Validate OpenSSF compliance files (like SECURITY-INSIGHTS.yml)
against their official schemas.

Example:
  baseline-init validate SECURITY-INSIGHTS.yml
  baseline-init validate .github/SECURITY-INSIGHTS.yml`,
	Args: cobra.ExactArgs(1),
	RunE: runValidate,
}

func init() {
	rootCmd.AddCommand(validateCmd)
}

func runValidate(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	// Verify file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("file does not exist: %s", filePath)
	}

	// Validate the file
	v := validator.New()
	result, err := v.ValidateFile(filePath)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	if result.IsValid {
		fmt.Printf("✓ %s is valid\n", filePath)
		return nil
	}

	fmt.Printf("✗ %s is invalid:\n", filePath)
	for _, e := range result.Errors {
		fmt.Printf("  - %s\n", e)
	}

	if len(result.Warnings) > 0 {
		fmt.Println("\nWarnings:")
		for _, w := range result.Warnings {
			fmt.Printf("  - %s\n", w)
		}
	}

	os.Exit(1)
	return nil
}
