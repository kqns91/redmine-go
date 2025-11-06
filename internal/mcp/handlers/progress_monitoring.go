package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/usecase"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

// RegisterProgressMonitoringTools registers all progress monitoring-related MCP tools.
func RegisterProgressMonitoringTools(server *mcp.Server, useCases *usecase.UseCases, cfg *config.Config) {
	const toolGroup = "progress_monitoring"

	// Analyze Project Health tool
	if cfg.IsToolEnabled(toolGroup, "analyze_project_health") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "analyze_project_health",
			Description: "Analyze project health by checking issue progress, delays, and critical path status. Returns summary statistics and lists of at-risk and delayed issues.",
		}, handleAnalyzeProjectHealth(useCases))
	}

	// Suggest Reschedule tool
	if cfg.IsToolEnabled(toolGroup, "suggest_reschedule") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "suggest_reschedule",
			Description: "Suggest rescheduling for delayed tasks based on dependencies and current progress. Optionally apply the suggested changes automatically.",
		}, handleSuggestReschedule(useCases))
	}

	// Adjust Estimates tool
	if cfg.IsToolEnabled(toolGroup, "adjust_estimates") {
		mcp.AddTool(server, &mcp.Tool{
			Name:        "adjust_estimates",
			Description: "Adjust estimated hours based on actual time entries and current progress. Provides forecasted completion dates.",
		}, handleAdjustEstimates(useCases))
	}
}

// AnalyzeProjectHealthArgs represents arguments for project health analysis
type AnalyzeProjectHealthArgs struct {
	ProjectID      int `json:"project_id" jsonschema:"Project ID (required)"`
	ThresholdDays  int `json:"threshold_days,omitempty" jsonschema:"Number of days to consider a task 'at risk' (default: 0, meaning any delay is flagged)"`
	IncludeSubtasks bool `json:"include_subtasks,omitempty" jsonschema:"Whether to include subtasks in analysis (default: true)"`
}

// ProjectHealthResult represents the result of project health analysis
type ProjectHealthResult struct {
	Summary         ProjectHealthSummary `json:"summary"`
	DelayedIssues   []IssueHealth        `json:"delayed_issues"`
	AtRiskIssues    []IssueHealth        `json:"at_risk_issues"`
	OnTrackIssues   []IssueHealth        `json:"on_track_issues"`
	CriticalPath    []IssueHealth        `json:"critical_path"`
	Recommendations []string             `json:"recommendations"`
}

// ProjectHealthSummary provides overall project statistics
type ProjectHealthSummary struct {
	TotalIssues        int     `json:"total_issues"`
	OnTrack            int     `json:"on_track"`
	AtRisk             int     `json:"at_risk"`
	Delayed            int     `json:"delayed"`
	Completed          int     `json:"completed"`
	AverageDelayDays   float64 `json:"average_delay_days"`
	TotalEstimatedTime float64 `json:"total_estimated_time"`
	TotalSpentTime     float64 `json:"total_spent_time"`
	ProgressPercentage float64 `json:"progress_percentage"`
	EstimatedCompletion string  `json:"estimated_completion,omitempty"`
	ProjectStatus      string  `json:"project_status"` // "on_schedule", "at_risk", "delayed"
}

// IssueHealth represents health information for a single issue
type IssueHealth struct {
	ID               int     `json:"id"`
	Subject          string  `json:"subject"`
	Status           string  `json:"status"`
	StartDate        string  `json:"start_date,omitempty"`
	DueDate          string  `json:"due_date,omitempty"`
	DoneRatio        int     `json:"done_ratio"`
	EstimatedHours   float64 `json:"estimated_hours"`
	SpentHours       float64 `json:"spent_hours,omitempty"`
	DelayDays        int     `json:"delay_days"`
	IsCriticalPath   bool    `json:"is_critical_path"`
	BlockedBy        []int   `json:"blocked_by,omitempty"`
	Blocks           []int   `json:"blocks,omitempty"`
	ImpactLevel      string  `json:"impact_level"` // "critical", "high", "medium", "low"
}

