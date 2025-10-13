package usecase

import (
	"github.com/kqns91/redmine-go/pkg/redmine"
)

// CategoryUseCase provides business logic for issue category operations.
type CategoryUseCase struct {
	client *redmine.Client
}

// NewCategoryUseCase creates a new CategoryUseCase instance.
func NewCategoryUseCase(client *redmine.Client) *CategoryUseCase {
	return &CategoryUseCase{
		client: client,
	}
}

// ListIssueCategories retrieves a list of issue categories for a project.
func (u *CategoryUseCase) ListIssueCategories(projectID string) (*redmine.IssueCategoriesResponse, error) {
	return u.client.ListIssueCategories(projectID)
}

// ShowIssueCategory retrieves a single issue category by ID.
func (u *CategoryUseCase) ShowIssueCategory(id int) (*redmine.IssueCategoryResponse, error) {
	return u.client.ShowIssueCategory(id)
}

// CreateIssueCategory creates a new issue category.
func (u *CategoryUseCase) CreateIssueCategory(projectID string, category redmine.IssueCategory) (*redmine.IssueCategoryResponse, error) {
	return u.client.CreateIssueCategory(projectID, category)
}

// UpdateIssueCategory updates an existing issue category.
func (u *CategoryUseCase) UpdateIssueCategory(id int, category redmine.IssueCategory) error {
	return u.client.UpdateIssueCategory(id, category)
}

// DeleteIssueCategory deletes an issue category.
func (u *CategoryUseCase) DeleteIssueCategory(id int, opts *redmine.DeleteIssueCategoryOptions) error {
	return u.client.DeleteIssueCategory(id, opts)
}
