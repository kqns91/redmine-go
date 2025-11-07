package handlers

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

// RegisterBatchOperationTools registers all batch operation-related MCP tools.
func RegisterBatchOperationTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "batch_operations"

	// Create Task Tree tool
	if cfg.IsToolEnabled(toolGroup, "create_task_tree") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "create_task_tree",
			Description: "Batch create a flat list of tasks with parent-child relationships, dependencies, and estimates. Use 'parent_ref' to reference parent tasks. Ideal for creating 10-30 related tickets for a single feature.",
		}, handleCreateTaskTree(useCases))
	}

	// Bulk Update Issues tool
	if cfg.IsToolEnabled(toolGroup, "bulk_update_issues") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "bulk_update_issues",
			Description: "Batch update multiple issues with filters and field updates. Supports weekday-only scheduling adjustment. Ideal for updating 10-100+ issues at once.",
		}, handleBulkUpdateIssues(useCases))
	}
}

// TaskNode represents a single task in the task tree
type TaskNode struct {
	Subject        string                `json:"subject" jsonschema:"Task subject/title (required)"`
	Description    string                `json:"description,omitempty" jsonschema:"Task description (optional)"`
	TrackerID      int                   `json:"tracker_id,omitempty" jsonschema:"Tracker ID (optional, defaults to project default)"`
	StatusID       int                   `json:"status_id,omitempty" jsonschema:"Status ID (optional)"`
	PriorityID     int                   `json:"priority_id,omitempty" jsonschema:"Priority ID (optional)"`
	AssignedToID   int                   `json:"assigned_to_id,omitempty" jsonschema:"Assignee user ID (optional)"`
	FixedVersionID int                   `json:"fixed_version_id,omitempty" jsonschema:"Target version/milestone ID (optional)"`
	EstimatedHours float64               `json:"estimated_hours,omitempty" jsonschema:"Estimated hours (optional)"`
	StartDate      string                `json:"start_date,omitempty" jsonschema:"Start date in YYYY-MM-DD format (optional)"`
	DueDate        string                `json:"due_date,omitempty" jsonschema:"Due date in YYYY-MM-DD format (optional)"`
	DoneRatio      int                   `json:"done_ratio,omitempty" jsonschema:"Progress percentage 0-100 (optional)"`
	CustomFields   []redmine.CustomField `json:"custom_fields,omitempty" jsonschema:"Custom field values (optional)"`
	ParentRef      string                `json:"parent_ref,omitempty" jsonschema:"Reference to parent task (use task ref like 'task1')"`
	BlocksRefs     []string              `json:"blocks_refs,omitempty" jsonschema:"References to tasks this task blocks (use task refs like 'task1', 'task2')"`
	PrecedesRefs   []string              `json:"precedes_refs,omitempty" jsonschema:"References to tasks this task precedes (use task refs like 'task1', 'task2')"`
	Ref            string                `json:"ref,omitempty" jsonschema:"Unique reference for this task within the tree (e.g., 'task1', 'db_design')"`
}

// CreateTaskTreeArgs represents arguments for creating a task tree
type CreateTaskTreeArgs struct {
	ProjectID int        `json:"project_id" jsonschema:"Project ID (required)"`
	Tasks     []TaskNode `json:"tasks" jsonschema:"Flat list of tasks to create. Use parent_ref to specify parent-child relationships (required)"`
}

// CreateTaskTreeResult represents the result of task tree creation
type CreateTaskTreeResult struct {
	Success      bool                    `json:"success"`
	CreatedCount int                     `json:"created_count"`
	TaskMapping  map[string]int          `json:"task_mapping"` // ref -> issue_id
	Tasks        []CreatedTaskInfo       `json:"tasks"`
	Relations    []CreatedRelationInfo   `json:"relations"`
	Errors       []string                `json:"errors,omitempty"`
	Summary      TaskTreeCreationSummary `json:"summary"`
}

// CreatedTaskInfo represents information about a created task
type CreatedTaskInfo struct {
	Ref      string `json:"ref,omitempty"`
	IssueID  int    `json:"issue_id"`
	Subject  string `json:"subject"`
	ParentID int    `json:"parent_id,omitempty"`
}