func handleAnalyzeProjectHealth(useCases *usecase.UseCases) func(context.Context, mcp.CallToolRequestParams) (interface{}, error) {
	return func(ctx context.Context, params mcp.CallToolRequestParams) (interface{}, error) {
		var args AnalyzeProjectHealthArgs
		if err := json.Unmarshal([]byte(params.Arguments.(json.RawMessage)), &args); err != nil {
			return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
		}

		if args.ProjectID == 0 {
			return nil, fmt.Errorf("project_id is required")
		}

		// Set defaults
		if args.ThresholdDays == 0 {
			args.ThresholdDays = 0
		}

		// Fetch all issues for the project
		listOpts := &redmine.ListIssuesOptions{
			ProjectID: fmt.Sprintf("%d", args.ProjectID),
			Limit:     1000, // Fetch a large number
		}

		issuesResp, err := useCases.RedmineClient.ListIssues(ctx, listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}

		result := analyzeIssues(issuesResp.Issues, args.ThresholdDays)

		return result, nil
	}
}

func analyzeIssues(issues []redmine.Issue, thresholdDays int) *ProjectHealthResult {
	result := &ProjectHealthResult{
		Summary: ProjectHealthSummary{
			TotalIssues: len(issues),
		},
		DelayedIssues:   []IssueHealth{},
		AtRiskIssues:    []IssueHealth{},
		OnTrackIssues:   []IssueHealth{},
		CriticalPath:    []IssueHealth{},
		Recommendations: []string{},
	}

	now := time.Now()
	totalDelay := 0
	delayCount := 0

	for _, issue := range issues {
		health := analyzeIssueHealth(issue, now, thresholdDays)

		// Categorize
		if health.DoneRatio == 100 {
			result.Summary.Completed++
		} else if health.DelayDays > thresholdDays {
			result.DelayedIssues = append(result.DelayedIssues, health)
			totalDelay += health.DelayDays
			delayCount++
		} else if health.DelayDays > 0 {
			result.AtRiskIssues = append(result.AtRiskIssues, health)
		} else {
			result.OnTrackIssues = append(result.OnTrackIssues, health)
		}

		// Track critical path (issues with blocks relations)
		if len(health.Blocks) > 0 {
			result.CriticalPath = append(result.CriticalPath, health)
		}

		// Aggregate time
		result.Summary.TotalEstimatedTime += health.EstimatedHours
		result.Summary.TotalSpentTime += health.SpentHours
	}

	// Calculate averages and percentages
	result.Summary.OnTrack = len(result.OnTrackIssues)
	result.Summary.AtRisk = len(result.AtRiskIssues)
	result.Summary.Delayed = len(result.DelayedIssues)

	if delayCount > 0 {
		result.Summary.AverageDelayDays = float64(totalDelay) / float64(delayCount)
	}

	if result.Summary.TotalEstimatedTime > 0 {
		result.Summary.ProgressPercentage = (result.Summary.TotalSpentTime / result.Summary.TotalEstimatedTime) * 100
	}

	// Determine project status
	if result.Summary.Delayed > 0 || result.Summary.AverageDelayDays > 5 {
		result.Summary.ProjectStatus = "delayed"
	} else if result.Summary.AtRisk > 0 || result.Summary.AverageDelayDays > 2 {
		result.Summary.ProjectStatus = "at_risk"
	} else {
		result.Summary.ProjectStatus = "on_schedule"
	}

	// Generate recommendations
	result.Recommendations = generateRecommendations(result)

	// Sort issues by delay (most delayed first)
	sort.Slice(result.DelayedIssues, func(i, j int) bool {
		return result.DelayedIssues[i].DelayDays > result.DelayedIssues[j].DelayDays
	})

	return result
}

