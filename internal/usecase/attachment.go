package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// AttachmentUseCase provides business logic for attachment operations.
type AttachmentUseCase struct {
	client *redmine.Client
}

// NewAttachmentUseCase creates a new AttachmentUseCase instance.
func NewAttachmentUseCase(client *redmine.Client) *AttachmentUseCase {
	return &AttachmentUseCase{
		client: client,
	}
}

// ShowAttachment retrieves details of a specific attachment.
func (u *AttachmentUseCase) ShowAttachment(ctx context.Context, id int) (*redmine.AttachmentResponse, error) {
	return u.client.ShowAttachment(ctx, id)
}

// UpdateAttachment updates an existing attachment.
func (u *AttachmentUseCase) UpdateAttachment(ctx context.Context, id int, attachment redmine.Attachment) error {
	return u.client.UpdateAttachment(ctx, id, attachment)
}

// DeleteAttachment deletes a specific attachment.
func (u *AttachmentUseCase) DeleteAttachment(ctx context.Context, id int) error {
	return u.client.DeleteAttachment(ctx, id)
}
