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

func TestClientDoKeepsAPIKeyOnSameHostRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/old" {
			http.Redirect(w, r, "/new", http.StatusFound)
			return
		}
		if apiKey := r.Header.Get("X-Redmine-Api-Key"); apiKey != "test-api-key" {
			t.Errorf("Expected API key to be kept on same-host redirect, got %q", apiKey)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	resp, err := client.do(context.Background(), http.MethodGet, server.URL+"/old", nil)
	if err != nil {
		t.Fatalf("do() failed: %v", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()
}

func TestClientDoStripsAPIKeyOnCrossHostRedirect(t *testing.T) {
	tests := []struct {
		name       string
		httpClient *http.Client
	}{
		{name: "default client", httpClient: nil},
		{name: "replaced client", httpClient: &http.Client{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			other := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if apiKey := r.Header.Get("X-Redmine-Api-Key"); apiKey != "" {
					t.Errorf("Expected API key to be stripped on cross-host redirect, got %q", apiKey)
				}
				w.WriteHeader(http.StatusOK)
			}))
			defer other.Close()

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				http.Redirect(w, r, other.URL+"/landing", http.StatusFound)
			}))
			defer server.Close()

			client := New(server.URL, "test-api-key")
			if tt.httpClient != nil {
				client.HTTPClient = tt.httpClient
			}
			resp, err := client.do(context.Background(), http.MethodGet, server.URL+"/test", nil)
			if err != nil {
				t.Fatalf("do() failed: %v", err)
			}
			//nolint:errcheck
			defer resp.Body.Close()
		})
	}
}

func TestClientDoCallsCustomCheckRedirect(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/new", http.StatusFound)
	}))
	defer server.Close()

	called := false
	client := New(server.URL, "test-api-key")
	client.HTTPClient = &http.Client{
		CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
			called = true
			return http.ErrUseLastResponse
		},
	}

	resp, err := client.do(context.Background(), http.MethodGet, server.URL+"/old", nil)
	if err != nil {
		t.Fatalf("do() failed: %v", err)
	}
	//nolint:errcheck
	defer resp.Body.Close()

	if !called {
		t.Error("Expected custom CheckRedirect to be called")
	}
	if resp.StatusCode != http.StatusFound {
		t.Errorf("Expected status code 302 from custom CheckRedirect, got %d", resp.StatusCode)
	}
}

func TestClientDoStopsAfterTooManyRedirects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/loop", http.StatusFound)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	_, err := client.do(context.Background(), http.MethodGet, server.URL+"/loop", nil)
	if err == nil {
		t.Error("Expected error for redirect loop, got nil")
	}
}

func TestHTTPClientCheckRedirectStripsAPIKeyOnDowngrade(t *testing.T) {
	tests := []struct {
		name    string
		from    string
		to      string
		wantKey bool
	}{
		{name: "https to https on same host", from: "https://redmine.example.com/a", to: "https://redmine.example.com/b", wantKey: true},
		{name: "https to http on same host", from: "https://redmine.example.com/a", to: "http://redmine.example.com/b", wantKey: false},
		{name: "http to https on same host", from: "http://redmine.example.com/a", to: "https://redmine.example.com/b", wantKey: true},
		{name: "different port", from: "https://redmine.example.com/a", to: "https://redmine.example.com:8443/b", wantKey: false},
		{name: "different host", from: "https://redmine.example.com/a", to: "https://attacker.example.net/b", wantKey: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client := New("https://redmine.example.com", "test-api-key")

			via, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tt.from, nil)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequestWithContext(context.Background(), http.MethodGet, tt.to, nil)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("X-Redmine-Api-Key", "test-api-key")

			if err := client.httpClient().CheckRedirect(req, []*http.Request{via}); err != nil {
				t.Fatalf("CheckRedirect failed: %v", err)
			}

			if got := req.Header.Get("X-Redmine-Api-Key") != ""; got != tt.wantKey {
				t.Errorf("API key kept = %v, want %v", got, tt.wantKey)
			}
		})
	}
}

func TestPathSegment(t *testing.T) {
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "identifier", in: "my-project_1", want: "my-project_1"},
		{name: "numeric id", in: "42", want: "42"},
		{name: "wiki title", in: "Getting_Started", want: "Getting_Started"},
		{name: "non-ASCII", in: "\u65e5\u672c\u8a9e\u30da\u30fc\u30b8", want: "%E6%97%A5%E6%9C%AC%E8%AA%9E%E3%83%9A%E3%83%BC%E3%82%B8"},
		{name: "space", in: "My Page", want: "My%20Page"},
		{name: "query and fragment", in: "Page?x=1#y", want: "Page%3Fx=1%23y"},
		{name: "percent", in: "100%", want: "100%25"},
		{name: "backslash", in: `a\b`, want: "a%5Cb"},
		{name: "empty", in: "", wantErr: true},
		{name: "dot", in: ".", wantErr: true},
		{name: "double dot", in: "..", wantErr: true},
		{name: "slash", in: "a/b", wantErr: true},
		{name: "traversal", in: "../../projects/victim", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pathSegment(tt.in)
			if tt.wantErr {
				if err == nil {
					t.Errorf("pathSegment(%q) = %q, want error", tt.in, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("pathSegment(%q) failed: %v", tt.in, err)
			}
			if got != tt.want {
				t.Errorf("pathSegment(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestPathSegmentsCannotChangeEndpoint(t *testing.T) {
	var gotPath, gotRawQuery string
	requests := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests++
		gotPath = r.URL.EscapedPath()
		gotRawQuery = r.URL.RawQuery
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	ctx := context.Background()

	// Traversal is rejected before any request is sent.
	if err := client.DeleteWikiPage(ctx, "demo", "../../../../projects/victim"); err == nil {
		t.Error("Expected error for traversal in wiki page name, got nil")
	}
	if err := client.DeleteProject(ctx, ".."); err == nil {
		t.Error("Expected error for '..' project identifier, got nil")
	}
	if requests != 0 {
		t.Fatalf("Expected no request for invalid path segments, got %d", requests)
	}

	// "?" and "#" stay inside the page name instead of starting a query.
	if err := client.DeleteWikiPage(ctx, "demo", "Page?foo=bar#x"); err != nil {
		t.Fatalf("DeleteWikiPage failed: %v", err)
	}
	if want := "/projects/demo/wiki/Page%3Ffoo=bar%23x.json"; gotPath != want {
		t.Errorf("Expected path %s, got %s", want, gotPath)
	}
	if gotRawQuery != "" {
		t.Errorf("Expected no query string, got %q", gotRawQuery)
	}
}
