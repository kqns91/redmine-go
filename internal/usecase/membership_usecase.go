package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// MembershipUseCase provides business logic for membership operations.
type MembershipUseCase struct {
	client *redmine.Client
}

// NewMembershipUseCase creates a new MembershipUseCase instance.
func NewMembershipUseCase(client *redmine.Client) *MembershipUseCase {
	return &MembershipUseCase{
		client: client,
	}
}

// ListMemberships retrieves paginated list of project memberships.
func (u *MembershipUseCase) ListMemberships(ctx context.Context, projectIDOrIdentifier string) (*redmine.MembershipsResponse, error) {
	return u.client.ListMemberships(ctx, projectIDOrIdentifier)
}

// ShowMembership retrieves specific membership details by ID.
func (u *MembershipUseCase) ShowMembership(ctx context.Context, id int) (*redmine.MembershipResponse, error) {
	return u.client.ShowMembership(ctx, id)
}

// CreateMembership adds a new project member.
func (u *MembershipUseCase) CreateMembership(ctx context.Context, projectIDOrIdentifier string, membership redmine.MembershipCreateUpdate) (*redmine.MembershipResponse, error) {
	return u.client.CreateMembership(ctx, projectIDOrIdentifier, membership)
}

// UpdateMembership updates membership roles.
func (u *MembershipUseCase) UpdateMembership(ctx context.Context, id int, roleIDs []int) error {
	return u.client.UpdateMembership(ctx, id, roleIDs)
}

// DeleteMembership deletes a membership.
func (u *MembershipUseCase) DeleteMembership(ctx context.Context, id int) error {
	return u.client.DeleteMembership(ctx, id)
}
