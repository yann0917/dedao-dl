# Contributing to dedao-dl

Thank you for your interest in contributing to dedao-dl! We welcome contributions from the community.

## How to Contribute

### Reporting Bugs
- Check if the bug has already been reported in the [Issues](https://github.com/yann0917/dedao-dl/issues)
- If not, create a new issue with:
  - Clear description of the bug
  - Steps to reproduce
  - Expected and actual behavior
  - Your environment (OS, Go version, etc.)

### Suggesting Enhancements
- Use the Issues tab to suggest new features
- Provide a clear description of the proposed feature
- Explain the use case and benefits

### Submitting Pull Requests
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin feature/your-feature`)
5. Open a Pull Request with a clear description of your changes

## Development Setup

1. Clone the repository
2. Ensure you have Go (>=1.18) installed
3. Install dependencies using Go Modules (`go mod tidy`)
4. Format your code with `gofmt` before committing (`gofmt -w .`)
5. Run tests with `go test ./...` 
6. Ensure your changes do not break existing tests

## Code Style

- Format all code with gofmt (`gofmt -w .`)
- Follow idiomatic Go code style conventions
- Write clear, descriptive commit messages
- Add comments for complex logic

## License

By contributing, you agree that your contributions will be licensed under the same license as the project.