func analyzeIssueHealth(issue redmine.Issue, now time.Time, thresholdDays int) IssueHealth {
	health := IssueHealth{
		ID:             issue.ID,
		Subject:        issue.Subject,
		Status:         issue.Status.Name,
		StartDate:      issue.StartDate,
		DueDate:        issue.DueDate,
		DoneRatio:      issue.DoneRatio,
		EstimatedHours: issue.EstimatedHours,
		SpentHours:     issue.SpentHours,
		DelayDays:      0,
		IsCriticalPath: false,
		BlockedBy:      []int{},
		Blocks:         []int{},
		ImpactLevel:    "low",
	}

	// Calculate delay
	if issue.DueDate != "" && issue.DoneRatio < 100 {
		dueDate, err := time.Parse("2006-01-02", issue.DueDate)
		if err == nil {
			delay := int(now.Sub(dueDate).Hours() / 24)
			if delay > 0 {
				health.DelayDays = delay
			}
		}
	}

	// Analyze relations
	for _, rel := range issue.Relations {
		if rel.RelationType == "blocks" {
			health.Blocks = append(health.Blocks, rel.IssueToID)
			health.IsCriticalPath = true
		} else if rel.RelationType == "blocked" {
			health.BlockedBy = append(health.BlockedBy, rel.IssueID)
		}
	}

	// Determine impact level
	if health.IsCriticalPath && health.DelayDays > thresholdDays {
		health.ImpactLevel = "critical"
	} else if health.DelayDays > thresholdDays+3 {
		health.ImpactLevel = "high"
	} else if health.DelayDays > thresholdDays {
		health.ImpactLevel = "medium"
	}

	return health
}

func generateRecommendations(result *ProjectHealthResult) []string {
	recommendations := []string{}

	if result.Summary.Delayed > 5 {
		recommendations = append(recommendations, fmt.Sprintf("⚠️ %d issues are delayed. Consider rescheduling or reassigning resources.", result.Summary.Delayed))
	}

	if len(result.CriticalPath) > 0 {
		delayedCritical := 0
		for _, issue := range result.CriticalPath {
			if issue.DelayDays > 0 {
				delayedCritical++
			}
		}
		if delayedCritical > 0 {
			recommendations = append(recommendations, fmt.Sprintf("🚨 %d critical path issues are delayed, which will cascade delays to dependent tasks.", delayedCritical))
		}
	}

	if result.Summary.TotalSpentTime > result.Summary.TotalEstimatedTime*1.2 {
		recommendations = append(recommendations, "⏱️ Actual time spent exceeds estimates by 20%+. Consider revising estimates for remaining tasks.")
	}

	if result.Summary.ProjectStatus == "on_schedule" {
		recommendations = append(recommendations, "✅ Project is on schedule. Continue monitoring progress regularly.")
	}

	return recommendations
}

// SuggestRescheduleArgs represents arguments for rescheduling suggestions
type SuggestRescheduleArgs struct {
	ProjectID       int  `json:"project_id" jsonschema:"Project ID (required)"`
	AutoApply       bool `json:"auto_apply,omitempty" jsonschema:"Automatically apply suggested reschedules (default: false)"`
	BufferDays      int  `json:"buffer_days,omitempty" jsonschema:"Number of buffer days to add to rescheduled tasks (default: 2)"`
	OnlyCriticalPath bool `json:"only_critical_path,omitempty" jsonschema:"Only reschedule critical path issues (default: false)"`
}

// RescheduleResult represents the result of reschedule suggestions
type RescheduleResult struct {
	Recommendations      []RescheduleRecommendation `json:"recommendations"`
	TotalAffectedIssues  int                        `json:"total_affected_issues"`
	NewProjectCompletion string                     `json:"new_project_completion,omitempty"`
	Applied              bool                       `json:"applied"`
	Errors               []string                   `json:"errors,omitempty"`
}

