package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListCustomFields(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/custom_fields.json" {
			t.Errorf("Expected path /custom_fields.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := CustomFieldsResponse{
			CustomFields: []CustomFieldDefinition{
				{ID: 1, Name: "Custom Field 1", FieldFormat: "string"},
				{ID: 2, Name: "Custom Field 2", FieldFormat: "int"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListCustomFields()
	if err != nil {
		t.Fatalf("ListCustomFields failed: %v", err)
	}

	if len(result.CustomFields) != 2 {
		t.Errorf("Expected 2 custom fields, got %d", len(result.CustomFields))
	}
	if result.CustomFields[0].Name != "Custom Field 1" {
		t.Errorf("Expected name 'Custom Field 1', got %s", result.CustomFields[0].Name)
	}
}
