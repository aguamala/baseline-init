// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package validator

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestValidator_ValidateSecurityInsights(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		wantValid   bool
		wantErrors  int
		wantWarnings int
	}{
		{
			name: "valid minimal file",
			content: `header:
  schema-version: '1.0.0'
  expiration-date: '2026-12-31T23:59:59Z'
  last-updated: '2025-01-01T00:00:00Z'
  last-reviewed: '2025-01-01T00:00:00Z'
  project-url: https://github.com/example/repo

project-lifecycle:
  status: active
  bug-fixes-only: false

contribution-policy:
  accepts-pull-requests: true
  accepts-automated-pull-requests: true

security-contacts:
  - type: email
    value: security@example.com

vulnerability-reporting:
  accepts-vulnerability-reports: true
`,
			wantValid:    true,
			wantErrors:   0,
			wantWarnings: 0,
		},
		{
			name: "missing required fields",
			content: `header:
  schema-version: ''
  expiration-date: ''
  project-url: ''

project-lifecycle:
  status: ''
`,
			wantValid:    false,
			wantErrors:   4, // schema-version, project-url, expiration-date, status
			wantWarnings: 0,
		},
		{
			name: "missing recommended fields",
			content: `header:
  schema-version: '1.0.0'
  expiration-date: '2026-12-31T23:59:59Z'
  project-url: https://github.com/example/repo

project-lifecycle:
  status: active
`,
			wantValid:    true,
			wantErrors:   0,
			wantWarnings: 3, // last-updated, last-reviewed, no security-contacts
		},
		{
			name: "expired file",
			content: `header:
  schema-version: '1.0.0'
  expiration-date: '2020-01-01T00:00:00Z'
  last-updated: '2020-01-01T00:00:00Z'
  last-reviewed: '2020-01-01T00:00:00Z'
  project-url: https://github.com/example/repo

project-lifecycle:
  status: active

security-contacts:
  - type: email
    value: security@example.com
`,
			wantValid:    true,
			wantErrors:   0,
			wantWarnings: 1, // expired
		},
		{
			name:         "invalid YAML",
			content:      `this is not: valid: yaml:`,
			wantValid:    false,
			wantErrors:   1,
			wantWarnings: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := New()
			result, err := v.validateSecurityInsights([]byte(tt.content))
			if err != nil {
				t.Fatalf("validateSecurityInsights() error = %v", err)
			}

			if result.IsValid != tt.wantValid {
				t.Errorf("IsValid = %v, want %v (errors: %v, warnings: %v)",
					result.IsValid, tt.wantValid, result.Errors, result.Warnings)
			}

			if len(result.Errors) != tt.wantErrors {
				t.Errorf("Errors count = %d, want %d (errors: %v)",
					len(result.Errors), tt.wantErrors, result.Errors)
			}

			if len(result.Warnings) < tt.wantWarnings {
				t.Errorf("Warnings count = %d, want at least %d (warnings: %v)",
					len(result.Warnings), tt.wantWarnings, result.Warnings)
			}
		})
	}
}

func TestValidator_ValidateFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validator-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	validContent := `header:
  schema-version: '1.0.0'
  expiration-date: '` + time.Now().AddDate(1, 0, 0).Format(time.RFC3339) + `'
  last-updated: '` + time.Now().Format(time.RFC3339) + `'
  last-reviewed: '` + time.Now().Format(time.RFC3339) + `'
  project-url: https://github.com/example/repo

project-lifecycle:
  status: active

security-contacts:
  - type: email
    value: security@example.com
`

	testFile := filepath.Join(tmpDir, "SECURITY-INSIGHTS.yml")
	if err := os.WriteFile(testFile, []byte(validContent), 0644); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	v := New()
	result, err := v.ValidateFile(testFile)
	if err != nil {
		t.Fatalf("ValidateFile() error = %v", err)
	}

	if !result.IsValid {
		t.Errorf("IsValid = false, want true (errors: %v)", result.Errors)
	}
}