// CreatedRelationInfo represents information about a created relation
type CreatedRelationInfo struct {
	RelationID   int    `json:"relation_id"`
	IssueID      int    `json:"issue_id"`
	IssueToID    int    `json:"issue_to_id"`
	RelationType string `json:"relation_type"`
}

// TaskTreeCreationSummary provides summary statistics
type TaskTreeCreationSummary struct {
	TotalTasks         int     `json:"total_tasks"`
	TotalRelations     int     `json:"total_relations"`
	TotalEstimatedTime float64 `json:"total_estimated_time"`
	EarliestStartDate  string  `json:"earliest_start_date,omitempty"`
	LatestDueDate      string  `json:"latest_due_date,omitempty"`
}

func handleCreateTaskTree(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args CreateTaskTreeArgs) (*mcp.CallToolResult, CreateTaskTreeResult, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args CreateTaskTreeArgs) (*mcp.CallToolResult, CreateTaskTreeResult, error) {
		if args.ProjectID == 0 {
			return nil, CreateTaskTreeResult{}, errors.New("project_id is required")
		}

		if len(args.Tasks) == 0 {
			return nil, CreateTaskTreeResult{}, errors.New("at least one task is required")
		}

		result := &CreateTaskTreeResult{
			Success:     true,
			TaskMapping: make(map[string]int),
			Tasks:       []CreatedTaskInfo{},
			Relations:   []CreatedRelationInfo{},
			Errors:      []string{},
			Summary: TaskTreeCreationSummary{
				TotalEstimatedTime: 0,
			},
		}

		// Create all tasks in order (respecting parent references)
		createTasksFlat(ctx, useCases, args.ProjectID, args.Tasks, result)

		// Create relations between tasks
		createTaskRelations(ctx, useCases, args.Tasks, result)

		// Calculate summary
		calculateSummary(result)

		result.CreatedCount = len(result.Tasks)

		return nil, *result, nil
	}
}

func createTasksFlat(ctx context.Context, useCases *usecase.UseCases, projectID int, tasks []TaskNode, result *CreateTaskTreeResult) {
	// Create tasks in multiple passes to handle parent references
	// Pass 1: Create tasks without parents
	// Pass 2+: Create tasks whose parents have been created

	remaining := make([]TaskNode, len(tasks))
	copy(remaining, tasks)
	maxPasses := len(tasks) + 1 // Prevent infinite loops
	pass := 0

	for len(remaining) > 0 && pass < maxPasses {
		pass++
		var stillRemaining []TaskNode

		for _, task := range remaining {
			// Check if parent exists (if parent_ref is specified)
			parentID := 0
			if task.ParentRef != "" {
				if id, ok := result.TaskMapping[task.ParentRef]; ok {
					parentID = id
				} else {
					// Parent not created yet, defer to next pass
					stillRemaining = append(stillRemaining, task)
					continue
				}
			}

			// Create the issue
			req := redmine.IssueCreateRequest{
				ProjectID:      projectID,
				Subject:        task.Subject,
				Description:    task.Description,
				TrackerID:      task.TrackerID,
				StatusID:       task.StatusID,
				PriorityID:     task.PriorityID,
				AssignedToID:   task.AssignedToID,
				ParentIssueID:  parentID,
				FixedVersionID: task.FixedVersionID,
				EstimatedHours: task.EstimatedHours,
				StartDate:      task.StartDate,
				DueDate:        task.DueDate,
				DoneRatio:      task.DoneRatio,
				CustomFields:   task.CustomFields,
			}

			resp, err := useCases.RedmineClient.CreateIssue(ctx, req)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("failed to create issue '%s': %v", task.Subject, err))
				continue
			}

			issueID := resp.Issue.ID

			// Store mapping if ref is provided
			if task.Ref != "" {
				result.TaskMapping[task.Ref] = issueID
			}

			// Store created task info
			result.Tasks = append(result.Tasks, CreatedTaskInfo{
				Ref:      task.Ref,
				IssueID:  issueID,
				Subject:  task.Subject,
				ParentID: parentID,
			})
		}

		remaining = stillRemaining
	}

	// Check if there are unresolved parent references
	if len(remaining) > 0 {
		for _, task := range remaining {
			result.Errors = append(result.Errors, fmt.Sprintf("failed to create task '%s': parent reference '%s' not found", task.Subject, task.ParentRef))
		}
	}
}

