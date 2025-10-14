package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// IssueUseCase provides business logic for issue operations.
type IssueUseCase struct {
	client *redmine.Client
}

// NewIssueUseCase creates a new IssueUseCase instance.
func NewIssueUseCase(client *redmine.Client) *IssueUseCase {
	return &IssueUseCase{
		client: client,
	}
}

// ListIssues retrieves a list of issues.
func (u *IssueUseCase) ListIssues(ctx context.Context, opts *redmine.ListIssuesOptions) (*redmine.IssuesResponse, error) {
	return u.client.ListIssues(ctx, opts)
}

// ShowIssue retrieves a single issue by ID.
func (u *IssueUseCase) ShowIssue(ctx context.Context, id int, opts *redmine.ShowIssueOptions) (*redmine.IssueResponse, error) {
	return u.client.ShowIssue(ctx, id, opts)
}

// CreateIssue creates a new issue.
func (u *IssueUseCase) CreateIssue(ctx context.Context, req redmine.IssueCreateRequest) (*redmine.IssueResponse, error) {
	return u.client.CreateIssue(ctx, req)
}

// UpdateIssue updates an existing issue.
func (u *IssueUseCase) UpdateIssue(ctx context.Context, id int, req redmine.IssueUpdateRequest) error {
	return u.client.UpdateIssue(ctx, id, req)
}

// DeleteIssue deletes an issue.
func (u *IssueUseCase) DeleteIssue(ctx context.Context, id int) error {
	return u.client.DeleteIssue(ctx, id)
}

// AddWatcher adds a watcher to an issue.
func (u *IssueUseCase) AddWatcher(ctx context.Context, issueID int, userID int) error {
	return u.client.AddWatcher(ctx, issueID, userID)
}

// RemoveWatcher removes a watcher from an issue.
func (u *IssueUseCase) RemoveWatcher(ctx context.Context, issueID int, userID int) error {
	return u.client.RemoveWatcher(ctx, issueID, userID)
}
