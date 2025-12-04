# baseline-init

A command-line tool for OpenSSF Baseline compliance checking and setup.

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/aguamala/baseline-init)](https://goreportcard.com/report/github.com/aguamala/baseline-init)

## Overview

`baseline-init` helps repositories achieve and maintain [OpenSSF Security Baseline](https://github.com/ossf/security-baseline) compliance by:

- **Checking** repositories for missing compliance requirements
- **Validating** existing compliance files against schemas
- **Auto-generating** compliant default files (SECURITY-INSIGHTS.yml, SECURITY.md, etc.)
- **Guiding** users through interactive setup

## Installation

### From Source

```bash
git clone https://github.com/aguamala/baseline-init.git
cd baseline-init
go build -o baseline-init
```

### Using Go Install

```bash
go install github.com/aguamala/baseline-init@latest
```

## Quick Start

### Check Compliance

Check if your repository meets OpenSSF baseline requirements:

```bash
baseline-init check
```

Output formats:
```bash
baseline-init check --format json   # JSON output
baseline-init check --format yaml   # YAML output
baseline-init check --format text   # Human-readable (default)
```

### Setup Compliance Files

#### Auto Mode (Quick Start)

Generate files with sensible defaults:

```bash
baseline-init setup --auto
```

#### Interactive Mode (Recommended)

Walk through guided setup with prompts:

```bash
baseline-init setup --interactive
```

The interactive mode will ask you for:
- Project URL
- Security contact email
- Project lifecycle stage
- Vulnerability reporting preferences
- Pull request policies
- Maintainer information

#### Force Overwrite

Overwrite existing files:

```bash
baseline-init setup --auto --force
```

### Validate Files

Validate compliance files against their schemas:

```bash
baseline-init validate SECURITY-INSIGHTS.yml
baseline-init validate .github/SECURITY-INSIGHTS.yml
```

## Commands

### `baseline-init check [path]`

Scan a repository for OpenSSF baseline compliance.

**Flags:**
- `-f, --format` - Output format: text, json, yaml (default: text)
- `-p, --path` - Path to repository (default: current directory)

**Example:**
```bash
baseline-init check /path/to/repo --format json
```

**Exit Codes:**
- `0` - Repository is compliant
- `1` - Repository is not compliant or error occurred

### `baseline-init setup [path]`

Generate OpenSSF baseline compliance files.

**Flags:**
- `--auto` - Auto-generate with defaults
- `--interactive` - Interactive setup mode
- `--force` - Overwrite existing files
- `-p, --path` - Path to repository (default: current directory)

**Example:**
```bash
baseline-init setup --interactive
baseline-init setup /path/to/repo --auto
```

### `baseline-init validate <file>`

Validate a compliance file against its schema.

**Example:**
```bash
baseline-init validate SECURITY-INSIGHTS.yml
```

### `baseline-init version`

Display version information.

## Generated Files

### SECURITY-INSIGHTS.yml

The primary compliance file containing security metadata following the [OpenSSF Security Insights specification](https://github.com/ossf/security-insights-spec).

**Location:** Repository root or `.github/SECURITY-INSIGHTS.yml`

**Key sections:**
- `header` - Metadata and versioning
- `project-lifecycle` - Project status and maintenance
- `contribution-policy` - PR and contribution policies
- `security-contacts` - Security team contact information
- `vulnerability-reporting` - Vulnerability disclosure policies
- `security-testing` - Security testing practices
- `dependencies` - Dependency management

### SECURITY.md

Security policy and vulnerability reporting instructions.

**Location:** Repository root, `.github/SECURITY.md`, or `docs/SECURITY.md`

## Compliance Requirements

The tool checks for the following OpenSSF baseline requirements:

### Required Files

- ✅ **SECURITY-INSIGHTS.yml** - Security metadata (High Priority)
- ✅ **LICENSE** - Open source license (High Priority)
- ✅ **SECURITY.md** - Security policy (Medium Priority)

### Recommended Files

- 📋 **CODE_OF_CONDUCT.md** - Code of conduct (Medium Priority)
- 📋 **CONTRIBUTING.md** - Contribution guidelines (Low Priority)

## CI/CD Integration

### GitHub Actions

Add compliance checking to your workflow:

```yaml
name: OpenSSF Baseline Check

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  compliance:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Install baseline-init
        run: go install github.com/aguamala/baseline-init@latest

      - name: Check compliance
        run: baseline-init check --format text
```

## Development

### Prerequisites

- Go 1.23 or later
- Git

### Building

```bash
go build -o baseline-init
```

### Running Tests

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

### Project Structure

```
baseline-init/
├── cmd/              # Command definitions (check, setup, validate)
├── pkg/
│   ├── checker/      # Compliance checking logic
│   ├── generator/    # File generation logic
│   ├── validator/    # YAML validation logic
│   ├── interactive/  # Interactive prompts
│   └── report/       # Output formatting
├── main.go          # Entry point
├── go.mod           # Go module definition
└── README.md        # This file
```

## Contributing

Contributions are welcome! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./...`
5. Submit a pull request

## License

Apache License 2.0 - see [LICENSE](LICENSE) for details.

## References

- [OpenSSF Security Baseline](https://github.com/ossf/security-baseline)
- [OpenSSF Security Insights Specification](https://github.com/ossf/security-insights-spec)
- [OpenSSF Best Practices Badge](https://bestpractices.coreinfrastructure.org/)

## Support

- **Issues**: [GitHub Issues](https://github.com/aguamala/baseline-init/issues)
- **Discussions**: [GitHub Discussions](https://github.com/aguamala/baseline-init/discussions)

## Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [promptui](https://github.com/manifoldco/promptui) - Interactive prompts
- [color](https://github.com/fatih/color) - Terminal colors
- [yaml.v3](https://gopkg.in/yaml.v3) - YAML parsing

Based on the [OpenSSF Security Baseline](https://github.com/ossf/security-baseline) and [Security Insights](https://github.com/ossf/security-insights-spec) projects.
