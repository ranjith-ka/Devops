# Changelog

All notable changes to this project will be documented in this file.

## [1.2.0] - 2023-10-XX

### Added [1.2.0]

- Introduced `/config` endpoint to fetch current application configuration.
- Added support for dynamic reloading of configuration without restarting the application.
- Implemented middleware for request logging and error handling.
- Added integration tests for `serve.go` endpoints.
- Added support for environment variable `ENABLE_METRICS` to toggle Prometheus metrics.
- Added support for installing the `gh-copilot-cli` GitHub CLI extension.

### Updated [1.2.0]

- Refactored configuration management to use `viper` for better flexibility.
- Enhanced logging to include request IDs for traceability.
- Updated Dockerfile to use multi-stage builds for smaller image size.
- Improved CI/CD pipeline to include integration tests and code coverage reports.
- Updated `Readme.md` with instructions for enabling/disabling metrics.

### Fixed [1.1.0]

- Fixed incorrect content type in `/metrics` endpoint response.
- Resolved issue with configuration file not being detected in certain environments.
- Fixed flaky tests in `random_test.go`.

## [1.1.0] - 2023-10-XX

### Added [1.1.0]

- Implemented `status` endpoint to check application health.
- Added support for structured logging using `logrus`.
- Introduced configuration file support for managing application settings.
- Added Dockerfile for containerizing the application.
- Added CI/CD pipeline configuration using GitHub Actions.
- Added unit tests for `random.go` and `serve.go`.
- Added support for environment variable `APP_CONFIG_PATH` for custom configuration file paths.

### Updated [1.0.0]

- Refactored `serve.go` to improve modularity and readability.
- Enhanced Prometheus metrics to include additional application metrics.
- Updated `Readme.md` with Docker usage instructions and CI/CD pipeline details.
- Improved error messages for better debugging.

### Fixed [1.0.0]

- Resolved issue with incorrect HTTP status codes in `/joke` endpoint.
- Fixed race condition in concurrent API requests.
- Addressed memory leak in `getCopilotResponse` function.

## [1.0.0] - YYYY-MM-DD

### Added

- Initial implementation of the `serve` command with endpoints:
  - `/`: Welcome message.
  - `/hello`: Displays the first HTTP program message.
  - `/hello2`: Displays the second HTTP program message.
  - `/headers`: Displays the request headers.
  - `/joke`: Fetches a random joke.
- Integration with Prometheus metrics at `/metrics`.
- Support for fetching random jokes from `https://icanhazdadjoke.com/`.
- Added Copilot Chat integration for generating programming jokes.
- Documentation for application usage in `doc/usage.md`.
- Instructions for using custom instructions in Copilot Chat in `Readme.md`.
- Added `random` command to fetch a random joke from the terminal.
- Added Copilot API integration to fetch programming jokes using prompts.
- Added support for environment variable `GITHUB_TOKEN` for Copilot API authentication.
- Added error handling for HTTP requests and API responses.

### Updated

- Enhanced error handling for HTTP responses and API requests.

### Fixed

- Minor formatting issues in the documentation.
- Fixed potential nil pointer dereference in HTTP response handling.
- Improved error handling for API requests and responses.
