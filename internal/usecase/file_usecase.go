package usecase

import (
	"context"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

// FileUseCase provides business logic for file operations.
type FileUseCase struct {
	client *redmine.Client
}

// NewFileUseCase creates a new FileUseCase instance.
func NewFileUseCase(client *redmine.Client) *FileUseCase {
	return &FileUseCase{
		client: client,
	}
}

// ListFiles retrieves files available for a specific project.
func (u *FileUseCase) ListFiles(ctx context.Context, projectIDOrIdentifier string) (*redmine.FilesResponse, error) {
	return u.client.ListFiles(ctx, projectIDOrIdentifier)
}

// UploadFile uploads a file to a specific project.
func (u *FileUseCase) UploadFile(ctx context.Context, projectIDOrIdentifier string, fileUpload redmine.FileUpload) error {
	return u.client.UploadFile(ctx, projectIDOrIdentifier, fileUpload)
}
