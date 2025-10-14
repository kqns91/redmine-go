package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// RoleUseCase provides business logic for role operations.
type RoleUseCase struct {
	client *redmine.Client
}

// NewRoleUseCase creates a new RoleUseCase instance.
func NewRoleUseCase(client *redmine.Client) *RoleUseCase {
	return &RoleUseCase{
		client: client,
	}
}

// ListRoles retrieves the list of all roles.
func (u *RoleUseCase) ListRoles(ctx context.Context) (*redmine.RolesResponse, error) {
	return u.client.ListRoles(ctx)
}

// ShowRole retrieves permissions for a specific role.
func (u *RoleUseCase) ShowRole(ctx context.Context, id int) (*redmine.RoleResponse, error) {
	return u.client.ShowRole(ctx, id)
}
