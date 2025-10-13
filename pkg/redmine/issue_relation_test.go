package redmine

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestListIssueRelations(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("Expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/issues/123/relations.json" {
			t.Errorf("Expected path /issues/123/relations.json, got %s", r.URL.Path)
		}
		if r.Header.Get("X-Redmine-API-Key") != "test-api-key" {
			t.Errorf("Expected API key header")
		}

		response := IssueRelationsResponse{
			Relations: []IssueRelation{
				{ID: 1, IssueID: 123, IssueToID: 456, RelationType: "relates"},
				{ID: 2, IssueID: 123, IssueToID: 789, RelationType: "blocks"},
			},
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	result, err := client.ListIssueRelations(123)
	if err != nil {
		t.Fatalf("ListIssueRelations failed: %v", err)
	}

	if len(result.Relations) != 2 {
		t.Errorf("Expected 2 relations, got %d", len(result.Relations))
	}
	if result.Relations[0].RelationType != "relates" {
		t.Errorf("Expected relation type 'relates', got %s", result.Relations[0].RelationType)
	}
}

func TestCreateIssueRelation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("Expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/issues/123/relations.json" {
			t.Errorf("Expected path /issues/123/relations.json, got %s", r.URL.Path)
		}

		var req IssueRelationRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}
		if req.Relation.IssueToID != 456 {
			t.Errorf("Expected issue_to_id 456, got %d", req.Relation.IssueToID)
		}
		if req.Relation.RelationType != "relates" {
			t.Errorf("Expected relation type 'relates', got %s", req.Relation.RelationType)
		}

		w.WriteHeader(http.StatusCreated)
		response := IssueRelationResponse{
			Relation: IssueRelation{
				ID:           1,
				IssueID:      123,
				IssueToID:    456,
				RelationType: req.Relation.RelationType,
			},
		}
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	relation := IssueRelation{
		IssueToID:    456,
		RelationType: "relates",
	}
	result, err := client.CreateIssueRelation(123, relation)
	if err != nil {
		t.Fatalf("CreateIssueRelation failed: %v", err)
	}

	if result.Relation.IssueToID != 456 {
		t.Errorf("Expected issue_to_id 456, got %d", result.Relation.IssueToID)
	}
}

func TestDeleteIssueRelation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/relations/1.json" {
			t.Errorf("Expected path /relations/1.json, got %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := New(server.URL, "test-api-key")
	err := client.DeleteIssueRelation(1)
	if err != nil {
		t.Fatalf("DeleteIssueRelation failed: %v", err)
	}
}
