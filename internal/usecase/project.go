package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// ProjectUseCase provides business logic for project operations.
type ProjectUseCase struct {
	client *redmine.Client
}

// NewProjectUseCase creates a new ProjectUseCase instance.
func NewProjectUseCase(client *redmine.Client) *ProjectUseCase {
	return &ProjectUseCase{
		client: client,
	}
}

// ListProjects retrieves a list of projects.
func (u *ProjectUseCase) ListProjects(ctx context.Context, opts *redmine.ListProjectsOptions) (*redmine.ProjectsResponse, error) {
	return u.client.ListProjects(ctx, opts)
}

// ShowProject retrieves a single project by ID or identifier.
func (u *ProjectUseCase) ShowProject(ctx context.Context, idOrIdentifier string, opts *redmine.ShowProjectOptions) (*redmine.ProjectResponse, error) {
	return u.client.ShowProject(ctx, idOrIdentifier, opts)
}

// CreateProject creates a new project.
func (u *ProjectUseCase) CreateProject(ctx context.Context, project redmine.Project) (*redmine.ProjectResponse, error) {
	return u.client.CreateProject(ctx, project)
}

// UpdateProject updates an existing project.
func (u *ProjectUseCase) UpdateProject(ctx context.Context, idOrIdentifier string, project redmine.Project) error {
	return u.client.UpdateProject(ctx, idOrIdentifier, project)
}

// DeleteProject deletes a project.
func (u *ProjectUseCase) DeleteProject(ctx context.Context, idOrIdentifier string) error {
	return u.client.DeleteProject(ctx, idOrIdentifier)
}

// ArchiveProject archives a project.
func (u *ProjectUseCase) ArchiveProject(ctx context.Context, idOrIdentifier string) error {
	return u.client.ArchiveProject(ctx, idOrIdentifier)
}

// UnarchiveProject unarchives a project.
func (u *ProjectUseCase) UnarchiveProject(ctx context.Context, idOrIdentifier string) error {
	return u.client.UnarchiveProject(ctx, idOrIdentifier)
}
