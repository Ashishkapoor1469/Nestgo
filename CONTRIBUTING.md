# Contributing to NestGo

Thank you for your interest in contributing to NestGo! We welcome contributions from everyone.

## How to Contribute

### 1. Report Bugs
If you find a bug, please open an issue using the Bug Report template. Include as much detail as possible, such as steps to reproduce and environment information.

### 2. Suggest Features
Have an idea for a new feature? Open an issue using the Feature Request template.

### 3. Submit Pull Requests
1. Fork the repository.
2. Create a new branch for your changes: `git checkout -b feature/my-new-feature`.
3. Make your changes and ensure tests pass.
4. Commit your changes: `git commit -m "feat: add my new feature"`.
5. Push to your branch: `git push origin feature/my-new-feature`.
6. Open a Pull Request against the `main` branch.

## Coding Standards
- Follow idiomatic Go patterns.
- Ensure `go fmt` is run on all files.
- Add tests for new features.
- Update documentation if applicable.

## Development Setup
1. Clone the repo.
2. Run `go mod tidy`.
3. Build the CLI: `go build -o nestgo ./cmd/nestgo`.

Thank you for making NestGo better!
