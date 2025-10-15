# redmine-go

Unofficial Redmine API implementation in Go.

[![Go Version](https://img.shields.io/badge/Go-1.25.2%2B-00ADD8?logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

English | [日本語](README.ja.md)

## Overview

`redmine-go` is a Go implementation of the Redmine REST API. It provides three ways to interact with Redmine:

- **SDK** - Go package for building applications that integrate with Redmine
- **CLI** - Command-line tool for managing Redmine from the terminal
- **MCP Server** - Server implementation for AI assistants using Model Context Protocol

All three are built on the same SDK foundation, supporting 22 Redmine REST APIs with 76 methods.

## Features

- Coverage of Redmine REST API (22 APIs, 76 methods)
- Context-aware operations using Go's context package
- Comprehensive error handling
- Multiple output formats (CLI: JSON, table, simple text)
- Flexible tool control (MCP: per-tool enable/disable)

## Installation

### SDK

```bash
go get github.com/kqns91/redmine-go
```

### CLI

```bash
go install github.com/kqns91/redmine-go/cmd/redmine@latest
```

### MCP Server

```bash
go install github.com/kqns91/redmine-go/cmd/mcp@latest
```

---

## SDK Usage

The SDK provides a Go client for interacting with Redmine REST API.

### Basic Example

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/kqns91/redmine-go/pkg/redmine"
)

func main() {
    client := redmine.New("https://your-redmine.com", "your-api-key")
    ctx := context.Background()

    // List projects
    projects, err := client.ListProjects(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }

    for _, project := range projects.Projects {
        fmt.Printf("%s (ID: %d)\n", project.Name, project.ID)
    }

    // Create an issue
    issue := redmine.Issue{
        ProjectID:   1,
        Subject:     "Sample issue",
        Description: "Issue description",
    }

    created, err := client.CreateIssue(ctx, issue)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Created issue #%d\n", created.Issue.ID)
}
```

### Supported APIs

The SDK supports Redmine REST APIs:

**Core Resources**
- Projects (CRUD, archive/unarchive)
- Issues (CRUD, watchers)
- Users (CRUD)
- Time Entries (CRUD)

**Project Management**
- Versions (CRUD)
- Issue Relations (CRUD)
- Memberships (CRUD)
- Issue Categories (CRUD)

**Content**
- Wiki Pages (CRUD)
- News (read)
- Files (read, upload)
- Attachments (read, update, delete)

**Administration**
- Groups (CRUD, user management)
- Roles (read)
- Trackers (read)
- Issue Statuses (read)
- Enumerations (priorities, activities, categories)

**Other**
- Custom Fields (read)
- Queries (read)
- Journals (read)
- My Account (read, update)
- Search

For detailed API documentation, see the [pkg/redmine](pkg/redmine) directory.

---

## CLI Usage

Command-line tool for managing Redmine from the terminal.

### Configuration

Run the config command to set up interactively:

```bash
redmine config set url https://your-redmine.com
redmine config set api_key your-api-key
```

View current configuration:

```bash
redmine config list
```

The configuration is stored at `~/.redmine/config.yaml`. You can also edit this file directly if needed.

Alternatively, you can use environment variables or command-line flags:

```bash
# Environment variables
export REDMINE_URL="https://your-redmine.com"
export REDMINE_API_KEY="your-api-key"

# Command-line flags
redmine --url https://your-redmine.com --api-key your-api-key <command>
```

### Getting Your API Key

1. Log in to your Redmine instance
2. Navigate to "My account" (top right)
3. Find "API access key" in the right sidebar
4. Click "Show" and copy the key

### Basic Commands

```bash
# Projects
redmine projects list
redmine projects show <project-id>

# Issues
redmine issues list --project <project-id>
redmine issues show <issue-id>
redmine issues create --project <project-id> --subject "Title" --description "Description"
redmine issues update <issue-id> --status <status-id> --assigned-to <user-id>

# Users
redmine users list
redmine users show <user-id>
redmine users current
```

### Output Formats

The CLI supports three output formats:

**Table format** (default)
```bash
redmine projects list --format table
```
Structured table with columns, suitable for terminal viewing.

**JSON format**
```bash
redmine projects list --format json
```
Machine-readable JSON output, useful for scripting and integration.

**Text format**
```bash
redmine projects list --format text
```
Plain text output with minimal formatting.

### Help

All commands provide detailed help:

```bash
redmine --help
redmine projects --help
redmine issues create --help
```

---

## MCP Server Usage

The MCP (Model Context Protocol) server enables AI assistants to interact with Redmine.

### Configuration

Add to your MCP client configuration file.

For example, with Claude Desktop:

**macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
**Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

Basic configuration (all tools enabled):

```json
{
  "mcpServers": {
    "redmine": {
      "command": "/path/to/mcp",
      "env": {
        "REDMINE_URL": "https://your-redmine.com",
        "REDMINE_API_KEY": "your-api-key"
      }
    }
  }
}
```

### Available Tools

The server provides 76 tools across 21 categories:

- Projects (7 tools)
- Issues (7 tools)
- Users (6 tools)
- Issue Categories (5 tools)
- Time Entries (5 tools)
- Versions (5 tools)
- Memberships (5 tools)
- Issue Relations (4 tools)
- Wiki Pages (4 tools)
- Attachments (3 tools)
- Enumerations (3 tools)
- Groups (7 tools)
- News (2 tools)
- Files (2 tools)
- Roles (2 tools)
- Metadata (2 tools)
- My Account (2 tools)
- Search (1 tool)
- Queries (1 tool)
- Custom Fields (1 tool)
- Journals (1 tool)

### Tool Control

You can control which tools are enabled using environment variables.

#### Enable Specific Tool Groups

Use `REDMINE_ENABLED_TOOLS` to specify which tool groups to enable:

```json
{
  "mcpServers": {
    "redmine": {
      "command": "/path/to/mcp",
      "env": {
        "REDMINE_URL": "https://your-redmine.com",
        "REDMINE_API_KEY": "your-api-key",
        "REDMINE_ENABLED_TOOLS": "projects,issues,search"
      }
    }
  }
}
```

Available tool groups:
`projects`, `issues`, `users`, `categories`, `time_entries`, `versions`, `memberships`, `issue_relations`, `wiki`, `attachments`, `enumerations`, `groups`, `news`, `files`, `roles`, `metadata`, `my_account`, `search`, `queries`, `custom_fields`, `journals`, `all`

#### Disable Specific Tools

Use `REDMINE_DISABLED_TOOLS` to disable individual tools:

```json
{
  "env": {
    "REDMINE_DISABLED_TOOLS": "delete_project,delete_issue,delete_user"
  }
}
```

This setting takes precedence over `REDMINE_ENABLED_TOOLS`.

#### Configuration Examples

**Read-only mode**

Disable all write operations:

```json
{
  "env": {
    "REDMINE_ENABLED_TOOLS": "all",
    "REDMINE_DISABLED_TOOLS": "create_project,update_project,delete_project,archive_project,unarchive_project,create_issue,update_issue,delete_issue,add_watcher,remove_watcher,create_user,update_user,delete_user,create_issue_category,update_issue_category,delete_issue_category,create_time_entry,update_time_entry,delete_time_entry,create_version,update_version,delete_version,create_membership,update_membership,delete_membership,create_issue_relation,delete_issue_relation,create_or_update_wiki_page,delete_wiki_page,update_attachment,delete_attachment,upload_file,create_group,update_group,delete_group,add_group_user,remove_group_user,update_my_account"
  }
}
```

**Project and issue management only**

```json
{
  "env": {
    "REDMINE_ENABLED_TOOLS": "projects,issues,search",
    "REDMINE_DISABLED_TOOLS": "delete_project,delete_issue"
  }
}
```

---

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Related Resources

- [Redmine REST API Documentation](https://www.redmine.org/projects/redmine/wiki/Rest_api)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
