# Contributing to Glint

Thank you for your interest in contributing to Glint! This document provides guidelines and instructions for contributing.

## Code of Conduct

By participating in this project, you agree to abide by our [Code of Conduct](CODE_OF_CONDUCT.md).

## Getting Started

1. **Fork the repository** on GitHub
2. **Clone your fork** locally:
   ```bash
   git clone https://github.com/YOUR-USERNAME/glint.git
   cd glint
   ```
3. **Add the upstream repository**:
   ```bash
   git remote add upstream https://github.com/droqsic/glint.git
   ```
4. **Create a branch** for your work:
   ```bash
   git checkout -b feature/your-feature-name
   ```

## Development Workflow

1. **Make your changes** in your feature branch
2. **Write or update tests** as needed
3. **Run tests** to ensure they pass:
   ```bash
   go test ./...
   ```
4. **Run benchmarks** if you're making performance-related changes:
   ```bash
   go test -bench=. ./...
   ```
5. **Format your code**:
   ```bash
   go fmt ./...
   ```
6. **Verify with go vet**:
   ```bash
   go vet ./...
   ```
7. **Commit your changes** with a clear commit message:
   ```bash
   git commit -m "Add feature: your feature description"
   ```
8. **Push to your fork**:
   ```bash
   git push origin feature/your-feature-name
   ```
9. **Create a Pull Request** from your fork to the main repository

## Pull Request Guidelines

When submitting a pull request:

1. **Include tests** for any new functionality
2. **Update documentation** as needed
3. **Follow Go coding conventions**
4. **Keep PRs focused** - submit separate PRs for unrelated changes
5. **Be responsive to feedback** - be willing to update your PR based on reviews

## Reporting Issues

When reporting issues:

1. **Use the issue template** provided
2. **Include reproduction steps** - how can we reproduce the issue?
3. **Include environment details** - Go version, OS, etc.
4. **Include logs or error messages** if applicable

## Contributing to Glint Features

Glint is focused on terminal color support detection and enabling. When contributing, keep these core features in mind:

### Color Support Detection

If you're improving color support detection:

1. Consider different terminal types and environments
2. Respect environment variables like `TERM`, `COLORTERM`, and `NO_COLOR`
3. Add tests for the new detection logic
4. Document the detection behavior

### Color Support Levels

If you're working with color support levels:

1. Ensure compatibility with the existing level definitions (`LevelNone`, `Level16`, `Level256`, `LevelTrue`)
2. Add tests for level detection in different environments
3. Document how the levels are determined

### Force Color Support

If you're enhancing the force color support functionality:

1. Ensure it works correctly across different platforms
2. Consider edge cases (redirected output, non-terminal environments)
3. Add tests for the forcing behavior
4. Document any limitations or platform-specific behavior

## License

By contributing to Glint, you agree that your contributions will be licensed under the project's [MIT License](LICENSE).
