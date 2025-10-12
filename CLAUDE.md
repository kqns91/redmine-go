# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is an unofficial Redmine API client SDK for Go. The goal is to provide comprehensive coverage of all Redmine REST API endpoints as documented at https://www.redmine.org/projects/redmine/wiki/Rest_api.

## Architecture

### Core Structure

The codebase follows a flat, service-oriented architecture under `./pkg/redmine/`:

- **client.go**: Contains the `Client` struct and `do()` method for making HTTP requests to Redmine API
  - All endpoint-specific methods should wrap `do()` for HTTP communication
  - Handles authentication via X-Redmine-API-Key header
  - Manages JSON content type and error responses

- **type.go**: Defines shared models used across multiple services
  - Contains common types like `CustomField`, `Resource`, etc.
  - Only types referenced by multiple service files belong here

- **Service files** (e.g., project.go): One file per Redmine service (projects, issues, users, etc.)
  - Each file contains service-specific models, request/response structs, and API endpoint methods
  - Methods are implemented as methods on the `Client` struct

### Development Workflow

When implementing a new service:

1. Check the Redmine API reference at https://www.redmine.org/projects/redmine/wiki/Rest_api
2. Create a todo list breaking down the implementation into steps
3. Create a new file in `./pkg/redmine/` named after the service (e.g., `issues.go`, `users.go`)
4. Define service-specific models, request structs, and response structs
5. Move any models used by multiple services to `type.go`
6. Implement endpoint methods that wrap the `Client.do()` method

## Commands

### Linting
```bash
golangci-lint run
```
All code must pass linting before committing. The configuration uses `default: all` linters with custom formatters including gci, gofmt, gofumpt, goimports, and swaggo.

### Testing
```bash
go test ./...
```

### Build
```bash
go build ./...
```

## Development Guidelines

- Implement services one at a time with appropriate granularity
- Always verify the API reference before implementing endpoints
- Create clear, properly-sized commits for each logical change
- Run `golangci-lint run` before every commit to ensure code quality
- Follow Go idiomatic patterns and naming conventions
- All code should use the shared `Client.do()` method for HTTP requests
