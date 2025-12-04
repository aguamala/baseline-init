// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/aguamala/baseline-init/pkg/generator"
	"github.com/aguamala/baseline-init/pkg/interactive"
	"github.com/spf13/cobra"
)

var (
	setupAuto        bool
	setupInteractive bool
	setupPath        string
	setupForce       bool
)

var setupCmd = &cobra.Command{
	Use:   "setup [path]",
	Short: "Setup OpenSSF baseline compliance files",
	Long: `Generate OpenSSF baseline compliance files for a repository.

The setup command can run in two modes:

1. Auto mode (--auto): Automatically generates files with sensible defaults
2. Interactive mode (--interactive): Walks you through customization

Example:
  baseline-init setup --auto
  baseline-init setup --interactive
  baseline-init setup --auto /path/to/repo
  baseline-init setup --auto --force  # Overwrite existing files`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)

	setupCmd.Flags().BoolVar(&setupAuto, "auto", false, "Auto-generate with defaults")
	setupCmd.Flags().BoolVar(&setupInteractive, "interactive", false, "Interactive setup mode")
	setupCmd.Flags().StringVarP(&setupPath, "path", "p", ".", "Path to repository")
	setupCmd.Flags().BoolVar(&setupForce, "force", false, "Overwrite existing files")

	setupCmd.MarkFlagsMutuallyExclusive("auto", "interactive")
}

func runSetup(cmd *cobra.Command, args []string) error {
	// Determine repository path
	repoPath := setupPath
	if len(args) > 0 {
		repoPath = args[0]
	}

	// Verify path exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("path does not exist: %s", repoPath)
	}

	// If neither mode specified, default to interactive
	if !setupAuto && !setupInteractive {
		setupInteractive = true
	}

	gen := generator.New(repoPath, setupForce)

	if setupInteractive {
		// Interactive mode: gather user input
		config, err := interactive.GatherConfiguration(repoPath)
		if err != nil {
			return fmt.Errorf("failed to gather configuration: %w", err)
		}

		if err := gen.GenerateWithConfig(config); err != nil {
			return fmt.Errorf("failed to generate files: %w", err)
		}
	} else {
		// Auto mode: generate with defaults
		if err := gen.GenerateDefaults(); err != nil {
			return fmt.Errorf("failed to generate files: %w", err)
		}
	}

	fmt.Println("\n✓ OpenSSF baseline compliance files generated successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. Review and customize the generated files")
	fmt.Println("  2. Run 'baseline-init check' to validate")
	fmt.Println("  3. Commit the files to your repository")

	return nil
}
