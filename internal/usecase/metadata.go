package usecase

import (
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
func (u *MetadataUseCase) ListTrackers() (*redmine.TrackersResponse, error) {
	return u.client.ListTrackers()
}

// ListIssueStatuses retrieves a list of issue statuses.
func (u *MetadataUseCase) ListIssueStatuses() (*redmine.IssueStatusesResponse, error) {
	return u.client.ListIssueStatuses()
}