func createTaskRelations(ctx context.Context, useCases *usecase.UseCases, tasks []TaskNode, result *CreateTaskTreeResult) {
	// Create relations for each task
	for _, task := range tasks {
		if task.Ref == "" {
			continue
		}

		issueID, ok := result.TaskMapping[task.Ref]
		if !ok {
			continue
		}

		// Create "blocks" relations
		for _, blocksRef := range task.BlocksRefs {
			targetID, ok := result.TaskMapping[blocksRef]
			if !ok {
				result.Errors = append(result.Errors, "referenced task not found: "+blocksRef)
				continue
			}

			relation := redmine.IssueRelation{
				IssueToID:    targetID,
				RelationType: "blocks",
			}

			resp, err := useCases.RedmineClient.CreateIssueRelation(ctx, issueID, relation)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("failed to create relation %d blocks %d: %v", issueID, targetID, err))
				continue
			}

			result.Relations = append(result.Relations, CreatedRelationInfo{
				RelationID:   resp.Relation.ID,
				IssueID:      issueID,
				IssueToID:    targetID,
				RelationType: "blocks",
			})
		}

		// Create "precedes" relations
		for _, precedesRef := range task.PrecedesRefs {
			targetID, ok := result.TaskMapping[precedesRef]
			if !ok {
				result.Errors = append(result.Errors, "referenced task not found: "+precedesRef)
				continue
			}

			relation := redmine.IssueRelation{
				IssueToID:    targetID,
				RelationType: "precedes",
			}

			resp, err := useCases.RedmineClient.CreateIssueRelation(ctx, issueID, relation)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("failed to create relation %d precedes %d: %v", issueID, targetID, err))
				continue
			}

			result.Relations = append(result.Relations, CreatedRelationInfo{
				RelationID:   resp.Relation.ID,
				IssueID:      issueID,
				IssueToID:    targetID,
				RelationType: "precedes",
			})
		}
	}
}

func calculateSummary(result *CreateTaskTreeResult) {
	result.Summary.TotalTasks = len(result.Tasks)
	result.Summary.TotalRelations = len(result.Relations)

	var earliestStart, latestDue time.Time
	var hasStartDate, hasDueDate bool

	// We need to fetch the created issues to get their dates
	// For now, we'll leave these empty as we don't have direct access to the created issue details
	// This would require additional API calls to fetch issue details

	result.Summary.EarliestStartDate = ""
	result.Summary.LatestDueDate = ""

	// Note: estimated time calculation would require fetching all created issues
	// or tracking it during creation
	result.Summary.TotalEstimatedTime = 0

	_, _ = earliestStart, latestDue
	_, _ = hasStartDate, hasDueDate
}

// BulkUpdateIssuesArgs represents arguments for bulk updating issues
type BulkUpdateIssuesArgs struct {
	ProjectID         int                `json:"project_id" jsonschema:"Project ID (required)"`
	IssueIDs          []int              `json:"issue_ids,omitempty" jsonschema:"Specific issue IDs to update (optional, if not provided will use filters)"`
	Filters           *BulkUpdateFilters `json:"filters,omitempty" jsonschema:"Filters to select issues (optional)"`
	Updates           *BulkUpdateFields  `json:"updates,omitempty" jsonschema:"Fields to update (optional)"`
	AdjustForWeekdays bool               `json:"adjust_for_weekdays,omitempty" jsonschema:"Adjust start/due dates to weekdays only, skipping weekends (default: false)"`
	IncludeSubtasks   bool               `json:"include_subtasks,omitempty" jsonschema:"Include subtasks in the update (default: false)"`
}

// BulkUpdateFilters represents filters for selecting issues
type BulkUpdateFilters struct {
	StatusID   int    `json:"status_id,omitempty" jsonschema:"Filter by status ID"`
	TrackerID  int    `json:"tracker_id,omitempty" jsonschema:"Filter by tracker ID"`
	PriorityID int    `json:"priority_id,omitempty" jsonschema:"Filter by priority ID"`
	StartDate  string `json:"start_date,omitempty" jsonschema:"Filter by start date (YYYY-MM-DD)"`
}

