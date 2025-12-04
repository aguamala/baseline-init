// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package interactive

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/aguamala/baseline-init/pkg/generator"
	"github.com/manifoldco/promptui"
)

// GatherConfiguration interactively gathers configuration from the user
func GatherConfiguration(repoPath string) (*generator.Config, error) {
	config := &generator.Config{}

	fmt.Println("🔧 OpenSSF Baseline Interactive Setup")
	fmt.Println("======================================")
	fmt.Println()

	// Project URL
	projectURL, err := detectGitRemote(repoPath)
	if err != nil {
		projectURL = ""
	}

	urlPrompt := promptui.Prompt{
		Label:   "Project URL",
		Default: projectURL,
	}
	config.ProjectURL, err = urlPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	// Project Name
	projectName := filepath.Base(repoPath)
	namePrompt := promptui.Prompt{
		Label:   "Project Name",
		Default: projectName,
	}
	config.ProjectName, err = namePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	// Security Email
	emailPrompt := promptui.Prompt{
		Label:   "Security Contact Email",
		Default: "security@example.com",
		Validate: func(input string) error {
			if !strings.Contains(input, "@") {
				return fmt.Errorf("invalid email address")
			}
			return nil
		},
	}
	config.SecurityEmail, err = emailPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	// Project Stage
	stagePrompt := promptui.Select{
		Label: "Project Lifecycle Stage",
		Items: []string{"active", "archived", "concept", "moved", "wip"},
	}
	_, config.ProjectStage, err = stagePrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	// Accepts Vulnerability Reports
	vulnPrompt := promptui.Select{
		Label: "Accept Vulnerability Reports",
		Items: []string{"Yes", "No"},
	}
	_, vulnResponse, err := vulnPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}
	config.AcceptsVulnReports = vulnResponse == "Yes"

	// Accepts Pull Requests
	prPrompt := promptui.Select{
		Label: "Accept Pull Requests",
		Items: []string{"Yes", "No"},
	}
	_, prResponse, err := prPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}
	config.AcceptsPullRequests = prResponse == "Yes"

	// Accepts Automated PRs
	autoPrPrompt := promptui.Select{
		Label: "Accept Automated Pull Requests (e.g., Dependabot)",
		Items: []string{"Yes", "No"},
	}
	_, autoPrResponse, err := autoPrPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}
	config.AcceptsAutomatedPR = autoPrResponse == "Yes"

	// Bug Fixes Only
	bugFixPrompt := promptui.Select{
		Label: "Bug Fixes Only (no new features)",
		Items: []string{"No", "Yes"},
	}
	_, bugFixResponse, err := bugFixPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}
	config.BugFixesOnly = bugFixResponse == "Yes"

	// Maintainers
	maintainerPrompt := promptui.Prompt{
		Label:   "GitHub Maintainer Username(s) (comma-separated)",
		Default: "maintainer",
	}
	maintainerInput, err := maintainerPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	maintainers := strings.Split(maintainerInput, ",")
	config.Maintainers = []string{}
	for _, m := range maintainers {
		m = strings.TrimSpace(m)
		if m != "" {
			if !strings.HasPrefix(m, "github:") {
				m = "github:" + m
			}
			config.Maintainers = append(config.Maintainers, m)
		}
	}

	// Distribution Points
	distPrompt := promptui.Prompt{
		Label:   "Distribution Points (URLs, comma-separated, or press Enter to skip)",
		Default: "",
	}
	distInput, err := distPrompt.Run()
	if err != nil {
		return nil, fmt.Errorf("prompt failed: %w", err)
	}

	if distInput != "" {
		distPoints := strings.Split(distInput, ",")
		config.DistributionPoints = []string{}
		for _, d := range distPoints {
			d = strings.TrimSpace(d)
			if d != "" {
				config.DistributionPoints = append(config.DistributionPoints, d)
			}
		}
	}

	fmt.Println()
	return config, nil
}

// detectGitRemote attempts to detect the Git remote URL
func detectGitRemote(repoPath string) (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = repoPath
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	url := strings.TrimSpace(string(output))

	// Convert SSH URL to HTTPS
	if strings.HasPrefix(url, "git@github.com:") {
		url = strings.Replace(url, "git@github.com:", "https://github.com/", 1)
		url = strings.TrimSuffix(url, ".git")
	}

	return url, nil
}
