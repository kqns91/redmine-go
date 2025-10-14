package redmine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListProjects(t *testing.T) {
	// モックサーバーを作成
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// リクエストの検証
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/projects.json" {
			t.Errorf("Expected path /projects.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-Api-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		// モックレスポンス
		response := ProjectsResponse{
			Projects: []Project{
				{ID: 1, Name: "Test Project 1", Identifier: "test-project-1"},
				{ID: 2, Name: "Test Project 2", Identifier: "test-project-2"},
			},
			TotalCount: 2,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// クライアントを作成
	client := New(server.URL, "test-api-key")

	// テスト実行
	result, err := client.ListProjects(context.Background(), nil)
	if err != nil {
		t.Fatalf("ListProjects failed: %v", err)
	}

	// 結果の検証
	if len(result.Projects) != 2 {
		t.Errorf("Expected 2 projects, got %d", len(result.Projects))
	}
	if result.Projects[0].Name != "Test Project 1" {
		t.Errorf("Expected project name 'Test Project 1', got '%s'", result.Projects[0].Name)
	}
	if result.TotalCount != 2 {
		t.Errorf("Expected total count 2, got %d", result.TotalCount)
	}
}

func TestListProjectsWithOptions(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// クエリパラメータの検証
		if r.URL.Query().Get("limit") != "10" {
			t.Errorf("Expected limit=10, got %s", r.URL.Query().Get("limit"))
		}
		if r.URL.Query().Get("offset") != "20" {
			t.Errorf("Expected offset=20, got %s", r.URL.Query().Get("offset"))
		}
		if r.URL.Query().Get("include") != "trackers" {
			t.Errorf("Expected include=trackers, got %s", r.URL.Query().Get("include"))
		}

		response := ProjectsResponse{Projects: []Project{}}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	opts := &ListProjectsOptions{
		Limit:   10,
		Offset:  20,
		Include: "trackers",
	}
	_, err := client.ListProjects(context.Background(), opts)
	if err != nil {
		t.Fatalf("ListProjects with options failed: %v", err)
	}
}

func TestShowProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/projects/test-project.json" {
			t.Errorf("Expected path /projects/test-project.json, got %s", r.URL.Path)
		}

		response := ProjectResponse{
			Project: Project{
				ID:          1,
				Name:        "Test Project",
				Identifier:  "test-project",
				Description: "Test Description",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	result, err := client.ShowProject(context.Background(), "test-project", nil)
	if err != nil {
		t.Fatalf("ShowProject failed: %v", err)
	}

	if result.Project.Name != "Test Project" {
		t.Errorf("Expected project name 'Test Project', got '%s'", result.Project.Name)
	}
}

func TestCreateProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}

		// リクエストボディの検証
		var req ProjectRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Project.Name != "New Project" {
			t.Errorf("Expected project name 'New Project', got '%s'", req.Project.Name)
		}

		// 201 Created レスポンス
		w.WriteHeader(http.StatusCreated)
		response := ProjectResponse{
			Project: Project{
				ID:         3,
				Name:       req.Project.Name,
				Identifier: req.Project.Identifier,
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	project := Project{
		Name:       "New Project",
		Identifier: "new-project",
	}
	result, err := client.CreateProject(context.Background(), project)
	if err != nil {
		t.Fatalf("CreateProject failed: %v", err)
	}

	if result.Project.ID != 3 {
		t.Errorf("Expected project ID 3, got %d", result.Project.ID)
	}
}

func TestCreateProjectError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 422 Unprocessable Entity レスポンス
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"errors":["Name can't be blank"]}`))
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	project := Project{
		Identifier: "test",
	}
	_, err := client.CreateProject(context.Background(), project)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}
}

func TestDeleteProject(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/projects/1.json" {
			t.Errorf("Expected path /projects/1.json, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")

	err := client.DeleteProject(context.Background(), "1")
	if err != nil {
		t.Fatalf("DeleteProject failed: %v", err)
	}
}
