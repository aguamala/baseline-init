// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package checker

import (
	"os"
	"path/filepath"
)

// Checker performs OpenSSF baseline compliance checks
type Checker struct {
	repoPath string
}

// CheckResult contains the results of a compliance check
type CheckResult struct {
	Path          string             `json:"path"`
	IsCompliant   bool               `json:"is_compliant"`
	Files         []FileCheck        `json:"files"`
	MissingFiles  []string           `json:"missing_files"`
	Recommendations []Recommendation `json:"recommendations"`
}

// FileCheck represents the status of a compliance file
type FileCheck struct {
	Name     string   `json:"name"`
	Path     string   `json:"path"`
	Exists   bool     `json:"exists"`
	Valid    bool     `json:"valid"`
	Errors   []string `json:"errors,omitempty"`
	Warnings []string `json:"warnings,omitempty"`
}

// Recommendation provides actionable guidance
type Recommendation struct {
	Priority    string `json:"priority"` // critical, high, medium, low
	Category    string `json:"category"`
	Description string `json:"description"`
	Action      string `json:"action"`
}

// New creates a new Checker instance
func New(repoPath string) *Checker {
	return &Checker{
		repoPath: repoPath,
	}
}

// Check performs a compliance check on the repository
func (c *Checker) Check() (*CheckResult, error) {
	result := &CheckResult{
		Path:          c.repoPath,
		Files:         []FileCheck{},
		MissingFiles:  []string{},
		Recommendations: []Recommendation{},
	}

	// Check for SECURITY-INSIGHTS.yml
	siCheck := c.checkSecurityInsights()
	result.Files = append(result.Files, siCheck)
	if !siCheck.Exists {
		result.MissingFiles = append(result.MissingFiles, "SECURITY-INSIGHTS.yml")
		result.Recommendations = append(result.Recommendations, Recommendation{
			Priority:    "high",
			Category:    "Security Metadata",
			Description: "SECURITY-INSIGHTS.yml file is missing",
			Action:      "Run 'baseline-init setup --auto' to generate this file",
		})
	}

	// Check for SECURITY.md
	securityMdCheck := c.checkSecurityPolicy()
	result.Files = append(result.Files, securityMdCheck)
	if !securityMdCheck.Exists {
		result.MissingFiles = append(result.MissingFiles, "SECURITY.md")
		result.Recommendations = append(result.Recommendations, Recommendation{
			Priority:    "medium",
			Category:    "Security Policy",
			Description: "SECURITY.md file is missing",
			Action:      "Create a SECURITY.md file documenting your security policy",
		})
	}

	// Check for LICENSE file
	licenseCheck := c.checkLicense()
	result.Files = append(result.Files, licenseCheck)
	if !licenseCheck.Exists {
		result.MissingFiles = append(result.MissingFiles, "LICENSE")
		result.Recommendations = append(result.Recommendations, Recommendation{
			Priority:    "high",
			Category:    "Legal",
			Description: "LICENSE file is missing",
			Action:      "Add an appropriate open source license to your repository",
		})
	}

	// Check for CODE_OF_CONDUCT.md
	cocCheck := c.checkCodeOfConduct()
	result.Files = append(result.Files, cocCheck)
	if !cocCheck.Exists {
		result.Recommendations = append(result.Recommendations, Recommendation{
			Priority:    "medium",
			Category:    "Community",
			Description: "CODE_OF_CONDUCT.md file is missing",
			Action:      "Consider adding a code of conduct for contributors",
		})
	}

	// Check for CONTRIBUTING.md
	contributingCheck := c.checkContributing()
	result.Files = append(result.Files, contributingCheck)
	if !contributingCheck.Exists {
		result.Recommendations = append(result.Recommendations, Recommendation{
			Priority:    "low",
			Category:    "Community",
			Description: "CONTRIBUTING.md file is missing",
			Action:      "Consider adding contribution guidelines",
		})
	}

	// Determine overall compliance
	result.IsCompliant = len(result.MissingFiles) == 0

	return result, nil
}

