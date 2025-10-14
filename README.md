# redmine-go

Unofficial Redmine MCP (Model Context Protocol) server implementation in Go.

## Overview

This project provides a comprehensive MCP server for Redmine, enabling AI assistants to interact with Redmine project management systems through the Model Context Protocol. The server exposes 76 tools covering all Redmine REST APIs.

## Features

### Implemented Services (76 Tools)

#### Projects (7 tools)
- `list_projects` - List all projects with filtering options
- `show_project` - Get detailed project information
- `create_project` - Create new projects
- `update_project` - Update existing projects
- `delete_project` - Delete projects
- `archive_project` - Archive projects
- `unarchive_project` - Unarchive projects

#### Issues (7 tools)
- `list_issues` - List issues with advanced filtering, sorting, and pagination
- `show_issue` - Get detailed issue information
- `create_issue` - Create new issues with full field support
- `update_issue` - Update existing issues
- `delete_issue` - Delete issues
- `add_watcher` - Add watchers to issues
- `remove_watcher` - Remove watchers from issues

#### Users (6 tools)
- `list_users` - List users with filtering options
- `show_user` - Get user details by ID
- `current_user` - Get current authenticated user
- `create_user` - Create new users (admin only)
- `update_user` - Update existing users (admin only)
- `delete_user` - Delete users (admin only)

#### Issue Categories (5 tools)
- `list_issue_categories` - List categories for a project
- `show_issue_category` - Get category details
- `create_issue_category` - Create new category
- `update_issue_category` - Update existing category
- `delete_issue_category` - Delete category with reassignment option

#### Time Entries (5 tools)
- `list_time_entries` - List time entries with filtering options
- `show_time_entry` - Get time entry details
- `create_time_entry` - Create new time entry
- `update_time_entry` - Update existing time entry
- `delete_time_entry` - Delete time entry

#### Versions (5 tools)
- `list_versions` - List project versions
- `show_version` - Get version details
- `create_version` - Create new version
- `update_version` - Update existing version
- `delete_version` - Delete version

#### Memberships (5 tools)
- `list_memberships` - List project memberships
- `show_membership` - Get membership details
- `create_membership` - Add user/group to project with roles
- `update_membership` - Update membership roles
- `delete_membership` - Remove user/group from project

#### Issue Relations (4 tools)
- `list_issue_relations` - List issue relations
- `show_issue_relation` - Get issue relation details
- `create_issue_relation` - Create issue relation (relates, duplicates, blocks, etc.)
- `delete_issue_relation` - Delete issue relation

#### Wiki Pages (4 tools)
- `list_wiki_pages` - List wiki pages in a project
- `show_wiki_page` - Get wiki page content
- `create_or_update_wiki_page` - Create or update wiki page
- `delete_wiki_page` - Delete wiki page

#### Attachments (3 tools)
- `show_attachment` - Get attachment details
- `update_attachment` - Update attachment metadata
- `delete_attachment` - Delete attachment

#### Enumerations (3 tools)
- `list_issue_priorities` - List issue priorities
- `list_time_entry_activities` - List time entry activities
- `list_document_categories` - List document categories

#### Groups (7 tools)
- `list_groups` - List all groups (admin only)
- `show_group` - Get group details (admin only)
- `create_group` - Create new group (admin only)
- `update_group` - Update existing group (admin only)
- `delete_group` - Delete group (admin only)
- `add_group_user` - Add user to group (admin only)
- `remove_group_user` - Remove user from group (admin only)

#### News (2 tools)
- `list_news` - List all news from all projects
- `list_project_news` - List news for a specific project

#### Files (2 tools)
- `list_files` - List files in a project
- `upload_file` - Upload file to a project

#### Roles (2 tools)
- `list_roles` - List all roles
- `show_role` - Get role details

#### Metadata (2 tools)
- `list_trackers` - List all available trackers (Bug, Feature, Support, etc.)
- `list_issue_statuses` - List all available issue statuses

