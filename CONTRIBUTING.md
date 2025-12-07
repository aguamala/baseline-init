# Contributing to baseline-init

Thank you for your interest in contributing to baseline-init! This document provides guidelines and instructions for contributing to the project.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Project Architecture](#project-architecture)
- [Making Changes](#making-changes)
- [Testing](#testing)
- [Submitting Changes](#submitting-changes)
- [Code Style](#code-style)
- [Common Gotchas](#common-gotchas)

## Code of Conduct

We expect all contributors to be respectful and professional. Please maintain a welcoming and inclusive environment for everyone.

## Getting Started

### Prerequisites

- **Go 1.23 or later** - Download from [golang.org](https://golang.org/dl/)
- **Git** - For version control
- **Make** (optional) - For using Makefile commands

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:
   ```bash
   git clone https://github.com/YOUR_USERNAME/baseline-init.git
   cd baseline-init
   ```
3. Add the upstream repository:
   ```bash
   git remote add upstream https://github.com/aguamala/baseline-init.git
   ```

## Development Setup

### Building the Project

```bash
# Standard build
go build -o baseline-init .

# Or use Makefile
make build

# Build with version info
make build VERSION=v1.0.0
```

### Running the Tool

```bash
# After building
./baseline-init check
./baseline-init setup --auto
./baseline-init validate SECURITY-INSIGHTS.yml

# Or use Makefile shortcuts
make run-check
make run-setup
```

### Dependencies

Install or update dependencies:

```bash
go mod download
go mod tidy
```

## Project Architecture

baseline-init follows a **layered architecture** with clear separation of concerns:

### Directory Structure

```
baseline-init/
├── cmd/              # CLI commands (check, setup, validate, version)
├── pkg/
│   ├── checker/      # Compliance file detection and checking
│   ├── generator/    # File generation (SECURITY-INSIGHTS.yml, SECURITY.md)
│   ├── validator/    # YAML schema validation (v1 and v2 support)
│   ├── interactive/  # User interaction and prompts
│   └── report/       # Output formatting (text, JSON, YAML)
├── main.go           # Entry point
└── go.mod            # Dependencies
```

### Key Design Patterns

1. **Multi-Location File Discovery**: The checker searches for compliance files in repository root, `.github/`, and `docs/` directories
2. **Schema Version Support**: The validator supports both v1.0.0 and v2.0.0 of Security Insights schema
3. **Template-Based Generation**: File generation uses Go's `fmt.Sprintf()` with heredoc-style strings
4. **Output Format Abstraction**: Multiple output formats (text, JSON, YAML) via strategy pattern

See [CLAUDE.md](CLAUDE.md) for detailed architectural documentation.

## Making Changes

### Creating a Branch

Create a descriptive branch name:

```bash
git checkout -b feature/add-new-validator
git checkout -b fix/checker-panic
git checkout -b docs/update-readme
```

### Development Workflow

1. **Make your changes** in the appropriate package
2. **Follow existing patterns** - look at similar code in the package
3. **Update tests** - add or modify tests for your changes
4. **Run tests** - ensure all tests pass
5. **Format code** - run `go fmt ./...`
6. **Vet code** - run `go vet ./...`

### Code Quality Commands

```bash
# Format and vet
make lint

# Or individually
go fmt ./...
go vet ./...
```

## Testing

### Running Tests

```bash
# Run all tests
go test ./...

# Verbose output
make test

# Run specific package tests
go test -v ./pkg/checker
go test -v ./pkg/validator

# Run a single test
go test -v ./pkg/checker -run TestChecker_Check
```

### Test Coverage

```bash
# Generate coverage report (creates coverage.html)
make test-coverage

# View coverage in terminal
go test -cover ./...
```

### Writing Tests

We follow these testing practices:

1. **Test files**: Place tests in `*_test.go` files alongside source code
2. **Test names**: Use pattern `TestPackage_Method` or `TestPackage_Scenario`
3. **Temp directories**: Use `os.MkdirTemp()` for file operations
4. **Independence**: Each test should be independent and clean up after itself
5. **Table-driven**: Use table-driven tests for multiple scenarios

**Example test structure:**

```go
func TestChecker_DetectFile(t *testing.T) {
    tests := []struct {
        name     string
        setup    func(dir string) error
        expected bool
    }{
        {
            name: "file in root",
            setup: func(dir string) error {
                return os.WriteFile(filepath.Join(dir, "SECURITY.md"), []byte("test"), 0644)
            },
            expected: true,
        },
        // More test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            dir := t.TempDir()
            if tt.setup != nil {
                if err := tt.setup(dir); err != nil {
                    t.Fatal(err)
                }
            }
            // Test logic here
        })
    }
}
```

### What to Test

- **Happy path**: Test expected behavior with valid inputs
- **Edge cases**: Test boundary conditions and unusual inputs
- **Error handling**: Test how code handles invalid inputs or errors
- **Multiple locations**: For checker, test files in root, `.github/`, and `docs/`
- **Schema versions**: For validator, test both v1 and v2 formats

## Submitting Changes

### Before Submitting

Ensure your changes:

- [ ] Pass all tests (`go test ./...`)
- [ ] Follow Go formatting (`go fmt ./...`)
- [ ] Pass go vet (`go vet ./...`)
- [ ] Include tests for new functionality
- [ ] Update documentation if needed
- [ ] Don't break existing functionality

### Pull Request Process

1. **Update your branch** with the latest upstream changes:
   ```bash
   git fetch upstream
   git rebase upstream/main
   ```

2. **Push your changes** to your fork:
   ```bash
   git push origin feature/your-feature-name
   ```

3. **Open a Pull Request** on GitHub with:
   - Clear title describing the change
   - Description of what changed and why
   - Reference to any related issues (e.g., "Fixes #123")
   - Screenshots or examples if applicable

4. **Respond to feedback** - maintainers may request changes

5. **Squash commits** if requested before merging

### Commit Messages

Write clear, descriptive commit messages:

```
Add support for custom schema validators

- Implement custom validator interface
- Add tests for validator plugin system
- Update documentation with examples
```

Format:
- First line: imperative mood, ~50 characters
- Blank line
- Detailed explanation if needed (wrap at 72 characters)

## Code Style

### Go Style Guidelines

Follow standard Go conventions:

- Use `gofmt` for formatting
- Follow [Effective Go](https://golang.org/doc/effective_go.html)
- Use meaningful variable names
- Add comments for exported functions and types
- Keep functions focused and small

### Package Guidelines

- **Single Responsibility**: Each package should have one clear purpose
- **No Circular Dependencies**: Packages should not depend on each other cyclically
- **Independent Packages**: Packages in `pkg/` should be usable independently
- **CLI Independence**: Commands in `cmd/` should not depend on each other

### Error Handling

- Return errors rather than panicking
- Wrap errors with context using `fmt.Errorf("context: %w", err)`
- Check all errors explicitly

### Documentation

- Add godoc comments for exported functions, types, and packages
- Update README.md if adding user-facing features
- Update CLAUDE.md if changing architecture or patterns

## Common Gotchas

Be aware of these project-specific details:

1. **File paths must be absolute** - All file operations use absolute paths; convert relative paths immediately
2. **Schema version flexibility** - In v2.0.0, `schema-version` accepts float or string
3. **Exit codes matter** - The check command exits with code 1 if non-compliant (enables CI/CD integration)
4. **Date format differences**:
   - v1.0.0: RFC3339 (`2025-12-03T19:46:39-06:00`)
   - v2.0.0: YYYY-MM-DD (`2025-12-03`)
5. **Maintainer format** - Internal format is `github:username`, v2.0.0 expands to full struct
6. **Overwrite protection** - Generator prompts before overwriting files unless `--force` is used

## Adding New Features

### Adding a New Compliance Check

1. Update `pkg/checker/checker.go`:
   - Add new `check*()` method following the existing pattern
   - Search multiple locations (root, `.github/`, `docs/`)
   - Add to `possiblePaths` array
2. Update `CheckResult` struct if needed
3. Add tests in `pkg/checker/checker_test.go`
4. Update documentation

### Adding a New Validator

1. Update `pkg/validator/validator.go`:
   - Add new validation function
   - Follow version detection pattern if applicable
2. Add tests with valid and invalid YAML examples
3. Update documentation

### Adding a New Output Format

1. Update `pkg/report/formatter.go`:
   - Add new `output*()` method
   - Update switch statement in `OutputCheckResult()`
2. Add format to command flags
3. Add tests
4. Update documentation

## Questions or Problems?

- **Bug Reports**: Open an [issue](https://github.com/aguamala/baseline-init/issues)
- **Feature Requests**: Open an [issue](https://github.com/aguamala/baseline-init/issues) with the "enhancement" label
- **Questions**: Start a [discussion](https://github.com/aguamala/baseline-init/discussions)

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.

---

Thank you for contributing to baseline-init!
