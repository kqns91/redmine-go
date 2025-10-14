package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListWikiPages(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test-project/wiki/index.json" {
			t.Errorf("Expected path /projects/test-project/wiki/index.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-Api-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := WikiPagesResponse{
			WikiPages: []WikiPageIndex{
				{Title: "Wiki Page 1", Version: 1},
				{Title: "Wiki Page 2", Version: 2},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListWikiPages(context.Background(), "test-project")
	if err != nil {
		t.Fatalf("ListWikiPages failed: %v", err)
	}

	if len(result.WikiPages) != 2 {
		t.Errorf("Expected 2 wiki pages, got %d", len(result.WikiPages))
	}
	if result.WikiPages[0].Title != "Wiki Page 1" {
		t.Errorf("Expected title 'Wiki Page 1', got %s", result.WikiPages[0].Title)
	}
}

func TestShowWikiPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test-project/wiki/TestPage.json" {
			t.Errorf("Expected path /projects/test-project/wiki/TestPage.json, got %s", r.URL.Path)
		}

		response := WikiPageResponse{
			WikiPage: WikiPage{
				Title:   "TestPage",
				Text:    "Test content",
				Version: 1,
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.GetWikiPage(context.Background(), "test-project", "TestPage", nil)
	if err != nil {
		t.Fatalf("ShowWikiPage failed: %v", err)
	}

	if result.WikiPage.Title != "TestPage" {
		t.Errorf("Expected title 'TestPage', got %s", result.WikiPage.Title)
	}
	if result.WikiPage.Text != "Test content" {
		t.Errorf("Expected text 'Test content', got %s", result.WikiPage.Text)
	}
}

func TestCreateWikiPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("Expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test-project/wiki/NewPage.json" {
			t.Errorf("Expected path /projects/test-project/wiki/NewPage.json, got %s", r.URL.Path)
		}

		var req WikiPageRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.WikiPage.Text != "New page content" {
			t.Errorf("Expected text 'New page content', got %s", req.WikiPage.Text)
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	wikiPage := WikiPageUpdate{
		Text: "New page content",
	}
	err := client.CreateOrUpdateWikiPage(context.Background(), "test-project", "NewPage", wikiPage)
	if err != nil {
		t.Fatalf("CreateOrUpdateWikiPage failed: %v", err)
	}
}

func TestDeleteWikiPage(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test-project/wiki/OldPage.json" {
			t.Errorf("Expected path /projects/test-project/wiki/OldPage.json, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	err := client.DeleteWikiPage(context.Background(), "test-project", "OldPage")
	if err != nil {
		t.Fatalf("DeleteWikiPage failed: %v", err)
	}
}
