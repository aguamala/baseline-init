// Copyright 2025 baseline-init Authors
// SPDX-License-Identifier: Apache-2.0

package checker

import (
	"os"
	"path/filepath"
	"testing"
)

func TestChecker_Check(t *testing.T) {
	// Create temporary test directory
	tmpDir, err := os.MkdirTemp("", "baseline-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name            string
		setupFiles      map[string]string
		wantCompliant   bool
		wantMissingLen  int
	}{
		{
			name:           "empty repository",
			setupFiles:     map[string]string{},
			wantCompliant:  false,
			wantMissingLen: 3, // SECURITY-INSIGHTS.yml, SECURITY.md, LICENSE
		},
		{
			name: "only SECURITY-INSIGHTS.yml",
			setupFiles: map[string]string{
				"SECURITY-INSIGHTS.yml": "test content",
			},
			wantCompliant:  false,
			wantMissingLen: 2, // SECURITY.md, LICENSE
		},
		{
			name: "all required files",
			setupFiles: map[string]string{
				"SECURITY-INSIGHTS.yml": "test content",
				"SECURITY.md":           "security policy",
				"LICENSE":               "license content",
			},
			wantCompliant:  true,
			wantMissingLen: 0,
		},
		{
			name: "SECURITY-INSIGHTS.yml in .github",
			setupFiles: map[string]string{
				".github/SECURITY-INSIGHTS.yml": "test content",
				"SECURITY.md":                   "security policy",
				"LICENSE":                       "license content",
			},
			wantCompliant:  true,
			wantMissingLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test directory
			testDir := filepath.Join(tmpDir, tt.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatalf("Failed to create test dir: %v", err)
			}

			// Create test files
			for path, content := range tt.setupFiles {
				fullPath := filepath.Join(testDir, path)
				dir := filepath.Dir(fullPath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("Failed to create directory %s: %v", dir, err)
				}
				if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
					t.Fatalf("Failed to write file %s: %v", fullPath, err)
				}
			}

			// Run checker
			c := New(testDir)
			result, err := c.Check()
			if err != nil {
				t.Fatalf("Check() error = %v", err)
			}

			if result.IsCompliant != tt.wantCompliant {
				t.Errorf("IsCompliant = %v, want %v", result.IsCompliant, tt.wantCompliant)
			}

			if len(result.MissingFiles) != tt.wantMissingLen {
				t.Errorf("MissingFiles length = %d, want %d (missing: %v)",
					len(result.MissingFiles), tt.wantMissingLen, result.MissingFiles)
			}
		})
	}
}

func TestChecker_CheckSecurityInsights(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "baseline-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name       string
		setupPath  string
		wantExists bool
	}{
		{
			name:       "no file",
			setupPath:  "",
			wantExists: false,
		},
		{
			name:       "root SECURITY-INSIGHTS.yml",
			setupPath:  "SECURITY-INSIGHTS.yml",
			wantExists: true,
		},
		{
			name:       ".github SECURITY-INSIGHTS.yml",
			setupPath:  ".github/SECURITY-INSIGHTS.yml",
			wantExists: true,
		},
		{
			name:       "root SECURITY-INSIGHTS.yaml",
			setupPath:  "SECURITY-INSIGHTS.yaml",
			wantExists: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testDir := filepath.Join(tmpDir, tt.name)
			if err := os.MkdirAll(testDir, 0755); err != nil {
				t.Fatalf("Failed to create test dir: %v", err)
			}

			if tt.setupPath != "" {
				fullPath := filepath.Join(testDir, tt.setupPath)
				dir := filepath.Dir(fullPath)
				if err := os.MkdirAll(dir, 0755); err != nil {
					t.Fatalf("Failed to create directory: %v", err)
				}
				if err := os.WriteFile(fullPath, []byte("test"), 0644); err != nil {
					t.Fatalf("Failed to write file: %v", err)
				}
			}

			c := New(testDir)
			result := c.checkSecurityInsights()

			if result.Exists != tt.wantExists {
				t.Errorf("Exists = %v, want %v", result.Exists, tt.wantExists)
			}
		})
	}
}
