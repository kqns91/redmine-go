package usecase

// UseCases holds all use case instances.
type UseCases struct {
	Project  *ProjectUseCase
	Issue    *IssueUseCase
	User     *UserUseCase
	Category *CategoryUseCase
	Search   *SearchUseCase
	Metadata *MetadataUseCase
}
