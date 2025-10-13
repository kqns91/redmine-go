package usecase

import (
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
func (u *UserUseCase) ListUsers(opts *redmine.ListUsersOptions) (*redmine.UsersResponse, error) {
	return u.client.ListUsers(opts)
}

// ShowUser retrieves a single user by ID.
func (u *UserUseCase) ShowUser(id int, opts *redmine.ShowUserOptions) (*redmine.UserResponse, error) {
	return u.client.ShowUser(id, opts)
}

// GetCurrentUser retrieves the current user.
func (u *UserUseCase) GetCurrentUser(opts *redmine.ShowUserOptions) (*redmine.UserResponse, error) {
	return u.client.GetCurrentUser(opts)
}

// CreateUser creates a new user.
func (u *UserUseCase) CreateUser(user redmine.User) (*redmine.UserResponse, error) {
	return u.client.CreateUser(user)
}

// UpdateUser updates an existing user.
func (u *UserUseCase) UpdateUser(id int, user redmine.User) error {
	return u.client.UpdateUser(id, user)
}

// DeleteUser deletes a user.
func (u *UserUseCase) DeleteUser(id int) error {
	return u.client.DeleteUser(id)
}
