package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// TimeEntryUseCase provides business logic for time entry operations.
type TimeEntryUseCase struct {
	client *redmine.Client
}

// NewTimeEntryUseCase creates a new TimeEntryUseCase instance.
func NewTimeEntryUseCase(client *redmine.Client) *TimeEntryUseCase {
	return &TimeEntryUseCase{
		client: client,
	}
}

// ListTimeEntries retrieves a list of time entries.
func (u *TimeEntryUseCase) ListTimeEntries(ctx context.Context, opts *redmine.ListTimeEntriesOptions) (*redmine.TimeEntriesResponse, error) {
	return u.client.ListTimeEntries(ctx, opts)
}

// ShowTimeEntry retrieves a single time entry by ID.
func (u *TimeEntryUseCase) ShowTimeEntry(ctx context.Context, id int) (*redmine.TimeEntryResponse, error) {
	return u.client.ShowTimeEntry(ctx, id)
}

// CreateTimeEntry creates a new time entry.
func (u *TimeEntryUseCase) CreateTimeEntry(ctx context.Context, req redmine.TimeEntryCreateRequest) (*redmine.TimeEntryResponse, error) {
	return u.client.CreateTimeEntry(ctx, req)
}

// UpdateTimeEntry updates an existing time entry.
func (u *TimeEntryUseCase) UpdateTimeEntry(ctx context.Context, id int, req redmine.TimeEntryUpdateRequest) error {
	return u.client.UpdateTimeEntry(ctx, id, req)
}

// DeleteTimeEntry deletes a time entry.
func (u *TimeEntryUseCase) DeleteTimeEntry(ctx context.Context, id int) error {
	return u.client.DeleteTimeEntry(ctx, id)
}
