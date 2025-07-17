# Contributing to PostgreSQL Migration Tool

Thank you for your interest in contributing to the PostgreSQL Migration Tool! We welcome contributions from the community and are pleased to have you join us.

## 🌟 Show Your Support

If you find this project helpful, please consider:

-   ⭐ **Starring the repository** - it helps others discover the project
-   🍴 **Forking the project** if you plan to contribute
-   📢 **Sharing the project** with others who might benefit

## 📋 Table of Contents

-   [Code of Conduct](#code-of-conduct)
-   [Getting Started](#getting-started)
-   [Development Setup](#development-setup)
-   [Project Structure](#project-structure)
-   [Coding Standards](#coding-standards)
-   [Making Changes](#making-changes)
-   [Pull Request Process](#pull-request-process)
-   [Issue Reporting](#issue-reporting)
-   [Documentation](#documentation)

## 📜 Code of Conduct

This project adheres to a code of conduct that we expect all contributors to follow:

-   **Be respectful** and inclusive in all interactions
-   **Be constructive** when providing feedback
-   **Be patient** with newcomers and questions
-   **Focus on what's best** for the community and project

## 🚀 Getting Started

### Prerequisites

Before contributing, ensure you have:

-   **Go 1.19+** installed
-   **Git** for version control
-   **Docker** (optional, for testing Docker functionality)
-   **PostgreSQL client tools** (optional, for testing host functionality)

### Development Setup

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:

    ```bash
    git clone https://github.com/YOUR_USERNAME/pg-migrator.git
    cd pg-migrator
    ```

3. **Add the upstream remote**:

    ```bash
    git remote add upstream https://github.com/ayussh-2/pg-migrator.git
    ```

4. **Install dependencies**:

    ```bash
    go mod download
    go mod tidy
    ```

5. **Verify the setup**:
    ```bash
    go build
    ```

## 🏗️ Project Structure

Understanding the codebase structure is crucial for effective contributions:

```
pg-migrator/
├── main.go                 # Application entry point and CLI setup
├── utils/                  # Utility packages
│   └── env.go             # Environment configuration and validation
├── migrator/              # Core migration logic
│   ├── migrator.go        # Main migration orchestrator and workflow
│   ├── tools.go           # PostgreSQL tool detection and validation
│   ├── connection.go      # Database connectivity testing
│   ├── dump.go            # Database dump creation and restoration
│   └── verify.go          # Migration verification and reporting
├── .env.sample           # Sample environment configuration
├── README.md             # Project documentation
├── CONTRIBUTING.md       # This file
├── LICENSE               # MIT License
└── go.mod               # Go module dependencies
```

### Module Responsibilities

-   **`main.go`**: Entry point, argument parsing, error handling
-   **`utils/env.go`**: Configuration loading, validation, defaults
-   **`migrator/migrator.go`**: Orchestrates the migration workflow
-   **`migrator/tools.go`**: Checks for required PostgreSQL tools
-   **`migrator/connection.go`**: Tests database connections
-   **`migrator/dump.go`**: Handles pg_dump and psql operations
-   **`migrator/verify.go`**: Verifies migration success and generates reports

## 📝 Coding Standards

### Go Code Style

Follow these conventions to maintain code consistency:

#### 1. **Naming Conventions**

```go
// ✅ Good: Clear, descriptive names
func createDatabaseDump() error
type MigrationConfig struct
var hostConnectionString string

// ❌ Bad: Unclear abbreviations
func createDbDmp() error
type MigConfig struct
var hostConnStr string
```

#### 2. **Function Comments**

All exported functions and types must have comments:

```go
// ✅ Good: Clear documentation
// createDump generates a SQL dump of the source database using pg_dump.
// It creates a timestamped backup directory and handles both Docker and host-based execution.
func (m *Migrator) createDump() error {
    // Implementation...
}

// ❌ Bad: Missing or unclear comments
func (m *Migrator) createDump() error {
    // Creates dump
}
```

#### 3. **Error Handling**

Always provide context in error messages:

```go
// ✅ Good: Descriptive error context
if err != nil {
    return fmt.Errorf("failed to create backup directory %s: %v", m.backupDir, err)
}

// ❌ Bad: Generic error handling
if err != nil {
    return err
}
```

#### 4. **Code Organization**

-   Keep functions focused and single-purpose
-   Use early returns to reduce nesting
-   Group related functionality together

```go
// ✅ Good: Early return, clear structure
func (m *Migrator) validateConfig() error {
    if m.config.HostHost == "" {
        return fmt.Errorf("HOST_HOST is required")
    }

    if m.config.TargetHost == "" {
        return fmt.Errorf("TARGET_HOST is required")
    }

    return nil
}
```

#### 5. **Constants and Variables**

```go
// ✅ Good: Use constants for magic numbers/strings
const (
    DefaultPostgreSQLPort = "5432"
    DefaultSSLMode       = "disable"
    BackupDirPermissions = 0755
)

// ✅ Good: Descriptive variable names
var (
    migrationStartTime = time.Now()
    backupDirectory    = fmt.Sprintf("migration_%s", timestamp)
)
```

### File Structure Standards

#### 1. **Import Organization**

```go
// Standard library imports
import (
    "fmt"
    "os"
    "path/filepath"
)

// Third-party imports
import (
    "github.com/joho/godotenv"
    "github.com/lib/pq"
)

// Local imports
import (
    "pg-migrator/utils"
)
```

#### 2. **Function Order**

-   Public functions first
-   Private functions after
-   Helper functions at the end

## 🔄 Making Changes

### 1. **Create a Feature Branch**

```bash
git checkout -b feature/your-feature-name
# or
git checkout -b fix/issue-description
```

### 2. **Make Your Changes**

-   Follow the coding standards above
-   Update documentation as needed
-   Test your changes thoroughly with real database scenarios

### 3. **Manual Testing**

Since automated tests are not yet implemented, please thoroughly test your changes:

```bash
# Test different scenarios manually
cp .env.sample .env.test
# Edit .env.test with test credentials
go run main.go

# Test various configurations:
# - Docker to Docker migrations
# - Host to Docker migrations
# - Different SSL modes
# - Error scenarios
```

### 4. **Commit Messages**

Use clear, descriptive commit messages:

```bash
# ✅ Good commit messages
git commit -m "feat: add support for custom SSL modes in configuration"
git commit -m "fix: resolve Docker container password passing issue"
git commit -m "docs: update README with new configuration examples"
git commit -m "refactor: simplify database connection validation logic"

# ❌ Bad commit messages
git commit -m "fix stuff"
git commit -m "update"
git commit -m "working"
```

#### Commit Message Format

```
type(scope): description

[optional body]

[optional footer]
```

**Types:**

-   `feat`: New feature
-   `fix`: Bug fix
-   `docs`: Documentation changes
-   `refactor`: Code refactoring
-   `test`: Adding tests
-   `chore`: Maintenance tasks

## 🔀 Pull Request Process

### 1. **Before Submitting**

-   [ ] Code follows the style guidelines
-   [ ] Self-review of your code completed
-   [ ] Documentation updated if needed
-   [ ] Manual testing completed across different scenarios

### 2. **Pull Request Template**

When creating a PR, include:

```markdown
## Description

Brief description of changes made.

## Type of Change

-   [ ] Bug fix (non-breaking change which fixes an issue)
-   [ ] New feature (non-breaking change which adds functionality)
-   [ ] Breaking change (fix or feature that would cause existing functionality to not work as expected)
-   [ ] Documentation update

## Testing

-   [ ] Manual testing completed
-   [ ] Tested with Docker configuration
-   [ ] Tested with host-only configuration
-   [ ] Tested with different database scenarios

## Checklist

-   [ ] Code follows project style guidelines
-   [ ] Self-review completed
-   [ ] Documentation updated
```

### 3. **Review Process**

-   PRs require at least one approval
-   Address all review comments
-   Keep discussions focused and constructive
-   Be responsive to feedback

## 🐛 Issue Reporting

### Bug Reports

When reporting bugs, include:

````markdown
**Bug Description**
Clear description of the bug.

**Environment**

-   OS: [e.g., Windows 10, macOS Big Sur, Ubuntu 20.04]
-   Go version: [e.g., 1.19.5]
-   Docker version: [e.g., 20.10.17] (if applicable)

**Configuration**

```env
# Sanitized .env configuration (remove sensitive data)
HOST_HOST=example.com
TARGET_DOCKER_CONTAINER=postgres-container
```
````

**Steps to Reproduce**

1. Step one
2. Step two
3. Error occurs

**Expected vs Actual Behavior**

-   Expected: What should happen
-   Actual: What actually happened

**Error Messages**

```
Paste any error messages here
```

**Additional Context**
Any other relevant information.

````

### Feature Requests
```markdown
**Feature Description**
Clear description of the proposed feature.

**Use Case**
Why would this feature be useful?

**Proposed Implementation**
Any ideas on how this could be implemented?

**Alternatives Considered**
Any alternative solutions you've considered?
````

## 📚 Documentation

### Code Documentation

-   Comment all exported functions and types
-   Use clear, descriptive comments
-   Include usage examples where helpful

### README Updates

When making changes that affect usage:

-   Update relevant sections in README.md
-   Add new configuration examples if needed
-   Update troubleshooting section for new issues

### Documentation Style

-   Use clear, concise language
-   Include practical examples
-   Keep documentation up-to-date with code changes

## 🎯 Areas for Contribution

We especially welcome contributions in these areas:

### 🔧 **Technical Improvements**

-   Performance optimizations
-   Better error handling and recovery
-   Enhanced Docker support
-   Cross-platform compatibility improvements

### 📖 **Documentation**

-   Additional configuration examples
-   Troubleshooting guides
-   Video tutorials or guides
-   Translation to other languages

### 🧪 **Testing**

-   Unit tests for existing functionality
-   Integration tests
-   Docker-specific test scenarios
-   Performance benchmarks

### 🌟 **Features**

-   Support for additional database types
-   GUI interface
-   Configuration validation improvements
-   Migration rollback capabilities
-   Progress bars and better user feedback

## 💡 Getting Help

If you need help while contributing:

1. **Check existing documentation** first
2. **Search closed issues** for similar problems
3. **Create a discussion** for general questions
4. **Join our community** discussions

## 🙏 Recognition

Contributors will be recognized in:

-   README.md contributors section
-   Release notes for significant contributions
-   Special thanks for major features or fixes

---

Thank you for contributing to the PostgreSQL Migration Tool! Your efforts help make database migrations easier for everyone. 🚀

**Don't forget to ⭐ star the repository if you find it useful!**
