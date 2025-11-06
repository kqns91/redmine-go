package handlers

import (
	"context"
	"encoding/json"
	"fmt"
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
			Description: "Batch create a tree of tasks with parent-child relationships, dependencies, and estimates. Ideal for creating 10-30 related tickets for a single feature.",
		}, handleCreateTaskTree(useCases))
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
	Children       []TaskNode            `json:"children,omitempty" jsonschema:"Child tasks (optional)"`
	BlocksRefs     []string              `json:"blocks_refs,omitempty" jsonschema:"References to tasks this task blocks (use task refs like 'task1', 'task2')"`
	PrecedesRefs   []string              `json:"precedes_refs,omitempty" jsonschema:"References to tasks this task precedes (use task refs like 'task1', 'task2')"`
	Ref            string                `json:"ref,omitempty" jsonschema:"Unique reference for this task within the tree (e.g., 'task1', 'db_design')"`
}

// CreateTaskTreeArgs represents arguments for creating a task tree
type CreateTaskTreeArgs struct {
	ProjectID int        `json:"project_id" jsonschema:"Project ID (required)"`
	Tasks     []TaskNode `json:"tasks" jsonschema:"List of top-level tasks to create (required)"`
}

// CreateTaskTreeResult represents the result of task tree creation
type CreateTaskTreeResult struct {
	Success      bool                      `json:"success"`
	CreatedCount int                       `json:"created_count"`
	TaskMapping  map[string]int            `json:"task_mapping"` // ref -> issue_id
	Tasks        []CreatedTaskInfo         `json:"tasks"`
	Relations    []CreatedRelationInfo     `json:"relations"`
	Errors       []string                  `json:"errors,omitempty"`
	Summary      TaskTreeCreationSummary   `json:"summary"`
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

func handleCreateTaskTree(useCases *usecase.UseCases) func(context.Context, mcp.CallToolRequestParams) (interface{}, error) {
	return func(ctx context.Context, params mcp.CallToolRequestParams) (interface{}, error) {
		var args CreateTaskTreeArgs
		if err := json.Unmarshal([]byte(params.Arguments.(json.RawMessage)), &args); err != nil {
			return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
		}

		if args.ProjectID == 0 {
			return nil, fmt.Errorf("project_id is required")
		}

		if len(args.Tasks) == 0 {
			return nil, fmt.Errorf("at least one task is required")
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

		// Create all tasks recursively
		for _, task := range args.Tasks {
			if err := createTaskNodeRecursive(ctx, useCases, args.ProjectID, 0, task, result); err != nil {
				result.Success = false
				result.Errors = append(result.Errors, err.Error())
			}
		}

		// Create relations between tasks
		if err := createTaskRelations(ctx, useCases, args.Tasks, result); err != nil {
			result.Success = false
			result.Errors = append(result.Errors, err.Error())
		}

		// Calculate summary
		calculateSummary(result)

		result.CreatedCount = len(result.Tasks)

		return result, nil
	}
}

func createTaskNodeRecursive(ctx context.Context, useCases *usecase.UseCases, projectID, parentID int, task TaskNode, result *CreateTaskTreeResult) error {
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
		return fmt.Errorf("failed to create issue '%s': %w", task.Subject, err)
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

	// Create children recursively
	for _, child := range task.Children {
		if err := createTaskNodeRecursive(ctx, useCases, projectID, issueID, child, result); err != nil {
			result.Errors = append(result.Errors, err.Error())
		}
	}

	return nil
}

func createTaskRelations(ctx context.Context, useCases *usecase.UseCases, tasks []TaskNode, result *CreateTaskTreeResult) error {
	// Build a flat list of all tasks including children
	var flatTasks []TaskNode
	var collectTasks func([]TaskNode)
	collectTasks = func(nodes []TaskNode) {
		for _, node := range nodes {
			flatTasks = append(flatTasks, node)
			if len(node.Children) > 0 {
				collectTasks(node.Children)
			}
		}
	}
	collectTasks(tasks)

	// Create relations
	for _, task := range flatTasks {
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
				result.Errors = append(result.Errors, fmt.Sprintf("referenced task not found: %s", blocksRef))
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
				result.Errors = append(result.Errors, fmt.Sprintf("referenced task not found: %s", precedesRef))
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

	return nil
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
