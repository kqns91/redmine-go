package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// MyAccountUseCase provides business logic for my account operations.
type MyAccountUseCase struct {
	client *redmine.Client
}

// NewMyAccountUseCase creates a new MyAccountUseCase instance.
func NewMyAccountUseCase(client *redmine.Client) *MyAccountUseCase {
	return &MyAccountUseCase{
		client: client,
	}
}

// ShowMyAccount retrieves current user's account details.
func (u *MyAccountUseCase) ShowMyAccount(ctx context.Context) (*redmine.MyAccountResponse, error) {
	return u.client.GetMyAccount(ctx)
}

// UpdateMyAccount updates current user's account details.
func (u *MyAccountUseCase) UpdateMyAccount(ctx context.Context, user redmine.User) error {
	return u.client.UpdateMyAccount(ctx, user)
}
