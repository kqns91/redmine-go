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
7. Create corresponding test file (e.g., `issues_test.go`) with comprehensive unit tests

### Testing Strategy

All API implementations should have corresponding unit tests:

- **Test files**: Name test files as `<service>_test.go` (e.g., `project_test.go`, `issue_test.go`)
- **HTTP mocking**: Use `httptest.NewServer` for reliable HTTP request/response mocking
- **Test coverage**: Aim for comprehensive coverage of all CRUD operations
- **Test structure**: Each test should validate:
  - HTTP method (GET, POST, PUT, DELETE)
  - Request path and query parameters
  - Request headers (X-Redmine-API-Key, Content-Type)
  - Request body structure (for POST/PUT)
  - Response body deserialization
  - Error handling

Example test pattern:
```go
func TestListProjects(t *testing.T) {
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Validate request
        if r.Method != http.MethodGet {
            t.Errorf("Expected GET request, got %s", r.Method)
        }
        if r.URL.Path != "/projects.json" {
            t.Errorf("Expected path /projects.json, got %s", r.URL.Path)
        }

        // Send mock response
        response := ProjectsResponse{
            Projects: []Project{{ID: 1, Name: "Test"}},
        }
        w.Header().Set("Content-Type", "application/json")
        _ = json.NewEncoder(w).Encode(response)
    }))
    defer server.Close()

    client := New(server.URL, "test-api-key")
    result, err := client.ListProjects(nil)
    if err != nil {
        t.Fatalf("ListProjects failed: %v", err)
    }
    // Assert results...
}
```

## Commands

### Linting
```bash
golangci-lint run
```
All code must pass linting before committing. The configuration uses `default: all` linters with custom formatters including gci, gofmt, gofumpt, goimports, and swaggo.

### Testing
```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./pkg/redmine -v

# Run specific test
go test ./pkg/redmine -v -run TestListProjects

# Run tests with coverage
go test ./pkg/redmine -coverprofile=coverage.out
go tool cover -html=coverage.out  # View coverage in browser
go tool cover -func=coverage.out  # View coverage summary
```

Current test coverage: ~49% with 60 comprehensive unit tests across all 22 Redmine REST APIs.

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
- Write comprehensive unit tests for all new API implementations
- Use `httptest.NewServer` for HTTP mocking in tests (avoid manual HTTP server setup)
- Validate HTTP methods, paths, headers, and request/response data in tests
- Test both success cases and error cases where applicable

## API Coverage

This SDK provides comprehensive coverage of all Redmine REST APIs:

### Stable APIs
- **Projects**: CRUD operations, project management
- **Issues**: CRUD operations, filtering, watchers, relations
- **Users**: User management, current user
- **Time Entries**: Time tracking, filtering

### Alpha/Beta APIs
- **Versions**: Project versions management
- **Memberships**: Project membership management
- **Wiki Pages**: Wiki CRUD operations
- **News**: News listing
- **Queries**: Issue queries
- **Custom Fields**: Custom field definitions
- **Search**: Global search
- **Files**: File management
- **My Account**: Current user account
- **Journals**: Issue journal entries
- **Issue Relations**: Issue relationship management
- **Attachments**: Attachment operations
- **Trackers**: Tracker listing
- **Issue Statuses**: Status definitions
- **Issue Categories**: Category management
- **Roles**: Role management
- **Groups**: User group management
- **Enumerations**: Issue priorities, time entry activities, document categories
