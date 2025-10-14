package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// QueryUseCase provides business logic for query operations.
type QueryUseCase struct {
	client *redmine.Client
}

// NewQueryUseCase creates a new QueryUseCase instance.
func NewQueryUseCase(client *redmine.Client) *QueryUseCase {
	return &QueryUseCase{
		client: client,
	}
}

// ListQueries retrieves all custom queries visible by the user.
func (u *QueryUseCase) ListQueries(ctx context.Context) (*redmine.QueriesResponse, error) {
	return u.client.ListQueries(ctx)
}
