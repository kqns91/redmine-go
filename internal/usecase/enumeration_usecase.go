package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// EnumerationUseCase provides business logic for enumeration operations.
type EnumerationUseCase struct {
	client *redmine.Client
}

// NewEnumerationUseCase creates a new EnumerationUseCase instance.
func NewEnumerationUseCase(client *redmine.Client) *EnumerationUseCase {
	return &EnumerationUseCase{
		client: client,
	}
}

// ListIssuePriorities retrieves the list of issue priorities.
func (u *EnumerationUseCase) ListIssuePriorities(ctx context.Context) (*redmine.EnumerationsResponse, error) {
	return u.client.ListIssuePriorities(ctx)
}

// ListTimeEntryActivities retrieves the list of time entry activities.
func (u *EnumerationUseCase) ListTimeEntryActivities(ctx context.Context) (*redmine.EnumerationsResponse, error) {
	return u.client.ListTimeEntryActivities(ctx)
}

// ListDocumentCategories retrieves the list of document categories.
func (u *EnumerationUseCase) ListDocumentCategories(ctx context.Context) (*redmine.EnumerationsResponse, error) {
	return u.client.ListDocumentCategories(ctx)
}
