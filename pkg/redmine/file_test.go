package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListFiles(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test-project/files.json" {
			t.Errorf("Expected path /projects/test-project/files.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := FilesResponse{
			Files: []File{
				{ID: 1, Filename: "file1.txt", Filesize: 1024},
				{ID: 2, Filename: "file2.pdf", Filesize: 2048},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListFiles("test-project")
	if err != nil {
		t.Fatalf("ListFiles failed: %v", err)
	}

	if len(result.Files) != 2 {
		t.Errorf("Expected 2 files, got %d", len(result.Files))
	}
	if result.Files[0].Filename != "file1.txt" {
		t.Errorf("Expected filename 'file1.txt', got %s", result.Files[0].Filename)
	}
}

func TestUploadFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/projects/test-project/files.json" {
			t.Errorf("Expected path /projects/test-project/files.json, got %s", r.URL.Path)
		}

		var req FileUploadRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.File.Token == "" {
			t.Errorf("Expected token to be set")
		}

		w.WriteHeader(http.StatusCreated)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	file := FileUpload{
		Token:    "test-token",
		Filename: "newfile.txt",
	}
	err := client.UploadFile("test-project", file)
	if err != nil {
		t.Fatalf("UploadFile failed: %v", err)
	}
}
