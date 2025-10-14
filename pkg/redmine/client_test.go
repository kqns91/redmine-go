package redmine

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestNew(t *testing.T) {
	client := New("https://example.com", "test-api-key")

	if client.baseURL != "https://example.com" {
		t.Errorf("Expected baseURL https://example.com, got %s", client.baseURL)
	}
	if client.apiKey != "test-api-key" {
		t.Errorf("Expected apiKey test-api-key, got %s", client.apiKey)
	}
	if client.HTTPClient == nil {
		t.Error("Expected HTTPClient to be initialized")
	}
}

func TestNewTrimsTrailingSlash(t *testing.T) {
	client := New("https://example.com/", "test-api-key")

	if client.baseURL != "https://example.com" {
		t.Errorf("Expected baseURL without trailing slash, got %s", client.baseURL)
	}
}

func TestClientDoSetsHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey := r.Header.Get("X-Redmine-Api-Key")
		if apiKey != "test-api-key" {
			t.Errorf("Expected X-Redmine-API-Key: test-api-key, got %s", apiKey)
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type: application/json, got %s", contentType)
		}

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	_, err := client.do(context.Background(), http.MethodGet, server.URL+"/test", nil)
	if err != nil {
		t.Fatalf("do() failed: %v", err)
	}
}

func TestClientDoHandles4xxError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(`{"error":"Not found"}`))
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	_, err := client.do(context.Background(), http.MethodGet, server.URL+"/test", nil)
	if err == nil {
		t.Error("Expected error for 404 response, got nil")
	}
}

func TestClientDoHandles5xxError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"Internal server error"}`))
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	_, err := client.do(context.Background(), http.MethodGet, server.URL+"/test", nil)
	if err == nil {
		t.Error("Expected error for 500 response, got nil")
	}
}

func TestClientDoHandlesNetworkError(t *testing.T) {
	// 存在しないサーバーに接続
	client := New("http://localhost:1", "test-api-key")
	_, err := client.do(context.Background(), http.MethodGet, "http://localhost:1/test", nil)
	if err == nil {
		t.Error("Expected network error, got nil")
	}
}

func TestClientDoSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"success":true}`))
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	resp, err := client.do(context.Background(), http.MethodGet, server.URL+"/test", nil)
	if err != nil {
		t.Fatalf("do() failed: %v", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}
