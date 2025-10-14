package usecase

// UseCases holds all use case instances.
type UseCases struct {
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
}
