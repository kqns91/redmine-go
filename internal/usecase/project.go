package usecase

import (
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
func (u *ProjectUseCase) ListProjects(opts *redmine.ListProjectsOptions) (*redmine.ProjectsResponse, error) {
	return u.client.ListProjects(opts)
}

// ShowProject retrieves a single project by ID or identifier.
func (u *ProjectUseCase) ShowProject(idOrIdentifier string, opts *redmine.ShowProjectOptions) (*redmine.ProjectResponse, error) {
	return u.client.ShowProject(idOrIdentifier, opts)
}

// CreateProject creates a new project.
func (u *ProjectUseCase) CreateProject(project redmine.Project) (*redmine.ProjectResponse, error) {
	return u.client.CreateProject(project)
}

// UpdateProject updates an existing project.
func (u *ProjectUseCase) UpdateProject(idOrIdentifier string, project redmine.Project) error {
	return u.client.UpdateProject(idOrIdentifier, project)
}

// DeleteProject deletes a project.
func (u *ProjectUseCase) DeleteProject(idOrIdentifier string) error {
	return u.client.DeleteProject(idOrIdentifier)
}

// ArchiveProject archives a project.
func (u *ProjectUseCase) ArchiveProject(idOrIdentifier string) error {
	return u.client.ArchiveProject(idOrIdentifier)
}

// UnarchiveProject unarchives a project.
func (u *ProjectUseCase) UnarchiveProject(idOrIdentifier string) error {
	return u.client.UnarchiveProject(idOrIdentifier)
}
