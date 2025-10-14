package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// SearchUseCase provides business logic for search operations.
type SearchUseCase struct {
	client *redmine.Client
}

// NewSearchUseCase creates a new SearchUseCase instance.
func NewSearchUseCase(client *redmine.Client) *SearchUseCase {
	return &SearchUseCase{
		client: client,
	}
}

// Search performs a search across Redmine.
func (u *SearchUseCase) Search(ctx context.Context, opts *redmine.SearchOptions) (*redmine.SearchResponse, error) {
	return u.client.Search(ctx, opts)
}
