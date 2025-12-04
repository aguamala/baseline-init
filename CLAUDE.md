# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

`baseline-init` is a Golang CLI tool that helps repositories achieve OpenSSF Security Baseline compliance. It checks for required security files, validates them against schemas, and auto-generates compliant files with sensible defaults.

## Build & Development Commands

### Building
```bash
# Standard build
go build -o baseline-init .

# Or use Makefile
make build

# Build with version info
make build VERSION=v1.0.0
```

### Testing
```bash
# Run all tests
go test ./...

# Verbose output
make test

# Generate coverage report (creates coverage.html)
make test-coverage

# Run specific package tests
go test -v ./pkg/checker
go test -v ./pkg/validator

# Run a single test
go test -v ./pkg/checker -run TestChecker_Check
```

### Code Quality
```bash
# Format and vet
make lint

# Or individually
go fmt ./...
go vet ./...
```

### Running the Tool
```bash
# Quick test of check command
make run-check

# Quick test of setup command
make run-setup

# Or run directly after building
./baseline-init check
./baseline-init setup --auto
./baseline-init validate SECURITY-INSIGHTS.yml
```

## Architecture

### Core Design Pattern

The tool follows a **layered architecture** with clear separation of concerns:

1. **CLI Layer** (`cmd/`)
   - Cobra commands handle user input and orchestrate operations
   - Each command (check, setup, validate) is independent
   - Commands use packages from `pkg/` but never depend on each other

2. **Business Logic Layer** (`pkg/`)
   - Self-contained packages with single responsibilities
   - No circular dependencies between packages
   - Each package can be used independently

3. **Data Flow**
   ```
   User Input → CLI Command → Package Logic → Output Formatter → Terminal
   ```

### Key Architectural Decisions

#### Multi-Location File Discovery
The checker (`pkg/checker/checker.go`) searches for compliance files in multiple locations:
- Repository root
- `.github/` directory
- `docs/` directory

This is implemented via the `check*()` helper methods (e.g., `checkSecurityInsights()`), which iterate through `possiblePaths` arrays. When adding new file checks, follow this pattern.

#### Schema Version Support
The validator (`pkg/validator/validator.go`) supports **both v1.0.0 and v2.0.0** of the Security Insights schema:

1. First unmarshals just the header to detect schema version
2. Routes to appropriate validator: `validateSecurityInsightsV1()` or `validateSecurityInsightsV2()`
3. v1 validation uses custom struct (`SecurityInsightsV1`)
4. v2 validation uses **official OpenSSF si-tooling structs** (`github.com/ossf/si-tooling/v2`)

**Current schema version generated**: 2.0.0 (as of latest update)

**Hybrid Approach for v2 Validation**:
- Uses `github.com/ossf/si-tooling/v2/si` for type-safe v2 validation
- Unmarshals YAML into official `si.SecurityInsights` struct
- Provides schema compliance via official OpenSSF type definitions
- Generation still uses template-based approach (see Template-Based Generation below)

When updating schema support:
- For v1: modify `SecurityInsightsV1` struct and `validateSecurityInsightsV1()`
- For v2: si-tooling structs auto-update with package upgrades
- Update the version detection logic in `validateSecurityInsights()`

#### Template-Based Generation
The generator (`pkg/generator/generator.go`) uses Go's `fmt.Sprintf()` with heredoc-style strings rather than separate template files. This keeps all generation logic in one place.

The `formatMaintainersV2()` helper demonstrates how to format complex nested YAML structures. When adding new generated files, follow this pattern.

#### Output Format Abstraction
The reporter (`pkg/report/formatter.go`) supports multiple output formats (text, JSON, YAML) via a strategy pattern:
- `OutputCheckResult()` dispatches to format-specific methods
- `outputText()` uses `github.com/fatih/color` for terminal colors
- JSON and YAML use standard library encoders

To add new output formats, add a new `output*()` method and update the switch statement.

## Package Responsibilities

### `pkg/checker`
- Scans repository for compliance files
- Returns structured `CheckResult` with file status and recommendations
- **Does not** validate file contents (that's `pkg/validator`'s job)
- Priority levels: critical, high, medium, low

### `pkg/generator`
- Creates SECURITY-INSIGHTS.yml (schema 2.0.0) and SECURITY.md
- Has two modes: auto (with defaults) and custom (with user config)
- Uses `Config` struct to pass parameters between packages
- Respects `force` flag to overwrite existing files