// RescheduleRecommendation represents a single reschedule recommendation
type RescheduleRecommendation struct {
	IssueID            int      `json:"issue_id"`
	Subject            string   `json:"subject"`
	CurrentDueDate     string   `json:"current_due_date"`
	RecommendedDueDate string   `json:"recommended_due_date"`
	Reason             string   `json:"reason"`
	CascadeImpact      []int    `json:"cascade_impact,omitempty"` // IDs of issues affected
}

func handleSuggestReschedule(useCases *usecase.UseCases) func(context.Context, mcp.CallToolRequestParams) (interface{}, error) {
	return func(ctx context.Context, params mcp.CallToolRequestParams) (interface{}, error) {
		var args SuggestRescheduleArgs
		if err := json.Unmarshal([]byte(params.Arguments.(json.RawMessage)), &args); err != nil {
			return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
		}

		if args.ProjectID == 0 {
			return nil, fmt.Errorf("project_id is required")
		}

		// Set defaults
		if args.BufferDays == 0 {
			args.BufferDays = 2
		}

		// Fetch all issues
		listOpts := &redmine.ListIssuesOptions{
			ProjectID: fmt.Sprintf("%d", args.ProjectID),
			Limit:     1000,
		}

		issuesResp, err := useCases.RedmineClient.ListIssues(ctx, listOpts)
		if err != nil {
			return nil, fmt.Errorf("failed to list issues: %w", err)
		}

		result := generateRescheduleRecommendations(issuesResp.Issues, args.BufferDays, args.OnlyCriticalPath)

		// Apply if requested
		if args.AutoApply {
			result.Applied = true
			for _, rec := range result.Recommendations {
				req := redmine.IssueUpdateRequest{
					DueDate: rec.RecommendedDueDate,
				}
				if err := useCases.RedmineClient.UpdateIssue(ctx, rec.IssueID, req); err != nil {
					result.Errors = append(result.Errors, fmt.Sprintf("Failed to update issue #%d: %v", rec.IssueID, err))
				}
			}
		}

		return result, nil
	}
}

func generateRescheduleRecommendations(issues []redmine.Issue, bufferDays int, onlyCriticalPath bool) *RescheduleResult {
	result := &RescheduleResult{
		Recommendations: []RescheduleRecommendation{},
		Applied:         false,
		Errors:          []string{},
	}

	now := time.Now()

	for _, issue := range issues {
		if issue.DoneRatio == 100 {
			continue // Skip completed issues
		}

		if issue.DueDate == "" {
			continue // Skip issues without due dates
		}

		dueDate, err := time.Parse("2006-01-02", issue.DueDate)
		if err != nil {
			continue
		}

		delay := int(now.Sub(dueDate).Hours() / 24)
		if delay <= 0 {
			continue // Not delayed
		}

		// Check if on critical path
		isOnCriticalPath := false
		cascadeImpact := []int{}
		for _, rel := range issue.Relations {
			if rel.RelationType == "blocks" {
				isOnCriticalPath = true
				cascadeImpact = append(cascadeImpact, rel.IssueToID)
			}
		}

		if onlyCriticalPath && !isOnCriticalPath {
			continue
		}

		// Calculate new due date
		newDueDate := now.AddDate(0, 0, bufferDays)
		if issue.EstimatedHours > 0 {
			// Add days based on estimated hours (assuming 8h workday)
			remainingDays := int((issue.EstimatedHours * (1 - float64(issue.DoneRatio)/100)) / 8)
			newDueDate = now.AddDate(0, 0, remainingDays+bufferDays)
		}

		rec := RescheduleRecommendation{
			IssueID:            issue.ID,
			Subject:            issue.Subject,
			CurrentDueDate:     issue.DueDate,
			RecommendedDueDate: newDueDate.Format("2006-01-02"),
			Reason:             fmt.Sprintf("Delayed by %d days", delay),
			CascadeImpact:      cascadeImpact,
		}

		if isOnCriticalPath {
			rec.Reason += " (critical path)"
		}

		result.Recommendations = append(result.Recommendations, rec)
	}

	result.TotalAffectedIssues = len(result.Recommendations)

	return result
}

