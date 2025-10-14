package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// GroupUseCase provides business logic for group operations.
type GroupUseCase struct {
	client *redmine.Client
}

// NewGroupUseCase creates a new GroupUseCase instance.
func NewGroupUseCase(client *redmine.Client) *GroupUseCase {
	return &GroupUseCase{
		client: client,
	}
}

// ListGroups retrieves the list of all groups (admin only).
func (u *GroupUseCase) ListGroups(ctx context.Context, opts *redmine.ListGroupsOptions) (*redmine.GroupsResponse, error) {
	return u.client.ListGroups(ctx, opts)
}

// ShowGroup retrieves group details (admin only).
func (u *GroupUseCase) ShowGroup(ctx context.Context, id int, opts *redmine.ShowGroupOptions) (*redmine.GroupResponse, error) {
	return u.client.ShowGroup(ctx, id, opts)
}

// CreateGroup creates a new group (admin only).
func (u *GroupUseCase) CreateGroup(ctx context.Context, group redmine.Group) (*redmine.GroupResponse, error) {
	return u.client.CreateGroup(ctx, group)
}

// UpdateGroup updates an existing group (admin only).
func (u *GroupUseCase) UpdateGroup(ctx context.Context, id int, group redmine.Group) error {
	return u.client.UpdateGroup(ctx, id, group)
}

// DeleteGroup deletes a group (admin only).
func (u *GroupUseCase) DeleteGroup(ctx context.Context, id int) error {
	return u.client.DeleteGroup(ctx, id)
}

// AddGroupUser adds a user to a group (admin only).
func (u *GroupUseCase) AddGroupUser(ctx context.Context, groupID int, userID int) error {
	return u.client.AddUserToGroup(ctx, groupID, userID)
}

// DeleteGroupUser removes a user from a group (admin only).
func (u *GroupUseCase) DeleteGroupUser(ctx context.Context, groupID int, userID int) error {
	return u.client.RemoveUserFromGroup(ctx, groupID, userID)
}
