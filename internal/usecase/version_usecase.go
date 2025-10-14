package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// VersionUseCase provides business logic for version operations.
type VersionUseCase struct {
	client *redmine.Client
}

// NewVersionUseCase creates a new VersionUseCase instance.
func NewVersionUseCase(client *redmine.Client) *VersionUseCase {
	return &VersionUseCase{
		client: client,
	}
}

// ListVersions retrieves versions for a specific project.
func (u *VersionUseCase) ListVersions(ctx context.Context, projectIDOrIdentifier string) (*redmine.VersionsResponse, error) {
	return u.client.ListVersions(ctx, projectIDOrIdentifier)
}

// ShowVersion retrieves a specific version by ID.
func (u *VersionUseCase) ShowVersion(ctx context.Context, id int) (*redmine.VersionResponse, error) {
	return u.client.ShowVersion(ctx, id)
}

// CreateVersion creates a new version for a project.
func (u *VersionUseCase) CreateVersion(ctx context.Context, projectIDOrIdentifier string, version redmine.Version) (*redmine.VersionResponse, error) {
	return u.client.CreateVersion(ctx, projectIDOrIdentifier, version)
}

// UpdateVersion updates an existing version.
func (u *VersionUseCase) UpdateVersion(ctx context.Context, id int, version redmine.Version) error {
	return u.client.UpdateVersion(ctx, id, version)
}

// DeleteVersion deletes a version.
func (u *VersionUseCase) DeleteVersion(ctx context.Context, id int) error {
	return u.client.DeleteVersion(ctx, id)
}