#### My Account (2 tools)
- `show_my_account` - Get current user's account details
- `update_my_account` - Update current user's account information

#### Search (1 tool)
- `search` - Universal search across issues, wiki pages, and attachments

#### Queries (1 tool)
- `list_queries` - List saved issue queries

#### Custom Fields (1 tool)
- `list_custom_fields` - List custom field definitions

#### Journals (1 tool)
- `show_journal` - Get journal entry details

## Architecture

The project follows a clean, layered architecture:

```
MCP Client (Claude) → cmd/mcp-server → internal/mcp → internal/usecase → pkg/redmine → Redmine API
```

### Directory Structure

- `cmd/mcp-server/` - MCP server entry point
- `internal/mcp/` - MCP-specific implementation
  - `handlers/` - Tool handlers organized by service
  - `server.go` - Server initialization and tool registration
- `internal/usecase/` - Business logic layer (reusable across MCP/CLI)
- `internal/config/` - Configuration management
- `pkg/redmine/` - Redmine API client SDK (76 methods, 22 APIs)

## Installation

```bash
go install github.com/kqns91/redmine-go/cmd/mcp-server@latest
```

Or build from source:

```bash
git clone https://github.com/kqns91/redmine-go.git
cd redmine-go
go build -o mcp-server ./cmd/mcp-server
```

## Configuration

### Required Environment Variables

Set the following environment variables:

```bash
export REDMINE_URL="https://your-redmine-instance.com"
export REDMINE_API_KEY="your-api-key-here"
```

To get your Redmine API key:
1. Log in to your Redmine instance
2. Go to "My account" (top right)
3. Click "Show" under "API access key" on the right sidebar
4. Copy the displayed key

### Optional: Tool Control

You can control which MCP tools are enabled using these optional environment variables:

#### REDMINE_ENABLED_TOOLS
Comma-separated list of tool groups to enable. If not set, all tools are enabled by default.

Available tool groups:
- `projects` - Project management tools (7 tools)
- `issues` - Issue management tools (7 tools)
- `users` - User management tools (6 tools)
- `categories` - Issue category tools (5 tools)
- `time_entries` - Time tracking tools (5 tools)
- `versions` - Version management tools (5 tools)
- `memberships` - Project membership tools (5 tools)
- `issue_relations` - Issue relation tools (4 tools)
- `wiki` - Wiki page tools (4 tools)
- `attachments` - Attachment tools (3 tools)
- `enumerations` - Enumeration tools (3 tools)
- `groups` - Group management tools (7 tools, admin only)
- `news` - News tools (2 tools)
- `files` - File tools (2 tools)
- `roles` - Role tools (2 tools)
- `metadata` - Tracker and status metadata (2 tools)
- `my_account` - Current user account tools (2 tools)
- `search` - Search functionality (1 tool)
- `queries` - Saved query tools (1 tool)
- `custom_fields` - Custom field tools (1 tool)
- `journals` - Journal tools (1 tool)
- `all` - Enable all tool groups (default behavior)

#### REDMINE_DISABLED_TOOLS
Comma-separated list of individual tool names to disable. Takes precedence over `REDMINE_ENABLED_TOOLS`.

**Priority:** `REDMINE_DISABLED_TOOLS` > `REDMINE_ENABLED_TOOLS` > Default (all enabled)

#### Configuration Examples

**Pattern 1: Default (all tools enabled)**
```bash
# No configuration needed - all 76 tools are enabled by default
export REDMINE_URL="https://your-redmine-instance.com"
export REDMINE_API_KEY="your-api-key-here"
```

**Pattern 2: Enable only specific tool groups**
```bash
export REDMINE_ENABLED_TOOLS="projects,issues,search"
# Enables 15 tools: 7 project tools + 7 issue tools + 1 search tool
```

