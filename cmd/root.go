// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Version information
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

var rootCmd = &cobra.Command{
	Use:   "baseline-init",
	Short: "OpenSSF Baseline compliance tool",
	Long: `baseline-init is a CLI tool that helps repositories achieve and maintain
OpenSSF baseline compliance by:

- Checking repositories for missing compliance requirements
- Validating existing compliance files
- Auto-generating compliant default files
- Providing interactive setup guidance

For more information about OpenSSF baseline, visit:
https://github.com/ossf/security-baseline`,
	Version: fmt.Sprintf("%s (commit: %s, built: %s)", Version, GitCommit, BuildDate),
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetVersionTemplate(`{{.Version}}
`)
}
