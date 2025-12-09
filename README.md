# redmine-go

Unofficial Redmine API client in Go

[![Go Version](https://img.shields.io/badge/Go-1.25.2%2B-00ADD8?logo=go)](https://go.dev/)
[![Go Reference](https://pkg.go.dev/badge/github.com/kqns91/redmine-go.svg)](https://pkg.go.dev/github.com/kqns91/redmine-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Ask DeepWiki](https://deepwiki.com/badge.svg)](https://deepwiki.com/kqns91/redmine-go)

English | [日本語](README.ja.md)

## Overview

`redmine-go` is a Redmine REST API client written in Go. It provides three ways to interact with Redmine:

- **SDK** - Go package for building applications that integrate with Redmine
- **CLI** - Command-line tool for managing Redmine from the terminal
- **MCP Server** - Server implementation for AI assistants using Model Context Protocol

All three are built on the same SDK foundation, supporting 22 Redmine REST APIs with 76 methods.

---

## SDK

Go client package for interacting with Redmine REST API.

### Installation

```bash
go get github.com/kqns91/redmine-go
```

### Basic Usage

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

The SDK supports the following Redmine REST APIs:

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

## CLI

Command-line tool for managing Redmine from the terminal.

### Installation

```bash
go install github.com/kqns91/redmine-go/cmd/redmine@latest
```

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

## MCP Server

The MCP (Model Context Protocol) server enables AI assistants to interact with Redmine.

### Installation

```bash
go install github.com/kqns91/redmine-go/cmd/mcp@latest
```

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

The server provides 80 tools across 23 categories:

**Core Resources**
- Projects (7 tools)
- Issues (7 tools)
- Users (6 tools)
- Issue Categories (5 tools)
- Time Entries (5 tools)
- Versions (5 tools)

**Advanced Operations**
- Batch Operations (1 tool) - Create multiple related tasks at once
- Progress Monitoring (3 tools) - Analyze project health, adjust estimates, suggest reschedules

**Project Management**
- Memberships (5 tools)
- Issue Relations (4 tools)
- Groups (7 tools)

**Content & Documentation**
- Wiki Pages (4 tools)
- Attachments (3 tools)
- News (2 tools)
- Files (2 tools)

**Metadata & Configuration**
- Enumerations (3 tools)
- Roles (2 tools)
- Metadata (2 tools)
- Custom Fields (1 tool)
- Queries (1 tool)

**User Account**
- My Account (2 tools)
- Search (1 tool)
- Journals (1 tool)

### Batch Operations

The `create_task_tree` tool enables efficient creation of multiple related tasks with dependencies:

```json
{
  "project_id": 1,
  "tasks": [
    {
      "ref": "backend",
      "subject": "Backend Development",
      "tracker_id": 1,
      "status_id": 1,
      "priority_id": 1,
      "assigned_to_id": 2,
      "estimated_hours": 24,
      "start_date": "2025-11-10",
      "due_date": "2025-11-13"
    },
    {
      "ref": "frontend",
      "subject": "Frontend Development",
      "parent_ref": "backend",
      "assigned_to_id": 3,
      "estimated_hours": 20,
      "start_date": "2025-11-14",
      "due_date": "2025-11-17",
      "blocks_refs": ["backend"]
    }
  ]
}
```

Features:
- Parent-child task relationships via `parent_ref`
- Task dependencies with `blocks_refs` and `precedes_refs`
- Automatic assignee distribution with `assigned_to_id`
- Estimated hours, start/due dates, custom fields support
- Ideal for creating 10-30 related tickets for a feature

### Progress Monitoring

Three tools help monitor and manage project progress:

**`analyze_project_health`** - Comprehensive project health analysis:
- Lists on-track, at-risk, and delayed issues
- Identifies critical path tasks
- Calculates delay days and impact levels
- Provides actionable recommendations

**`adjust_estimates`** - Smart estimate adjustments:
- Analyzes actual time entries vs. estimates
- Forecasts completion dates based on current progress
- Includes child issues in calculations
- Helps maintain realistic schedules

**`suggest_reschedule`** - Automatic rescheduling:
- Detects delayed tasks and dependency conflicts
- Suggests new dates with configurable buffer days
- Can auto-apply changes or just preview
- Supports critical-path-only mode

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
`projects`, `issues`, `users`, `categories`, `time_entries`, `versions`, `memberships`, `issue_relations`, `wiki`, `attachments`, `enumerations`, `groups`, `news`, `files`, `roles`, `metadata`, `my_account`, `search`, `queries`, `custom_fields`, `journals`, `batch_operations`, `progress_monitoring`, `all`

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
