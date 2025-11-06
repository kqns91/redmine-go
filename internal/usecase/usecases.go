package usecase

import "github.com/kqns91/redmine-go/pkg/redmine"

// UseCases holds all use case instances.
type UseCases struct {
	RedmineClient *redmine.Client // Direct access to Redmine client for batch operations
	Project       *ProjectUseCase
	Issue         *IssueUseCase
	User          *UserUseCase
	Category      *CategoryUseCase
	Search        *SearchUseCase
	Metadata      *MetadataUseCase
	TimeEntry     *TimeEntryUseCase
	Version       *VersionUseCase
	IssueRelation *IssueRelationUseCase
	Attachment    *AttachmentUseCase
	Membership    *MembershipUseCase
	Group         *GroupUseCase
	Wiki          *WikiUseCase
	News          *NewsUseCase
	File          *FileUseCase
	Query         *QueryUseCase
	CustomField   *CustomFieldUseCase
	Journal       *JournalUseCase
	Role          *RoleUseCase
	Enumeration   *EnumerationUseCase
	MyAccount     *MyAccountUseCase
}
