package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// IssueRelationUseCase provides business logic for issue relation operations.
type IssueRelationUseCase struct {
	client *redmine.Client
}

// NewIssueRelationUseCase creates a new IssueRelationUseCase instance.
func NewIssueRelationUseCase(client *redmine.Client) *IssueRelationUseCase {
	return &IssueRelationUseCase{
		client: client,
	}
}

// ListIssueRelations retrieves all relations for a specific issue.
func (u *IssueRelationUseCase) ListIssueRelations(ctx context.Context, issueID int) (*redmine.IssueRelationsResponse, error) {
	return u.client.ListIssueRelations(ctx, issueID)
}

// ShowIssueRelation retrieves details of a specific relation.
func (u *IssueRelationUseCase) ShowIssueRelation(ctx context.Context, id int) (*redmine.IssueRelationResponse, error) {
	return u.client.ShowIssueRelation(ctx, id)
}

// CreateIssueRelation creates a new issue relation.
func (u *IssueRelationUseCase) CreateIssueRelation(ctx context.Context, issueID int, relation redmine.IssueRelation) (*redmine.IssueRelationResponse, error) {
	return u.client.CreateIssueRelation(ctx, issueID, relation)
}

// DeleteIssueRelation deletes a specific issue relation.
func (u *IssueRelationUseCase) DeleteIssueRelation(ctx context.Context, id int) error {
	return u.client.DeleteIssueRelation(ctx, id)
}