// checkSecurityInsights checks for SECURITY-INSIGHTS.yml file
func (c *Checker) checkSecurityInsights() FileCheck {
	possiblePaths := []string{
		filepath.Join(c.repoPath, "SECURITY-INSIGHTS.yml"),
		filepath.Join(c.repoPath, ".github", "SECURITY-INSIGHTS.yml"),
		filepath.Join(c.repoPath, "SECURITY-INSIGHTS.yaml"),
		filepath.Join(c.repoPath, ".github", "SECURITY-INSIGHTS.yaml"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return FileCheck{
				Name:   "SECURITY-INSIGHTS.yml",
				Path:   path,
				Exists: true,
				Valid:  true, // TODO: Add actual validation
			}
		}
	}

	return FileCheck{
		Name:   "SECURITY-INSIGHTS.yml",
		Path:   "",
		Exists: false,
		Valid:  false,
	}
}

// checkSecurityPolicy checks for SECURITY.md file
func (c *Checker) checkSecurityPolicy() FileCheck {
	possiblePaths := []string{
		filepath.Join(c.repoPath, "SECURITY.md"),
		filepath.Join(c.repoPath, ".github", "SECURITY.md"),
		filepath.Join(c.repoPath, "docs", "SECURITY.md"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return FileCheck{
				Name:   "SECURITY.md",
				Path:   path,
				Exists: true,
				Valid:  true,
			}
		}
	}

	return FileCheck{
		Name:   "SECURITY.md",
		Path:   "",
		Exists: false,
		Valid:  false,
	}
}

// checkLicense checks for LICENSE file
func (c *Checker) checkLicense() FileCheck {
	possiblePaths := []string{
		filepath.Join(c.repoPath, "LICENSE"),
		filepath.Join(c.repoPath, "LICENSE.md"),
		filepath.Join(c.repoPath, "LICENSE.txt"),
		filepath.Join(c.repoPath, "COPYING"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return FileCheck{
				Name:   "LICENSE",
				Path:   path,
				Exists: true,
				Valid:  true,
			}
		}
	}

	return FileCheck{
		Name:   "LICENSE",
		Path:   "",
		Exists: false,
		Valid:  false,
	}
}

// checkCodeOfConduct checks for CODE_OF_CONDUCT.md file
func (c *Checker) checkCodeOfConduct() FileCheck {
	possiblePaths := []string{
		filepath.Join(c.repoPath, "CODE_OF_CONDUCT.md"),
		filepath.Join(c.repoPath, ".github", "CODE_OF_CONDUCT.md"),
		filepath.Join(c.repoPath, "docs", "CODE_OF_CONDUCT.md"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return FileCheck{
				Name:   "CODE_OF_CONDUCT.md",
				Path:   path,
				Exists: true,
				Valid:  true,
			}
		}
	}

	return FileCheck{
		Name:   "CODE_OF_CONDUCT.md",
		Path:   "",
		Exists: false,
		Valid:  false,
	}
}

// checkContributing checks for CONTRIBUTING.md file
func (c *Checker) checkContributing() FileCheck {
	possiblePaths := []string{
		filepath.Join(c.repoPath, "CONTRIBUTING.md"),
		filepath.Join(c.repoPath, ".github", "CONTRIBUTING.md"),
		filepath.Join(c.repoPath, "docs", "CONTRIBUTING.md"),
	}

	for _, path := range possiblePaths {
		if _, err := os.Stat(path); err == nil {
			return FileCheck{
				Name:   "CONTRIBUTING.md",
				Path:   path,
				Exists: true,
				Valid:  true,
			}
		}
	}

	return FileCheck{
		Name:   "CONTRIBUTING.md",
		Path:   "",
		Exists: false,
		Valid:  false,
	}
}
