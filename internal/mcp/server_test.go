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

		switch {
		case strings.HasPrefix(tool.Name, "list_"), strings.HasPrefix(tool.Name, "show_"):
			if !a.ReadOnlyHint {
				t.Errorf("%s: expected ReadOnlyHint", tool.Name)
			}
		case strings.HasPrefix(tool.Name, "delete_"), strings.HasPrefix(tool.Name, "update_"):
			if a.ReadOnlyHint || a.DestructiveHint == nil || !*a.DestructiveHint {
				t.Errorf("%s: expected DestructiveHint", tool.Name)
			}
		case strings.HasPrefix(tool.Name, "create_") && tool.Name != "create_or_update_wiki_page":
			if a.ReadOnlyHint || a.DestructiveHint == nil || *a.DestructiveHint {
				t.Errorf("%s: expected a non-destructive write tool", tool.Name)
			}
		}
	}
}
