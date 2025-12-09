# Contributing to baseline-init

Thank you for your interest in contributing to `baseline-init`! We welcome contributions from the community to help improve OpenSSF baseline compliance tooling.

## Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Getting Started](#getting-started)
- [Development Setup](#development-setup)
- [Development Workflow](#development-workflow)
- [Testing Guidelines](#testing-guidelines)
- [Code Style](#code-style)
- [Submitting Changes](#submitting-changes)
- [Project Architecture](#project-architecture)
- [Getting Help](#getting-help)

## Code of Conduct

This project adheres to the Contributor Covenant [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report unacceptable behavior to the project maintainers.

## Getting Started

### Prerequisites

- **Go 1.23 or later** - [Download Go](https://golang.org/dl/)
- **Git** - [Install Git](https://git-scm.com/downloads)
- Familiarity with Go development and testing
- Understanding of OpenSSF Security Baseline concepts (helpful but not required)

### Fork and Clone

1. Fork the repository on GitHub
2. Clone your fork locally:

```bash
git clone https://github.com/YOUR-USERNAME/baseline-init.git
cd baseline-init
```

3. Add the upstream repository:

```bash
git remote add upstream https://github.com/aguamala/baseline-init.git
```

## Development Setup

### Build the Project

```bash
# Standard build
go build -o baseline-init .

# Or use the Makefile
make build
```

### Run the Tool Locally

```bash
# Check command
./baseline-init check

# Setup command
./baseline-init setup --auto

# Validate command
./baseline-init validate SECURITY-INSIGHTS.yml
```

### Install Dependencies

Dependencies are managed via Go modules:

```bash
go mod download
go mod tidy
```

## Development Workflow

### 1. Create a Feature Branch

```bash
git checkout -b feature/your-feature-name
```

Branch naming conventions:
- `feature/` - New features
- `fix/` - Bug fixes
- `docs/` - Documentation changes
- `refactor/` - Code refactoring
- `test/` - Test additions or fixes

### 2. Make Your Changes

- Write clear, concise code
- Follow the existing code structure and patterns
- Add tests for new functionality
- Update documentation as needed

### 3. Run Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
make test

# Run tests with coverage
make test-coverage
```

The coverage report will be generated in `coverage.html`.

### 4. Run Code Quality Checks

```bash
# Format code
go fmt ./...

# Vet code
go vet ./...

# Or use the Makefile
make lint
```

### 5. Commit Your Changes

Write clear, descriptive commit messages:

```bash
git commit -m "Add feature: description of what you added"
git commit -m "Fix: description of what you fixed"
git commit -m "Docs: description of documentation changes"
```

Commit message guidelines:
- Use the imperative mood ("Add feature" not "Added feature")
- Keep the first line under 72 characters
- Reference issues and pull requests when applicable
- Provide context in the commit body for complex changes

### 6. Keep Your Branch Updated

```bash
git fetch upstream
git rebase upstream/main
```

### 7. Push to Your Fork

```bash
git push origin feature/your-feature-name
```

## Testing Guidelines

### Test Structure

- Tests are located in `*_test.go` files alongside source code
- Each package should have comprehensive test coverage
- Use table-driven tests for multiple scenarios

### Test Naming

Follow the pattern: `TestPackage_Method` or `TestPackage_Scenario`

```go
func TestChecker_Check(t *testing.T) { ... }
func TestValidator_ValidateSecurityInsights(t *testing.T) { ... }
```

### Test Best Practices

1. **Use temp directories** for file operations:
```go
tmpDir := t.TempDir()
```

2. **Test both success and failure cases**:
```go
tests := []struct {
    name    string
    input   string
    want    bool
    wantErr bool
}{
    {"valid input", "test", true, false},
    {"invalid input", "", false, true},
}
```

3. **Clean up resources**:
- Use `t.Cleanup()` or `defer` statements
- Remove temporary files and directories

4. **Test file locations**:
- For the checker, test files in root, `.github/`, and `docs/`
- Ensure path handling works across platforms

### Running Specific Tests

```bash
# Run tests for a specific package
go test -v ./pkg/checker

# Run a single test
go test -v ./pkg/checker -run TestChecker_Check

# Run tests with race detection
go test -race ./...
```

## Code Style

### General Guidelines

1. **Follow Go conventions**:
   - Use `gofmt` and `go vet`
   - Follow [Effective Go](https://golang.org/doc/effective_go)
   - Use meaningful variable and function names

2. **Package organization**:
   - Each package has a single, well-defined responsibility
   - No circular dependencies between packages
   - Packages should be independently usable

3. **Error handling**:
   - Always check and handle errors
   - Provide context with error wrapping: `fmt.Errorf("context: %w", err)`
   - Return errors rather than panicking

4. **Comments**:
   - Document all exported functions, types, and constants
   - Use `//` for inline comments
   - Keep comments concise and up-to-date

### Architecture Patterns

Follow these patterns when adding features:

#### Multi-Location File Discovery

When checking for compliance files, search multiple locations:

```go
possiblePaths := []string{
    filepath.Join(repoPath, "FILENAME"),
    filepath.Join(repoPath, ".github", "FILENAME"),
    filepath.Join(repoPath, "docs", "FILENAME"),
}
```

#### Template-Based Generation

Use `fmt.Sprintf()` with heredoc-style strings for file generation:

```go
content := fmt.Sprintf(`---
key: %s
value: %s
---`, key, value)
```

#### Output Format Abstraction

Support multiple output formats (text, JSON, YAML) via a dispatch pattern in the reporter.

## Submitting Changes

### Pull Request Process

1. **Ensure all tests pass**:
```bash
go test ./...
make lint
```

2. **Update documentation**:
   - Update README.md if adding new features
   - Update CLAUDE.md if changing architecture
   - Add inline code comments

3. **Create a Pull Request**:
   - Use a clear, descriptive title
   - Reference related issues
   - Describe what changed and why
   - Include testing evidence

### Pull Request Template

```markdown
## Description
Brief description of changes

## Related Issues
Fixes #123

## Changes Made
- Added feature X
- Fixed bug Y
- Updated documentation

## Testing
- [ ] All tests pass
- [ ] Added tests for new functionality
- [ ] Manual testing completed

## Checklist
- [ ] Code follows project style guidelines
- [ ] Documentation updated
- [ ] No breaking changes (or documented)
```

### Review Process

- Maintainers will review your PR
- Address feedback and update your branch
- Once approved, a maintainer will merge your PR

## Project Architecture

### Core Packages

#### `pkg/checker`
- Scans repositories for compliance files
- Returns structured `CheckResult` with file status
- Does not validate file contents

#### `pkg/generator`
- Creates SECURITY-INSIGHTS.yml (schema 2.0.0) and SECURITY.md
- Supports auto and interactive modes
- Includes overwrite protection

#### `pkg/validator`
- Validates YAML syntax and schema compliance
- Auto-detects schema version (v1.0.0 vs v2.0.0)
- Uses official OpenSSF si-tooling for v2 validation

#### `pkg/interactive`
- Collects user input via `promptui`
- Validates email format and Git URLs
- Returns `generator.Config` struct

#### `pkg/report`
- Formats compliance results for output
- Supports text, JSON, and YAML formats
- Color-coded terminal output

### Data Flow

```
User Input → CLI Command → Package Logic → Output Formatter → Terminal
```

### Key Design Principles

1. **Layered architecture** with clear separation of concerns
2. **No circular dependencies** between packages
3. **Self-contained packages** that can be used independently
4. **CLI commands** orchestrate but don't contain business logic

## Getting Help

### Resources

- **Documentation**: [README.md](README.md)
- **Architecture Guide**: [CLAUDE.md](CLAUDE.md)
- **Issues**: [GitHub Issues](https://github.com/aguamala/baseline-init/issues)
- **Discussions**: [GitHub Discussions](https://github.com/aguamala/baseline-init/discussions)

### Asking Questions

- Check existing issues and discussions first
- Provide context and examples when asking questions
- Be respectful and patient

### Reporting Bugs

When reporting bugs, include:
- Go version (`go version`)
- Operating system and version
- Steps to reproduce
- Expected vs actual behavior
- Relevant logs or error messages

### Suggesting Features

When suggesting features:
- Describe the use case
- Explain why it's valuable
- Consider backward compatibility
- Propose implementation approach (optional)

## Recognition

Contributors are recognized in:
- Git commit history
- GitHub contributor graph
- Release notes for significant contributions

Thank you for contributing to `baseline-init`!
