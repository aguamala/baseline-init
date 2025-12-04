// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"fmt"
	"os"
	"strings"
	"time"

	sitooling "github.com/ossf/si-tooling/v2/si"
	"gopkg.in/yaml.v3"
)

// Validator validates compliance files
type Validator struct{}

// ValidationResult contains validation results
type ValidationResult struct {
	IsValid  bool     `json:"is_valid"`
	Errors   []string `json:"errors"`
	Warnings []string `json:"warnings"`
}

// SecurityInsights represents the SECURITY-INSIGHTS.yml structure (v1.0.0)
type SecurityInsightsV1 struct {
	Header struct {
		SchemaVersion  string `yaml:"schema-version"`
		ExpirationDate string `yaml:"expiration-date"`
		LastUpdated    string `yaml:"last-updated"`
		LastReviewed   string `yaml:"last-reviewed"`
		ProjectURL     string `yaml:"project-url"`
	} `yaml:"header"`
	ProjectLifecycle struct {
		Status       string `yaml:"status"`
		BugFixesOnly bool   `yaml:"bug-fixes-only"`
	} `yaml:"project-lifecycle"`
	ContributionPolicy struct {
		AcceptsPullRequests          bool `yaml:"accepts-pull-requests"`
		AcceptsAutomatedPullRequests bool `yaml:"accepts-automated-pull-requests"`
	} `yaml:"contribution-policy"`
	SecurityContacts []struct {
		Type  string `yaml:"type"`
		Value string `yaml:"value"`
	} `yaml:"security-contacts"`
	VulnerabilityReporting struct {
		AcceptsVulnerabilityReports bool `yaml:"accepts-vulnerability-reports"`
	} `yaml:"vulnerability-reporting"`
}

// SecurityInsightsV2 represents the SECURITY-INSIGHTS.yml structure (v2.0.0)
type SecurityInsightsV2 struct {
	Header struct {
		SchemaVersion interface{} `yaml:"schema-version"`
		LastUpdated   string      `yaml:"last-updated"`
		LastReviewed  string      `yaml:"last-reviewed"`
		URL           string      `yaml:"url"`
	} `yaml:"header"`
	Project struct {
		Name           string `yaml:"name"`
		Administrators []struct {
			Name  string `yaml:"name"`
			Email string `yaml:"email"`
		} `yaml:"administrators"`
		VulnerabilityReporting struct {
			ReportsAccepted bool `yaml:"reports-accepted"`
		} `yaml:"vulnerability-reporting"`
	} `yaml:"project"`
	Repository struct {
		URL                           string `yaml:"url"`
		Status                        string `yaml:"status"`
		AcceptsChangeRequest          bool   `yaml:"accepts-change-request"`
		AcceptsAutomatedChangeRequest bool   `yaml:"accepts-automated-change-request"`
	} `yaml:"repository"`
}

// New creates a new Validator instance
func New() *Validator {
	return &Validator{}
}

// ValidateFile validates a compliance file
func (v *Validator) ValidateFile(path string) (*ValidationResult, error) {
	// Read file
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Determine file type based on name
	filename := strings.ToLower(path)
	if strings.Contains(filename, "security-insights") {
		return v.validateSecurityInsights(data)
	}

	return nil, fmt.Errorf("unknown file type: %s", path)
}

// validateSecurityInsights validates SECURITY-INSIGHTS.yml
func (v *Validator) validateSecurityInsights(data []byte) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// First, detect schema version
	var header struct {
		Header struct {
			SchemaVersion interface{} `yaml:"schema-version"`
		} `yaml:"header"`
	}
	if err := yaml.Unmarshal(data, &header); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid YAML: %v", err))
		return result, nil
	}

	// Determine version and validate accordingly
	schemaVersion := fmt.Sprintf("%v", header.Header.SchemaVersion)

	if strings.HasPrefix(schemaVersion, "2.") {
		return v.validateSecurityInsightsV2(data)
	}

	// Default to v1 validation
	return v.validateSecurityInsightsV1(data)
}