**Pattern 3: Enable all groups but disable specific tools**
```bash
export REDMINE_DISABLED_TOOLS="delete_project,delete_issue,delete_user,delete_group"
# Enables 72 tools (all tools except the 4 delete operations)
```

**Pattern 4: Read-only mode (disable all write operations)**
```bash
export REDMINE_ENABLED_TOOLS="all"
export REDMINE_DISABLED_TOOLS="create_project,update_project,delete_project,archive_project,unarchive_project,create_issue,update_issue,delete_issue,add_watcher,remove_watcher,create_user,update_user,delete_user,create_issue_category,update_issue_category,delete_issue_category,create_time_entry,update_time_entry,delete_time_entry,create_version,update_version,delete_version,create_membership,update_membership,delete_membership,create_issue_relation,delete_issue_relation,create_or_update_wiki_page,delete_wiki_page,update_attachment,delete_attachment,upload_file,create_group,update_group,delete_group,add_group_user,remove_group_user,update_my_account"
# Enables only read operations (list, show, get)
```

**Pattern 5: Minimal setup (projects and issues only, no destructive operations)**
```bash
export REDMINE_ENABLED_TOOLS="projects,issues,search"
export REDMINE_DISABLED_TOOLS="delete_project,delete_issue"
# Enables 13 tools: project/issue management + search, excluding delete operations
```

## Usage

### Running the MCP Server

```bash
./mcp-server
```

The server communicates via stdio using the Model Context Protocol.

### MCP Client Configuration

To use with Claude Desktop or other MCP clients, add to your MCP settings configuration:

**Basic Configuration (all tools enabled):**
```json
{
  "mcpServers": {
    "redmine": {
      "command": "/path/to/mcp-server",
      "env": {
        "REDMINE_URL": "https://your-redmine-instance.com",
        "REDMINE_API_KEY": "your-api-key-here"
      }
    }
  }
}
```

**With Tool Control (selective tools):**
```json
{
  "mcpServers": {
    "redmine": {
      "command": "/path/to/mcp-server",
      "env": {
        "REDMINE_URL": "https://your-redmine-instance.com",
        "REDMINE_API_KEY": "your-api-key-here",
        "REDMINE_ENABLED_TOOLS": "projects,issues,search",
        "REDMINE_DISABLED_TOOLS": "delete_project,delete_issue"
      }
    }
  }
}
```

## Development

### Requirements

- Go 1.23 or later
- golangci-lint for code quality checks

### Running Tests

```bash
go test ./...
```

### Linting

```bash
golangci-lint run
```

### Project Guidelines

- All code must pass `golangci-lint run` before committing
- Maintain test coverage for the SDK layer (`pkg/redmine`)
- Follow the layered architecture pattern
- Keep MCP handlers thin - business logic belongs in `internal/usecase`
- Use typed handlers with JSON schema annotations for better tool descriptions

## Dependencies

- [go-sdk v1.0.0](https://github.com/modelcontextprotocol/go-sdk) - MCP protocol implementation
- Standard Go library for HTTP client and JSON handling

## API Coverage

The underlying Redmine API client SDK (`pkg/redmine`) provides comprehensive coverage:
- 76 client methods
- 22 Redmine REST APIs
- 48.9% test coverage with 60 tests

Supported APIs include: Projects, Issues, Users, Time Entries, News, Wiki Pages, Files, My Account, Journals, Queries, Custom Fields, Search, Attachments, Issue Relations, Versions, Project Memberships, Issue Categories, Trackers, Issue Statuses, Groups, Roles, and Enumerations.

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please ensure:
1. All tests pass (`go test ./...`)
2. Code passes linting (`golangci-lint run`)
3. Follow the existing architecture patterns
4. Add tests for new functionality

## Related Projects

- [Redmine REST API Documentation](https://www.redmine.org/projects/redmine/wiki/Rest_api)
- [Model Context Protocol](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk)
