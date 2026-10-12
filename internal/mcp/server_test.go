package mcp

import (
	"context"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/mcp/config"
)

// listTools starts the server with every tool enabled and returns the tools
// as a client sees them.
func listTools(t *testing.T) []*mcp.Tool {
	t.Helper()

	server, err := NewServer(&config.Config{RedmineURL: "http://localhost:1", APIKey: "test-api-key"})
	if err != nil {
		t.Fatalf("NewServer failed: %v", err)
	}

	ctx := context.Background()
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	serverSession, err := server.Connect(ctx, serverTransport, nil)
	if err != nil {
		t.Fatalf("server.Connect failed: %v", err)
	}
	//nolint:errcheck
	defer serverSession.Close()

	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "0"}, nil)
	clientSession, err := client.Connect(ctx, clientTransport, nil)
	if err != nil {
		t.Fatalf("client.Connect failed: %v", err)
	}
	//nolint:errcheck
	defer clientSession.Close()

	result, err := clientSession.ListTools(ctx, nil)
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	return result.Tools
}

// annotationKind classifies annotations into the kinds used by the handlers.
func annotationKind(a *mcp.ToolAnnotations) string {
	switch {
	case a.ReadOnlyHint:
		return "read-only"
	case a.DestructiveHint == nil:
		return "unspecified"
	case *a.DestructiveHint && a.IdempotentHint:
		return "destructive-idempotent"
	case *a.DestructiveHint:
		return "destructive"
	case a.IdempotentHint:
		return "additive-idempotent"
	default:
		return "additive"
	}
}

// wantKindByName holds the expected kind of tools whose name does not follow
// the list_/show_/create_/update_/delete_ convention checked below.
var wantKindByName = map[string]string{
	"search":                     "read-only",
	"current_user":               "read-only",
	"analyze_project_health":     "read-only",
	"adjust_estimates":           "read-only",
	"upload_file":                "additive",
	"add_watcher":                "additive-idempotent",
	"add_group_user":             "additive-idempotent",
	"create_or_update_wiki_page": "destructive",
	"suggest_reschedule":         "destructive",
	"remove_watcher":             "destructive-idempotent",
	"remove_group_user":          "destructive-idempotent",
	"archive_project":            "destructive-idempotent",
	"unarchive_project":          "destructive-idempotent",
}

func TestAllToolsHaveAnnotations(t *testing.T) {
	tools := listTools(t)
	if len(tools) == 0 {
		t.Fatal("Expected tools to be registered, got none")
	}

	for _, tool := range tools {
		a := tool.Annotations
		if a == nil {
			t.Errorf("%s: missing annotations", tool.Name)
			continue
		}

		// Every tool that is not read-only must state whether it is destructive,
		// since clients otherwise assume it is.
		if !a.ReadOnlyHint && a.DestructiveHint == nil {
			t.Errorf("%s: DestructiveHint must be set for a tool that is not read-only", tool.Name)
		}

		want, ok := wantKindByName[tool.Name]
		if !ok {
			switch {
			case strings.HasPrefix(tool.Name, "list_"), strings.HasPrefix(tool.Name, "show_"):
				want = "read-only"
			case strings.HasPrefix(tool.Name, "create_"):
				want = "additive"
			case strings.HasPrefix(tool.Name, "update_"), strings.HasPrefix(tool.Name, "delete_"):
				want = "destructive-idempotent"
			default:
				t.Errorf("%s: add the expected annotation kind to wantKindByName", tool.Name)
				continue
			}
		}

		if got := annotationKind(a); got != want {
			t.Errorf("%s: annotation kind = %s, want %s", tool.Name, got, want)
		}
	}
}
