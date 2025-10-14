package usecase

import (
	"context"

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
func (u *CategoryUseCase) ListIssueCategories(ctx context.Context, projectID string) (*redmine.IssueCategoriesResponse, error) {
	return u.client.ListIssueCategories(ctx, projectID)
}

// ShowIssueCategory retrieves a single issue category by ID.
func (u *CategoryUseCase) ShowIssueCategory(ctx context.Context, id int) (*redmine.IssueCategoryResponse, error) {
	return u.client.ShowIssueCategory(ctx, id)
}

// CreateIssueCategory creates a new issue category.
func (u *CategoryUseCase) CreateIssueCategory(ctx context.Context, projectID string, category redmine.IssueCategory) (*redmine.IssueCategoryResponse, error) {
	return u.client.CreateIssueCategory(ctx, projectID, category)
}

// UpdateIssueCategory updates an existing issue category.
func (u *CategoryUseCase) UpdateIssueCategory(ctx context.Context, id int, category redmine.IssueCategory) error {
	return u.client.UpdateIssueCategory(ctx, id, category)
}

// DeleteIssueCategory deletes an issue category.
func (u *CategoryUseCase) DeleteIssueCategory(ctx context.Context, id int, opts *redmine.DeleteIssueCategoryOptions) error {
	return u.client.DeleteIssueCategory(ctx, id, opts)
}
