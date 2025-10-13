package usecase

import (
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
func (u *IssueUseCase) ListIssues(opts *redmine.ListIssuesOptions) (*redmine.IssuesResponse, error) {
	return u.client.ListIssues(opts)
}

// ShowIssue retrieves a single issue by ID.
func (u *IssueUseCase) ShowIssue(id int, opts *redmine.ShowIssueOptions) (*redmine.IssueResponse, error) {
	return u.client.ShowIssue(id, opts)
}

// CreateIssue creates a new issue.
func (u *IssueUseCase) CreateIssue(issue redmine.Issue) (*redmine.IssueResponse, error) {
	return u.client.CreateIssue(issue)
}

// UpdateIssue updates an existing issue.
func (u *IssueUseCase) UpdateIssue(id int, issue redmine.Issue) error {
	return u.client.UpdateIssue(id, issue)
}

// DeleteIssue deletes an issue.
func (u *IssueUseCase) DeleteIssue(id int) error {
	return u.client.DeleteIssue(id)
}

// AddWatcher adds a watcher to an issue.
func (u *IssueUseCase) AddWatcher(issueID int, userID int) error {
	return u.client.AddWatcher(issueID, userID)
}

// RemoveWatcher removes a watcher from an issue.
func (u *IssueUseCase) RemoveWatcher(issueID int, userID int) error {
	return u.client.RemoveWatcher(issueID, userID)
}