### `pkg/validator`
- Validates YAML syntax and schema compliance
- Auto-detects schema version for backward compatibility
- **v2 validation**: Uses official `github.com/ossf/si-tooling/v2` structs for type-safe schema validation
- **v1 validation**: Uses custom `SecurityInsightsV1` struct
- Returns `ValidationResult` with errors (fail) and warnings (pass but improve)
- Date validation: v1 uses RFC3339, v2 uses YYYY-MM-DD

### `pkg/interactive`
- Collects user input via `promptui` library
- Validates email format, detects Git remote URLs
- Returns `generator.Config` struct
- Single function: `GatherConfiguration()`

### `pkg/report`
- Formats compliance results for output
- Color-coded terminal output with priorities
- Groups recommendations by priority (critical → high → medium → low)

## Key Data Structures

### `checker.CheckResult`
The central data structure passed between checker → reporter:
```go
type CheckResult struct {
    Path            string             // Repository path
    IsCompliant     bool               // Overall compliance status
    Files           []FileCheck        // Status of each file
    MissingFiles    []string           // Quick list of missing files
    Recommendations []Recommendation   // Actionable next steps
}
```

### `generator.Config`
Configuration passed from interactive mode → generator:
```go
type Config struct {
    ProjectURL              string
    ProjectName             string
    SecurityEmail           string
    AcceptsVulnReports      bool
    AcceptsPullRequests     bool
    AcceptsAutomatedPR      bool
    ProjectStage            string  // active, archived, concept, moved, wip
    BugFixesOnly            bool
    Maintainers             []string // format: "github:username"
    DistributionPoints      []string
}
```

## Testing Philosophy

- **Unit tests** in `*_test.go` files alongside source
- Tests use temp directories (`os.MkdirTemp`) for file operations
- Each test is independent and cleans up after itself
- Test names follow pattern: `TestPackage_Method` or `TestPackage_Scenario`

When adding tests:
1. Test happy path and edge cases
2. Test multiple file locations for checker
3. Test both valid and invalid YAML for validator
4. Use table-driven tests for multiple scenarios

## Common Gotchas

1. **File paths must be absolute**: All file operations use absolute paths, not relative paths. Commands accept relative paths but convert them immediately.

2. **Schema version type**: In v2.0.0, `schema-version` can be a float (2.0) or string ("2.0.0"). The validator uses `interface{}` to handle both.

3. **Exit codes matter**: The check command exits with code 1 if non-compliant, enabling CI/CD integration. Don't change this behavior.

4. **Date formats differ by version**:
   - v1.0.0: RFC3339 (`2025-12-03T19:46:39-06:00`)
   - v2.0.0: YYYY-MM-DD (`2025-12-03`)

5. **Maintainer format**: Internal format is `github:username`, but v2.0.0 schema expands this to a full struct with name, email, affiliation, social, and primary fields.

## OpenSSF Compliance Context

The tool checks for these files (in priority order):
1. **SECURITY-INSIGHTS.yml** (High) - Machine-readable security metadata
2. **LICENSE** (High) - Open source license
3. **SECURITY.md** (Medium) - Human-readable security policy
4. **CODE_OF_CONDUCT.md** (Medium) - Community guidelines
5. **CONTRIBUTING.md** (Low) - Contribution guidelines

Reference specifications:
- [OpenSSF Security Baseline](https://github.com/ossf/security-baseline)
- [Security Insights Spec v2.0.0](https://github.com/ossf/security-insights-spec)

## Dependencies

### Core Dependencies
- `github.com/spf13/cobra` - CLI framework
- `github.com/fatih/color` - Terminal color output
- `github.com/manifoldco/promptui` - Interactive prompts
- `gopkg.in/yaml.v3` - YAML parsing

### OpenSSF Integration
- `github.com/ossf/si-tooling/v2` - Official OpenSSF Security Insights tooling
  - **Purpose**: Type-safe v2.0.0 schema validation
  - **Usage**: Validation only (not generation)
  - **Version**: v2.0.4
  - **Why**: Ensures schema compliance with official OpenSSF type definitions
  - **Note**: Generation uses templates, validation uses si-tooling (hybrid approach)

## Version Management

Version information is injected at build time via ldflags:
```bash
-ldflags "-X github.com/aguamala/baseline-init/cmd.Version=v1.0.0
          -X github.com/aguamala/baseline-init/cmd.GitCommit=$(git rev-parse HEAD)
          -X github.com/aguamala/baseline-init/cmd.BuildDate=$(date -u '+%Y-%m-%d_%H:%M:%S')"
```

These values are displayed via `baseline-init --version`.
