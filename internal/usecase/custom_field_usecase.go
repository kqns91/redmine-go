package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// CustomFieldUseCase provides business logic for custom field operations.
type CustomFieldUseCase struct {
	client *redmine.Client
}

// NewCustomFieldUseCase creates a new CustomFieldUseCase instance.
func NewCustomFieldUseCase(client *redmine.Client) *CustomFieldUseCase {
	return &CustomFieldUseCase{
		client: client,
	}
}

// ListCustomFields retrieves all custom fields definitions (requires admin privileges).
func (u *CustomFieldUseCase) ListCustomFields(ctx context.Context) (*redmine.CustomFieldsResponse, error) {
	return u.client.ListCustomFields(ctx)
}