// AdjustEstimatesArgs represents arguments for adjusting estimates
type AdjustEstimatesArgs struct {
	IssueID         int  `json:"issue_id" jsonschema:"Issue ID (required)"`
	IncludeChildren bool `json:"include_children,omitempty" jsonschema:"Include child issues in calculation (default: true)"`
}

// EstimateAdjustmentResult represents the result of estimate adjustment
type EstimateAdjustmentResult struct {
	IssueID                  int     `json:"issue_id"`
	Subject                  string  `json:"subject"`
	OriginalEstimate         float64 `json:"original_estimate"`
	HoursSpent               float64 `json:"hours_spent"`
	DoneRatio                int     `json:"done_ratio"`
	CalculatedRemaining      float64 `json:"calculated_remaining"`
	RecommendedTotalEstimate float64 `json:"recommended_total_estimate"`
	CompletionForecast       string  `json:"completion_forecast,omitempty"`
	EfficiencyRatio          float64 `json:"efficiency_ratio"` // actual vs estimated
	Children                 []EstimateAdjustmentResult `json:"children,omitempty"`
}

func handleAdjustEstimates(useCases *usecase.UseCases) func(context.Context, mcp.CallToolRequestParams) (interface{}, error) {
	return func(ctx context.Context, params mcp.CallToolRequestParams) (interface{}, error) {
		var args AdjustEstimatesArgs
		if err := json.Unmarshal([]byte(params.Arguments.(json.RawMessage)), &args); err != nil {
			return nil, fmt.Errorf("failed to unmarshal arguments: %w", err)
		}

		if args.IssueID == 0 {
			return nil, fmt.Errorf("issue_id is required")
		}

		// Fetch issue details
		issueResp, err := useCases.RedmineClient.ShowIssue(ctx, args.IssueID, nil)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch issue: %w", err)
		}

		result := calculateEstimateAdjustment(issueResp.Issue, args.IncludeChildren)

		return result, nil
	}
}

func calculateEstimateAdjustment(issue redmine.Issue, includeChildren bool) *EstimateAdjustmentResult {
	result := &EstimateAdjustmentResult{
		IssueID:          issue.ID,
		Subject:          issue.Subject,
		OriginalEstimate: issue.EstimatedHours,
		HoursSpent:       issue.SpentHours,
		DoneRatio:        issue.DoneRatio,
		Children:         []EstimateAdjustmentResult{},
	}

	// Calculate remaining work
	if issue.DoneRatio > 0 && issue.DoneRatio < 100 {
		progressRatio := float64(issue.DoneRatio) / 100.0
		if progressRatio > 0 {
			// Estimate total needed based on current efficiency
			result.RecommendedTotalEstimate = issue.SpentHours / progressRatio
			result.CalculatedRemaining = result.RecommendedTotalEstimate - issue.SpentHours
		}
	}

	// Calculate efficiency ratio
	if issue.EstimatedHours > 0 {
		result.EfficiencyRatio = issue.SpentHours / issue.EstimatedHours
	}

	// Forecast completion date (assuming 8h workday, 5 days/week)
	if result.CalculatedRemaining > 0 {
		workdaysRemaining := int(result.CalculatedRemaining / 8)
		weeksRemaining := workdaysRemaining / 5
		extraDays := workdaysRemaining % 5

		completionDate := time.Now().AddDate(0, 0, weeksRemaining*7+extraDays)
		result.CompletionForecast = completionDate.Format("2006-01-02")
	}

	// Process children if requested
	if includeChildren && len(issue.Children) > 0 {
		for _, child := range issue.Children {
			childResult := calculateEstimateAdjustment(child, true)
			result.Children = append(result.Children, *childResult)
		}
	}

	return result
}