// validateSecurityInsightsV1 validates SECURITY-INSIGHTS.yml schema v1.0.0
func (v *Validator) validateSecurityInsightsV1(data []byte) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	var si SecurityInsightsV1
	if err := yaml.Unmarshal(data, &si); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid YAML: %v", err))
		return result, nil
	}

	// Validate required fields
	if si.Header.SchemaVersion == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Missing required field: header.schema-version")
	}

	if si.Header.ProjectURL == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Missing required field: header.project-url")
	}

	if si.Header.ExpirationDate == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Missing required field: header.expiration-date")
	} else {
		// Validate expiration date format and check if expired
		expirationDate, err := time.Parse(time.RFC3339, si.Header.ExpirationDate)
		if err != nil {
			result.Warnings = append(result.Warnings, "Invalid expiration-date format (should be RFC3339)")
		} else if time.Now().After(expirationDate) {
			result.Warnings = append(result.Warnings, "File has expired - please update expiration-date")
		}
	}

	if si.Header.LastUpdated == "" {
		result.Warnings = append(result.Warnings, "Missing recommended field: header.last-updated")
	}

	if si.Header.LastReviewed == "" {
		result.Warnings = append(result.Warnings, "Missing recommended field: header.last-reviewed")
	}

	if si.ProjectLifecycle.Status == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Missing required field: project-lifecycle.status")
	} else {
		validStatuses := []string{"active", "archived", "concept", "moved", "wip"}
		isValid := false
		for _, status := range validStatuses {
			if si.ProjectLifecycle.Status == status {
				isValid = true
				break
			}
		}
		if !isValid {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Unusual project-lifecycle.status: %s (expected one of: %s)",
					si.ProjectLifecycle.Status, strings.Join(validStatuses, ", ")))
		}
	}

	if len(si.SecurityContacts) == 0 {
		result.Warnings = append(result.Warnings, "No security-contacts specified")
	} else {
		for i, contact := range si.SecurityContacts {
			if contact.Type == "" {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Security contact %d missing type", i))
			}
			if contact.Value == "" {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Security contact %d missing value", i))
			}
		}
	}

	return result, nil
}

// validateSecurityInsightsV2 validates SECURITY-INSIGHTS.yml schema v2.0.0
// Uses the official OpenSSF si-tooling library for schema validation
func (v *Validator) validateSecurityInsightsV2(data []byte) (*ValidationResult, error) {
	result := &ValidationResult{
		IsValid:  true,
		Errors:   []string{},
		Warnings: []string{},
	}

	// Use official si-tooling structs for validation
	var insights sitooling.SecurityInsights
	if err := yaml.Unmarshal(data, &insights); err != nil {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Schema validation failed: %v", err))
		return result, nil
	}

	// Validate schema version
	if !strings.HasPrefix(insights.Header.SchemaVersion, "2.") {
		result.IsValid = false
		result.Errors = append(result.Errors, fmt.Sprintf("Invalid schema version: %s (expected 2.x.x)", insights.Header.SchemaVersion))
		return result, nil
	}

	// insights is now a validated sitooling.SecurityInsights struct
	// Add our own custom checks on top of the official validation

	// Check header fields
	if insights.Header.LastUpdated == "" {
		result.Warnings = append(result.Warnings, "Missing recommended field: header.last-updated")
	}

	if insights.Header.LastReviewed == "" {
		result.Warnings = append(result.Warnings, "Missing recommended field: header.last-reviewed")
	}

	if insights.Header.URL == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Missing required field: header.url")
	}

	// Check project section
	if insights.Project.Name == "" {
		result.Warnings = append(result.Warnings, "Missing recommended field: project.name")
	}

	if len(insights.Project.Administrators) == 0 {
		result.Warnings = append(result.Warnings, "No project administrators specified")
	} else {
		for i, admin := range insights.Project.Administrators {
			if admin.Name == "" {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Administrator %d missing name", i))
			}
			if admin.Email == "" {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("Administrator %d missing email", i))
			}
		}
	}

	// Check repository section
	if insights.Repository.URL == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Missing required field: repository.url")
	}

	if insights.Repository.Status == "" {
		result.IsValid = false
		result.Errors = append(result.Errors, "Missing required field: repository.status")
	} else {
		validStatuses := []string{"active", "archived", "concept", "moved", "wip"}
		isValid := false
		for _, status := range validStatuses {
			if insights.Repository.Status == status {
				isValid = true
				break
			}
		}
		if !isValid {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("Unusual repository.status: %s (expected one of: %s)",
					insights.Repository.Status, strings.Join(validStatuses, ", ")))
		}
	}

	return result, nil
}
