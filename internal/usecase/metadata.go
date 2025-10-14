package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// MetadataUseCase provides business logic for metadata operations.
type MetadataUseCase struct {
	client *redmine.Client
}

// NewMetadataUseCase creates a new MetadataUseCase instance.
func NewMetadataUseCase(client *redmine.Client) *MetadataUseCase {
	return &MetadataUseCase{
		client: client,
	}
}

// ListTrackers retrieves a list of trackers.
func (u *MetadataUseCase) ListTrackers(ctx context.Context) (*redmine.TrackersResponse, error) {
	return u.client.ListTrackers(ctx)
}

// ListIssueStatuses retrieves a list of issue statuses.
func (u *MetadataUseCase) ListIssueStatuses(ctx context.Context) (*redmine.IssueStatusesResponse, error) {
	return u.client.ListIssueStatuses(ctx)
}
