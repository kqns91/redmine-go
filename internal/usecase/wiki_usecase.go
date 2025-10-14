package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// WikiUseCase provides business logic for wiki operations.
type WikiUseCase struct {
	client *redmine.Client
}

// NewWikiUseCase creates a new WikiUseCase instance.
func NewWikiUseCase(client *redmine.Client) *WikiUseCase {
	return &WikiUseCase{
		client: client,
	}
}

// ListWikiPages retrieves wiki pages index for a project.
func (u *WikiUseCase) ListWikiPages(ctx context.Context, projectIDOrIdentifier string) (*redmine.WikiPagesResponse, error) {
	return u.client.ListWikiPages(ctx, projectIDOrIdentifier)
}

// ShowWikiPage retrieves a specific wiki page.
func (u *WikiUseCase) ShowWikiPage(ctx context.Context, projectIDOrIdentifier string, pageName string, opts *redmine.GetWikiPageOptions) (*redmine.WikiPageResponse, error) {
	return u.client.GetWikiPage(ctx, projectIDOrIdentifier, pageName, opts)
}

// CreateOrUpdateWikiPage creates or updates a wiki page.
func (u *WikiUseCase) CreateOrUpdateWikiPage(ctx context.Context, projectIDOrIdentifier string, pageName string, page redmine.WikiPageUpdate) error {
	return u.client.CreateOrUpdateWikiPage(ctx, projectIDOrIdentifier, pageName, page)
}

// DeleteWikiPage deletes a wiki page.
func (u *WikiUseCase) DeleteWikiPage(ctx context.Context, projectIDOrIdentifier string, pageName string) error {
	return u.client.DeleteWikiPage(ctx, projectIDOrIdentifier, pageName)
}
