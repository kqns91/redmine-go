package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// JournalUseCase provides business logic for journal operations.
type JournalUseCase struct {
	client *redmine.Client
}

// NewJournalUseCase creates a new JournalUseCase instance.
func NewJournalUseCase(client *redmine.Client) *JournalUseCase {
	return &JournalUseCase{
		client: client,
	}
}

// ShowJournal retrieves a specific journal entry by ID.
// Note: Journals are typically accessed through issues with include=journals parameter.
func (u *JournalUseCase) ShowJournal(ctx context.Context, id int) (*redmine.JournalResponse, error) {
	return u.client.ShowJournal(ctx, id)
}