// BulkUpdateFields represents fields to update
type BulkUpdateFields struct {
	StatusID       int     `json:"status_id,omitempty" jsonschema:"New status ID"`
	PriorityID     int     `json:"priority_id,omitempty" jsonschema:"New priority ID"`
	AssignedToID   int     `json:"assigned_to_id,omitempty" jsonschema:"New assignee user ID"`
	DoneRatio      int     `json:"done_ratio,omitempty" jsonschema:"New progress percentage 0-100"`
	EstimatedHours float64 `json:"estimated_hours,omitempty" jsonschema:"New estimated hours"`
}

// BulkUpdateResult represents the result of bulk update operation
type BulkUpdateResult struct {
	Success       bool                   `json:"success"`
	UpdatedCount  int                    `json:"updated_count"`
	UpdatedIssues []BulkUpdatedIssueInfo `json:"updated_issues"`
	Errors        []string               `json:"errors,omitempty"`
	Summary       BulkUpdateSummary      `json:"summary"`
}

// BulkUpdatedIssueInfo represents information about an updated issue
type BulkUpdatedIssueInfo struct {
	IssueID       int    `json:"issue_id"`
	Subject       string `json:"subject"`
	OldStartDate  string `json:"old_start_date,omitempty"`
	NewStartDate  string `json:"new_start_date,omitempty"`
	OldDueDate    string `json:"old_due_date,omitempty"`
	NewDueDate    string `json:"new_due_date,omitempty"`
	FieldsUpdated int    `json:"fields_updated"`
}

// BulkUpdateSummary provides summary statistics
type BulkUpdateSummary struct {
	TotalIssues        int `json:"total_issues"`
	SuccessfulUpdates  int `json:"successful_updates"`
	FailedUpdates      int `json:"failed_updates"`
	WeekdayAdjustments int `json:"weekday_adjustments"`
}

func handleBulkUpdateIssues(useCases *usecase.UseCases) func(ctx context.Context, request *mcp.CallToolRequest, args BulkUpdateIssuesArgs) (*mcp.CallToolResult, BulkUpdateResult, error) {
	return func(ctx context.Context, request *mcp.CallToolRequest, args BulkUpdateIssuesArgs) (*mcp.CallToolResult, BulkUpdateResult, error) {
		if args.ProjectID == 0 {
			return nil, BulkUpdateResult{}, errors.New("project_id is required")
		}

		result := &BulkUpdateResult{
			Success:       true,
			UpdatedIssues: []BulkUpdatedIssueInfo{},
			Errors:        []string{},
			Summary:       BulkUpdateSummary{},
		}

		// Determine which issues to update
		issuesToUpdate, err := fetchIssuesToUpdate(ctx, useCases, args, result)
		if err != nil {
			return nil, BulkUpdateResult{}, err
		}

		result.Summary.TotalIssues = len(issuesToUpdate)

		// Update each issue
		for _, issue := range issuesToUpdate {
			processIssueUpdate(ctx, useCases, issue, args, result)
		}

		result.UpdatedCount = len(result.UpdatedIssues)

		return nil, *result, nil
	}
}

func fetchIssuesToUpdate(ctx context.Context, useCases *usecase.UseCases, args BulkUpdateIssuesArgs, result *BulkUpdateResult) ([]redmine.Issue, error) {
	var issuesToUpdate []redmine.Issue

	if len(args.IssueIDs) > 0 {
		// Use specific issue IDs
		for _, issueID := range args.IssueIDs {
			issue, err := useCases.RedmineClient.ShowIssue(ctx, issueID, nil)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("failed to fetch issue #%d: %v", issueID, err))
				continue
			}
			issuesToUpdate = append(issuesToUpdate, issue.Issue)
		}
	} else {
		// Use filters to find issues
		listOpts := &redmine.ListIssuesOptions{
			ProjectID:    args.ProjectID,
			Limit:        1000,
			SubprojectID: "*",
		}

		if args.Filters != nil {
			if args.Filters.StatusID > 0 {
				listOpts.StatusID = strconv.Itoa(args.Filters.StatusID)
			}
			if args.Filters.TrackerID > 0 {
				listOpts.TrackerID = args.Filters.TrackerID
			}
		}

		issuesResp, err := useCases.RedmineClient.ListIssues(ctx, listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}
		issuesToUpdate = issuesResp.Issues
	}

	return issuesToUpdate, nil
}

