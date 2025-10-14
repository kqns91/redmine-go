package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// NewsUseCase provides business logic for news operations.
type NewsUseCase struct {
	client *redmine.Client
}

// NewNewsUseCase creates a new NewsUseCase instance.
func NewNewsUseCase(client *redmine.Client) *NewsUseCase {
	return &NewsUseCase{
		client: client,
	}
}

// ListNews retrieves all news across all projects with pagination.
func (u *NewsUseCase) ListNews(ctx context.Context, opts *redmine.ListNewsOptions) (*redmine.NewsResponse, error) {
	return u.client.ListNews(ctx, opts)
}

// ListProjectNews retrieves all news from a specific project with pagination.
func (u *NewsUseCase) ListProjectNews(ctx context.Context, projectIDOrIdentifier string, opts *redmine.ListNewsOptions) (*redmine.NewsResponse, error) {
	return u.client.ListProjectNews(ctx, projectIDOrIdentifier, opts)
}
