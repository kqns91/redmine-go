package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// UserUseCase provides business logic for user operations.
type UserUseCase struct {
	client *redmine.Client
}

// NewUserUseCase creates a new UserUseCase instance.
func NewUserUseCase(client *redmine.Client) *UserUseCase {
	return &UserUseCase{
		client: client,
	}
}

// ListUsers retrieves a list of users.
func (u *UserUseCase) ListUsers(ctx context.Context, opts *redmine.ListUsersOptions) (*redmine.UsersResponse, error) {
	return u.client.ListUsers(ctx, opts)
}

// ShowUser retrieves a single user by ID.
func (u *UserUseCase) ShowUser(ctx context.Context, id int, opts *redmine.ShowUserOptions) (*redmine.UserResponse, error) {
	return u.client.ShowUser(ctx, id, opts)
}

// GetCurrentUser retrieves the current user.
func (u *UserUseCase) GetCurrentUser(ctx context.Context, opts *redmine.ShowUserOptions) (*redmine.UserResponse, error) {
	return u.client.GetCurrentUser(ctx, opts)
}

// CreateUser creates a new user.
func (u *UserUseCase) CreateUser(ctx context.Context, user redmine.User) (*redmine.UserResponse, error) {
	return u.client.CreateUser(ctx, user)
}

// UpdateUser updates an existing user.
func (u *UserUseCase) UpdateUser(ctx context.Context, id int, user redmine.User) error {
	return u.client.UpdateUser(ctx, id, user)
}

// DeleteUser deletes a user.
func (u *UserUseCase) DeleteUser(ctx context.Context, id int) error {
	return u.client.DeleteUser(ctx, id)
}