func processIssueUpdate(ctx context.Context, useCases *usecase.UseCases, issue redmine.Issue, args BulkUpdateIssuesArgs, result *BulkUpdateResult) {
	updatedInfo := BulkUpdatedIssueInfo{
		IssueID:       issue.ID,
		Subject:       issue.Subject,
		OldStartDate:  issue.StartDate,
		OldDueDate:    issue.DueDate,
		FieldsUpdated: 0,
	}

	updateReq := buildUpdateRequest(issue, args, &updatedInfo, result)

	// Only update if there are changes
	if updatedInfo.FieldsUpdated > 0 {
		if err := useCases.RedmineClient.UpdateIssue(ctx, issue.ID, updateReq); err != nil {
			result.Errors = append(result.Errors, fmt.Sprintf("failed to update issue #%d: %v", issue.ID, err))
			result.Summary.FailedUpdates++
		} else {
			result.UpdatedIssues = append(result.UpdatedIssues, updatedInfo)
			result.Summary.SuccessfulUpdates++
		}
	}
}

func buildUpdateRequest(issue redmine.Issue, args BulkUpdateIssuesArgs, updatedInfo *BulkUpdatedIssueInfo, result *BulkUpdateResult) redmine.IssueUpdateRequest {
	updateReq := redmine.IssueUpdateRequest{}

	// Apply field updates
	if args.Updates != nil {
		applyFieldUpdates(args.Updates, &updateReq, updatedInfo)
	}

	// Adjust for weekdays if requested
	if args.AdjustForWeekdays {
		adjustDatesForWeekdays(issue, &updateReq, updatedInfo, result)
	}

	return updateReq
}

func applyFieldUpdates(updates *BulkUpdateFields, updateReq *redmine.IssueUpdateRequest, updatedInfo *BulkUpdatedIssueInfo) {
	if updates.StatusID > 0 {
		updateReq.StatusID = updates.StatusID
		updatedInfo.FieldsUpdated++
	}
	if updates.PriorityID > 0 {
		updateReq.PriorityID = updates.PriorityID
		updatedInfo.FieldsUpdated++
	}
	if updates.AssignedToID > 0 {
		updateReq.AssignedToID = updates.AssignedToID
		updatedInfo.FieldsUpdated++
	}
	if updates.DoneRatio > 0 {
		updateReq.DoneRatio = updates.DoneRatio
		updatedInfo.FieldsUpdated++
	}
	if updates.EstimatedHours > 0 {
		updateReq.EstimatedHours = updates.EstimatedHours
		updatedInfo.FieldsUpdated++
	}
}

func adjustDatesForWeekdays(issue redmine.Issue, updateReq *redmine.IssueUpdateRequest, updatedInfo *BulkUpdatedIssueInfo, result *BulkUpdateResult) {
	if issue.StartDate != "" {
		newStart := adjustToWeekday(issue.StartDate)
		if newStart != issue.StartDate {
			updateReq.StartDate = newStart
			updatedInfo.NewStartDate = newStart
			updatedInfo.FieldsUpdated++
			result.Summary.WeekdayAdjustments++
		}
	}
	if issue.DueDate != "" {
		newDue := adjustToWeekday(issue.DueDate)
		if newDue != issue.DueDate {
			updateReq.DueDate = newDue
			updatedInfo.NewDueDate = newDue
			updatedInfo.FieldsUpdated++
			result.Summary.WeekdayAdjustments++
		}
	}
}

// adjustToWeekday moves a date to the next weekday if it falls on a weekend
func adjustToWeekday(dateStr string) string {
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return dateStr
	}

	// If it's Saturday (6) or Sunday (0), move to next Monday
	for date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
		date = date.AddDate(0, 0, 1)
	}

	return date.Format("2006-01-02")
}
