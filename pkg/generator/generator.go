// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package generator

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/fatih/color"
)

// Generator handles creation of compliance files
type Generator struct {
	repoPath string
	force    bool
}

// Config contains configuration for file generation
type Config struct {
	ProjectURL              string
	ProjectName             string
	SecurityEmail           string
	AcceptsVulnReports      bool
	AcceptsPullRequests     bool
	AcceptsAutomatedPR      bool
	ProjectStage            string
	BugFixesOnly            bool
	Maintainers             []string
	DistributionPoints      []string
}

// New creates a new Generator instance
func New(repoPath string, force bool) *Generator {
	return &Generator{
		repoPath: repoPath,
		force:    force,
	}
}

// GenerateDefaults generates files with default values
func (g *Generator) GenerateDefaults() error {
	config := &Config{
		ProjectURL:              "https://github.com/example/repo",
		ProjectName:             filepath.Base(g.repoPath),
		SecurityEmail:           "security@example.com",
		AcceptsVulnReports:      true,
		AcceptsPullRequests:     true,
		AcceptsAutomatedPR:      true,
		ProjectStage:            "active",
		BugFixesOnly:            false,
		Maintainers:             []string{"github:maintainer"},
		DistributionPoints:      []string{},
	}

	return g.GenerateWithConfig(config)
}

// GenerateWithConfig generates files with provided configuration
func (g *Generator) GenerateWithConfig(config *Config) error {
	green := color.New(color.FgGreen).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()

	// Ensure .github directory exists
	githubDir := filepath.Join(g.repoPath, ".github")
	if err := os.MkdirAll(githubDir, 0755); err != nil {
		return fmt.Errorf("failed to create .github directory: %w", err)
	}

	// Generate SECURITY-INSIGHTS.yml
	siPath := filepath.Join(g.repoPath, "SECURITY-INSIGHTS.yml")
	if _, err := os.Stat(siPath); err == nil && !g.force {
		fmt.Printf("%s SECURITY-INSIGHTS.yml already exists (use --force to overwrite)\n", yellow("⚠"))
	} else {
		if err := g.generateSecurityInsights(siPath, config); err != nil {
			return fmt.Errorf("failed to generate SECURITY-INSIGHTS.yml: %w", err)
		}
		fmt.Printf("%s Generated SECURITY-INSIGHTS.yml\n", green("✓"))
	}

	// Generate SECURITY.md if it doesn't exist
	securityMdPath := filepath.Join(g.repoPath, "SECURITY.md")
	if _, err := os.Stat(securityMdPath); err == nil && !g.force {
		fmt.Printf("%s SECURITY.md already exists (use --force to overwrite)\n", yellow("⚠"))
	} else {
		if err := g.generateSecurityMd(securityMdPath, config); err != nil {
			return fmt.Errorf("failed to generate SECURITY.md: %w", err)
		}
		fmt.Printf("%s Generated SECURITY.md\n", green("✓"))
	}

	return nil
}

// generateSecurityInsights creates SECURITY-INSIGHTS.yml file
func (g *Generator) generateSecurityInsights(path string, config *Config) error {
	// Format dates as YYYY-MM-DD (schema 2.0.0 format)
	lastUpdated := time.Now().Format("2006-01-02")
	lastReviewed := time.Now().Format("2006-01-02")

	// Format maintainers for the new schema
	maintainersSection := formatMaintainersV2(config.Maintainers, config.SecurityEmail)

	content := fmt.Sprintf(`# OpenSSF Security Insights
# Schema version 2.0.0
# For more information, see: https://github.com/ossf/security-insights-spec

header:
  schema-version: 2.0.0
  last-updated: '%s'
  last-reviewed: '%s'
  url: %s
  comment: |
    This file provides security insights for the project.

project:
  name: %s
  administrators:
%s
  vulnerability-reporting:
    reports-accepted: %t
    bug-bounty-available: false

repository:
  url: %s
  status: %s
  accepts-change-request: %t
  accepts-automated-change-request: %t
  core-team:
%s
  license:
    url: %s/blob/main/LICENSE
    expression: Apache-2.0
  security:
    assessments:
      self:
        comment: |
          Self assessment has not yet been completed.
`, lastUpdated, lastReviewed, config.ProjectURL, config.ProjectName,
		maintainersSection, config.AcceptsVulnReports,
		config.ProjectURL, config.ProjectStage, config.AcceptsPullRequests,
		config.AcceptsAutomatedPR, maintainersSection, config.ProjectURL)

	return os.WriteFile(path, []byte(content), 0644)
}

// generateSecurityMd creates SECURITY.md file
func (g *Generator) generateSecurityMd(path string, config *Config) error {
	content := fmt.Sprintf(`# Security Policy

## Supported Versions

We release patches for security vulnerabilities. Which versions are eligible for
receiving such patches depends on the CVSS v3.0 Rating:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |

## Reporting a Vulnerability

Please report security vulnerabilities to: %s

We will acknowledge your email within 48 hours, and will send a more detailed response
within 7 days indicating the next steps in handling your report.

After the initial reply to your report, we will endeavor to keep you informed of the
progress being made towards a fix and full announcement.

## Disclosure Policy

When we receive a security bug report, we will:

1. Confirm the problem and determine the affected versions.
2. Audit code to find any potential similar problems.
3. Prepare fixes for all releases still under maintenance.

## Comments on this Policy

If you have suggestions on how this process could be improved, please submit a pull
request or open an issue.
`, config.SecurityEmail)

	return os.WriteFile(path, []byte(content), 0644)
}

// formatMaintainersList formats maintainers for YAML (legacy 1.0.0 format)
func formatMaintainersList(maintainers []string) string {
	if len(maintainers) == 0 {
		return "    - github:maintainer"
	}

	result := ""
	for _, m := range maintainers {
		result += fmt.Sprintf("    - %s\n", m)
	}
	return result[:len(result)-1] // Remove trailing newline
}

// formatMaintainersV2 formats maintainers for schema 2.0.0
func formatMaintainersV2(maintainers []string, email string) string {
	if len(maintainers) == 0 {
		return `    - name: Maintainer
      affiliation: Organization
      email: ` + email + `
      social: https://github.com/maintainer
      primary: true`
	}

	result := ""
	for i, m := range maintainers {
		// Extract username from github:username format
		username := m
		if len(m) > 7 && m[:7] == "github:" {
			username = m[7:]
		}

		primary := "false"
		if i == 0 {
			primary = "true"
		}

		result += fmt.Sprintf(`    - name: %s
      affiliation: Organization
      email: %s
      social: https://github.com/%s
      primary: %s
`, username, email, username, primary)
	}
	return result[:len(result)-1] // Remove trailing newline
}

// formatDistributionPoints formats distribution points for YAML
func formatDistributionPoints(points []string) string {
	if len(points) == 0 {
		return "  - https://github.com/example/repo/releases"
	}

	result := ""
	for _, p := range points {
		result += fmt.Sprintf("  - %s\n", p)
	}
	return result[:len(result)-1] // Remove trailing newline
}
