package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var issueCmd = &cobra.Command{
	Use:   "issue",
	Short: "Manage Redmine issues",
	Long:  `チケットの作成、取得、更新、削除などの操作を行います。`,
}

var issueListCmd = &cobra.Command{
	Use:   "list",
	Short: "List issues",
	Long:  `チケットをリスト表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetInt("project-id")
		subprojectID, _ := cmd.Flags().GetString("subproject-id")
		trackerID, _ := cmd.Flags().GetInt("tracker-id")
		statusID, _ := cmd.Flags().GetString("status-id")
		assignedToID, _ := cmd.Flags().GetString("assigned-to-id")
		issueID, _ := cmd.Flags().GetString("issue-id")
		parentID, _ := cmd.Flags().GetInt("parent-id")
		createdOn, _ := cmd.Flags().GetString("created-on")
		updatedOn, _ := cmd.Flags().GetString("updated-on")
		include, _ := cmd.Flags().GetString("include")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		sort, _ := cmd.Flags().GetString("sort")
		format, _ := cmd.Flags().GetString("format")

		opts := &redmine.ListIssuesOptions{
			ProjectID:    projectID,
			SubprojectID: subprojectID,
			TrackerID:    trackerID,
			StatusID:     statusID,
			AssignedToID: assignedToID,
			IssueID:      issueID,
			ParentID:     parentID,
			CreatedOn:    createdOn,
			UpdatedOn:    updatedOn,
			Include:      include,
			Limit:        limit,
			Offset:       offset,
			Sort:         sort,
		}

		result, err := client.ListIssues(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("チケットの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatIssuesTable(result.Issues)
		case formatText:
			return formatIssuesText(result.Issues)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

var issueGetCmd = &cobra.Command{
	Use:   "get [issue_id]",
	Short: "Get an issue by ID",
	Long:  `指定したIDのチケットを取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なissue_id: %w", err)
		}

		include, _ := cmd.Flags().GetString("include")
		format, _ := cmd.Flags().GetString("format")

		opts := &redmine.ShowIssueOptions{
			Include: include,
		}

		result, err := client.ShowIssue(context.Background(), id, opts)
		if err != nil {
			return fmt.Errorf("チケットの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatIssueDetail(&result.Issue)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
	},
}

var issueCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new issue",
	Long:  `新しいチケットを作成します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		projectID, _ := cmd.Flags().GetInt("project-id")
		trackerID, _ := cmd.Flags().GetInt("tracker-id")
		statusID, _ := cmd.Flags().GetInt("status-id")
		priorityID, _ := cmd.Flags().GetInt("priority-id")
		categoryID, _ := cmd.Flags().GetInt("category-id")
		fixedVersionID, _ := cmd.Flags().GetInt("fixed-version-id")
		parentIssueID, _ := cmd.Flags().GetInt("parent-issue-id")
		subject, _ := cmd.Flags().GetString("subject")
		description, _ := cmd.Flags().GetString("description")
		assignedToID, _ := cmd.Flags().GetInt("assigned-to-id")
		startDate, _ := cmd.Flags().GetString("start-date")
		dueDate, _ := cmd.Flags().GetString("due-date")
		doneRatio, _ := cmd.Flags().GetInt("done-ratio")
		estimatedHours, _ := cmd.Flags().GetFloat64("estimated-hours")
		isPrivate, _ := cmd.Flags().GetBool("is-private")
		watcherUserIDsStr, _ := cmd.Flags().GetString("watcher-user-ids")

		if projectID == 0 {
			return errors.New("--project-id フラグは必須です")
		}
		if trackerID == 0 {
			return errors.New("--tracker-id フラグは必須です")
		}
		if subject == "" {
			return errors.New("--subject フラグは必須です")
		}

		req := redmine.IssueCreateRequest{
			ProjectID:      projectID,
			TrackerID:      trackerID,
			Subject:        subject,
			StatusID:       statusID,
			PriorityID:     priorityID,
			CategoryID:     categoryID,
			FixedVersionID: fixedVersionID,
			ParentIssueID:  parentIssueID,
			AssignedToID:   assignedToID,
			Description:    description,
			StartDate:      startDate,
			DueDate:        dueDate,
			DoneRatio:      doneRatio,
			EstimatedHours: estimatedHours,
			IsPrivate:      isPrivate,
		}

		// Parse watcher user IDs if provided
		if watcherUserIDsStr != "" {
			watcherIDs, err := parseIntSlice(watcherUserIDsStr)
			if err != nil {
				return fmt.Errorf("無効なwatcher-user-ids: %w", err)
			}
			req.WatcherUserIDs = watcherIDs
		}

		result, err := client.CreateIssue(context.Background(), req)
		if err != nil {
			return fmt.Errorf("チケットの作成に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var issueUpdateCmd = &cobra.Command{
	Use:   "update [issue_id]",
	Short: "Update an existing issue",
	Long:  `既存のチケットを更新します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なissue_id: %w", err)
		}

		subject, _ := cmd.Flags().GetString("subject")
		description, _ := cmd.Flags().GetString("description")
		statusID, _ := cmd.Flags().GetInt("status-id")
		priorityID, _ := cmd.Flags().GetInt("priority-id")
		categoryID, _ := cmd.Flags().GetInt("category-id")
		fixedVersionID, _ := cmd.Flags().GetInt("fixed-version-id")
		parentIssueID, _ := cmd.Flags().GetInt("parent-issue-id")
		assignedToID, _ := cmd.Flags().GetInt("assigned-to-id")
		startDate, _ := cmd.Flags().GetString("start-date")
		dueDate, _ := cmd.Flags().GetString("due-date")
		doneRatio, _ := cmd.Flags().GetInt("done-ratio")
		estimatedHours, _ := cmd.Flags().GetFloat64("estimated-hours")
		isPrivate, _ := cmd.Flags().GetBool("is-private")
		notes, _ := cmd.Flags().GetString("notes")
		privateNotes, _ := cmd.Flags().GetBool("private-notes")

		req := redmine.IssueUpdateRequest{
			Subject:        subject,
			StatusID:       statusID,
			PriorityID:     priorityID,
			CategoryID:     categoryID,
			FixedVersionID: fixedVersionID,
			ParentIssueID:  parentIssueID,
			AssignedToID:   assignedToID,
			Description:    description,
			StartDate:      startDate,
			DueDate:        dueDate,
			DoneRatio:      doneRatio,
			EstimatedHours: estimatedHours,
			IsPrivate:      isPrivate,
			Notes:          notes,
			PrivateNotes:   privateNotes,
		}

		err = client.UpdateIssue(context.Background(), id, req)
		if err != nil {
			return fmt.Errorf("チケットの更新に失敗しました: %w", err)
		}

		fmt.Println("チケットを更新しました")
		return nil
	},
}

var issueDeleteCmd = &cobra.Command{
	Use:   "delete [issue_id]",
	Short: "Delete an issue",
	Long:  `チケットを削除します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なissue_id: %w", err)
		}

		err = client.DeleteIssue(context.Background(), id)
		if err != nil {
			return fmt.Errorf("チケットの削除に失敗しました: %w", err)
		}

		fmt.Println("チケットを削除しました")
		return nil
	},
}

var issueAddWatcherCmd = &cobra.Command{
	Use:   "add-watcher [issue_id] [user_id]",
	Short: "Add a watcher to an issue",
	Long:  `チケットにウォッチャーを追加します。`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なissue_id: %w", err)
		}

		userID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("無効なuser_id: %w", err)
		}

		err = client.AddWatcher(context.Background(), issueID, userID)
		if err != nil {
			return fmt.Errorf("ウォッチャーの追加に失敗しました: %w", err)
		}

		fmt.Println("ウォッチャーを追加しました")
		return nil
	},
}

var issueRemoveWatcherCmd = &cobra.Command{
	Use:   "remove-watcher [issue_id] [user_id]",
	Short: "Remove a watcher from an issue",
	Long:  `チケットからウォッチャーを削除します。`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なissue_id: %w", err)
		}

		userID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("無効なuser_id: %w", err)
		}

		err = client.RemoveWatcher(context.Background(), issueID, userID)
		if err != nil {
			return fmt.Errorf("ウォッチャーの削除に失敗しました: %w", err)
		}

		fmt.Println("ウォッチャーを削除しました")
		return nil
	},
}

// formatIssuesTable formats issues in table format.
func formatIssuesTable(issues []redmine.Issue) error {
	if len(issues) == 0 {
		fmt.Println("チケットが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Project", "Tracker", "Status", "Priority", "Subject", "Assigned", "Updated"}
	rows := make([][]string, 0, len(issues))

	for _, issue := range issues {
		assignedTo := "-"
		if issue.AssignedTo.Name != "" {
			assignedTo = issue.AssignedTo.Name
		}

		rows = append(rows, []string{
			strconv.Itoa(issue.ID),
			issue.Project.Name,
			issue.Tracker.Name,
			issue.Status.Name,
			issue.Priority.Name,
			formatter.TruncateString(issue.Subject, 40),
			formatter.TruncateString(assignedTo, 15),
			issue.UpdatedOn,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatIssuesText formats issues in simple text format.
func formatIssuesText(issues []redmine.Issue) error {
	if len(issues) == 0 {
		fmt.Println("チケットが見つかりませんでした。")
		return nil
	}

	for _, issue := range issues {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(issue.ID)))
		fmt.Println(formatter.FormatKeyValue("Project", issue.Project.Name))
		fmt.Println(formatter.FormatKeyValue("Tracker", issue.Tracker.Name))
		fmt.Println(formatter.FormatKeyValue("Status", issue.Status.Name))
		fmt.Println(formatter.FormatKeyValue("Priority", issue.Priority.Name))
		fmt.Println(formatter.FormatKeyValue("Subject", issue.Subject))
		if issue.AssignedTo.Name != "" {
			fmt.Println(formatter.FormatKeyValue("Assigned To", issue.AssignedTo.Name))
		}
		if issue.Description != "" {
			fmt.Println(formatter.FormatKeyValue("Description", formatter.TruncateString(issue.Description, 100)))
		}
		fmt.Println(formatter.FormatKeyValue("Updated", issue.UpdatedOn))
		fmt.Println()
	}

	return nil
}

// formatIssueDetail formats a single issue in detailed text format.
func formatIssueDetail(issue *redmine.Issue) error {
	// Title
	fmt.Println(formatter.FormatTitle("Issue #" + strconv.Itoa(issue.ID) + ": " + issue.Subject))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(issue.ID)))
	fmt.Println(formatter.FormatKeyValue("Project", issue.Project.Name))
	fmt.Println(formatter.FormatKeyValue("Tracker", issue.Tracker.Name))
	fmt.Println(formatter.FormatKeyValue("Status", issue.Status.Name))
	fmt.Println(formatter.FormatKeyValue("Priority", issue.Priority.Name))
	fmt.Println(formatter.FormatKeyValue("Subject", issue.Subject))

	if issue.Description != "" {
		fmt.Println()
		fmt.Println(formatter.FormatSection("説明"))
		fmt.Println(issue.Description)
	}

	// Assignment & Dates
	fmt.Println()
	fmt.Println(formatter.FormatSection("担当・期日"))
	if issue.AssignedTo.Name != "" {
		fmt.Println(formatter.FormatKeyValue("Assigned To", issue.AssignedTo.Name))
	} else {
		fmt.Println(formatter.FormatKeyValue("Assigned To", "-"))
	}
	if issue.Author.Name != "" {
		fmt.Println(formatter.FormatKeyValue("Author", issue.Author.Name))
	}
	if issue.StartDate != "" {
		fmt.Println(formatter.FormatKeyValue("Start Date", issue.StartDate))
	}
	if issue.DueDate != "" {
		fmt.Println(formatter.FormatKeyValue("Due Date", issue.DueDate))
	}
	if issue.DoneRatio > 0 {
		fmt.Println(formatter.FormatKeyValue("Done Ratio", strconv.Itoa(issue.DoneRatio)+"%"))
	}
	if issue.EstimatedHours > 0 {
		fmt.Println(formatter.FormatKeyValue("Estimated Hours", fmt.Sprintf("%.2f", issue.EstimatedHours)))
	}

	// Timestamps
	fmt.Println()
	fmt.Println(formatter.FormatSection("作成・更新"))
	fmt.Println(formatter.FormatKeyValue("Created", issue.CreatedOn))
	fmt.Println(formatter.FormatKeyValue("Updated", issue.UpdatedOn))
	if issue.ClosedOn != "" {
		fmt.Println(formatter.FormatKeyValue("Closed", issue.ClosedOn))
	}

	// Journals (if included)
	formatJournals(issue.Journals)

	return nil
}

// formatJournals formats journal entries for display
func formatJournals(journals []redmine.Journal) {
	if len(journals) == 0 {
		return
	}

	fmt.Println()
	fmt.Println(formatter.FormatSection("履歴"))
	for i, journal := range journals {
		if i > 0 {
			fmt.Println()
		}
		fmt.Println(formatter.FormatKeyValue("Journal ID", strconv.Itoa(journal.ID)))
		fmt.Println(formatter.FormatKeyValue("User", journal.User.Name))
		fmt.Println(formatter.FormatKeyValue("Created", journal.CreatedOn))

		if journal.Notes != "" {
			fmt.Println(formatter.FormatKeyValue("Notes", journal.Notes))
		}

		if len(journal.Details) > 0 {
			fmt.Println("  Changes:")
			for _, detail := range journal.Details {
				changeDesc := fmt.Sprintf("    - %s: %s → %s", detail.Name, detail.OldValue, detail.NewValue)
				if detail.Property != "" {
					changeDesc = fmt.Sprintf("    - [%s] %s: %s → %s", detail.Property, detail.Name, detail.OldValue, detail.NewValue)
				}
				fmt.Println(changeDesc)
			}
		}
	}
}

// includeOptionsForIssueList returns valid include options for issue list command
func includeOptionsForIssueList() []string {
	return []string{"attachments", "relations"}
}

// includeOptionsForIssueGet returns valid include options for issue get command
func includeOptionsForIssueGet() []string {
	return []string{"children", "attachments", "relations", "changesets", "journals", "watchers", "allowed_statuses"}
}

// parseIntSlice parses a comma-separated string of integers
func parseIntSlice(s string) ([]int, error) {
	if s == "" {
		return nil, nil
	}
	parts := strings.Split(s, ",")
	result := make([]int, 0, len(parts))
	for _, part := range parts {
		id, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, fmt.Errorf("無効な数値: %s", part)
		}
		result = append(result, id)
	}
	return result, nil
}

//nolint:funlen // Flag definitions are necessarily verbose
func init() {
	rootCmd.AddCommand(issueCmd)

	// Subcommands
	issueCmd.AddCommand(issueListCmd)
	issueCmd.AddCommand(issueGetCmd)
	issueCmd.AddCommand(issueCreateCmd)
	issueCmd.AddCommand(issueUpdateCmd)
	issueCmd.AddCommand(issueDeleteCmd)
	issueCmd.AddCommand(issueAddWatcherCmd)
	issueCmd.AddCommand(issueRemoveWatcherCmd)

	// Flags for list command
	issueListCmd.Flags().Int("project-id", 0, "プロジェクトID")
	issueListCmd.Flags().String("subproject-id", "", "サブプロジェクトID")
	issueListCmd.Flags().Int("tracker-id", 0, "トラッカーID")
	issueListCmd.Flags().String("status-id", "", "ステータスID")
	issueListCmd.Flags().String("assigned-to-id", "", "担当者ID")
	issueListCmd.Flags().String("issue-id", "", "特定のissue IDでフィルター")
	issueListCmd.Flags().Int("parent-id", 0, "親issueでフィルター")
	issueListCmd.Flags().String("created-on", "", "作成日でフィルター (例: >=2024-01-01)")
	issueListCmd.Flags().String("updated-on", "", "更新日でフィルター (例: >=2024-01-01)")
	issueListCmd.Flags().String("include", "", "追加で取得する情報 (attachments, relations)")
	issueListCmd.Flags().Int("limit", 0, "取得する最大件数")
	issueListCmd.Flags().Int("offset", 0, "取得開始位置のオフセット")
	issueListCmd.Flags().String("sort", "", "ソート順")
	issueListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")

	// Register flag completion for list command
	_ = issueListCmd.RegisterFlagCompletionFunc("include", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return includeOptionsForIssueList(), cobra.ShellCompDirectiveNoFileComp
	})

	// Flags for get command
	issueGetCmd.Flags().String("include", "", "追加で取得する情報 (children, attachments, relations, changesets, journals, watchers, allowed_statuses)")
	issueGetCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")

	// Register flag completion for get command
	_ = issueGetCmd.RegisterFlagCompletionFunc("include", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return includeOptionsForIssueGet(), cobra.ShellCompDirectiveNoFileComp
	})

	// Flags for create command
	issueCreateCmd.Flags().Int("project-id", 0, "プロジェクトID (必須)")
	issueCreateCmd.Flags().Int("tracker-id", 0, "トラッカーID (必須)")
	issueCreateCmd.Flags().Int("status-id", 0, "ステータスID")
	issueCreateCmd.Flags().Int("priority-id", 0, "優先度ID")
	issueCreateCmd.Flags().Int("category-id", 0, "カテゴリID")
	issueCreateCmd.Flags().Int("fixed-version-id", 0, "バージョン/マイルストーンID")
	issueCreateCmd.Flags().Int("parent-issue-id", 0, "親issue ID")
	issueCreateCmd.Flags().String("subject", "", "件名 (必須)")
	issueCreateCmd.Flags().String("description", "", "説明")
	issueCreateCmd.Flags().Int("assigned-to-id", 0, "担当者ID")
	issueCreateCmd.Flags().String("start-date", "", "開始日 (YYYY-MM-DD)")
	issueCreateCmd.Flags().String("due-date", "", "期日 (YYYY-MM-DD)")
	issueCreateCmd.Flags().Int("done-ratio", 0, "進捗率 (0-100)")
	issueCreateCmd.Flags().Float64("estimated-hours", 0, "予定工数")
	issueCreateCmd.Flags().Bool("is-private", false, "プライベート設定")
	issueCreateCmd.Flags().String("watcher-user-ids", "", "ウォッチャーのユーザーIDリスト (カンマ区切り, 例: 1,2,3)")

	// Flags for update command
	issueUpdateCmd.Flags().String("subject", "", "件名")
	issueUpdateCmd.Flags().String("description", "", "説明")
	issueUpdateCmd.Flags().Int("status-id", 0, "ステータスID")
	issueUpdateCmd.Flags().Int("priority-id", 0, "優先度ID")
	issueUpdateCmd.Flags().Int("category-id", 0, "カテゴリID")
	issueUpdateCmd.Flags().Int("fixed-version-id", 0, "バージョン/マイルストーンID")
	issueUpdateCmd.Flags().Int("parent-issue-id", 0, "親issue ID")
	issueUpdateCmd.Flags().Int("assigned-to-id", 0, "担当者ID")
	issueUpdateCmd.Flags().String("start-date", "", "開始日 (YYYY-MM-DD)")
	issueUpdateCmd.Flags().String("due-date", "", "期日 (YYYY-MM-DD)")
	issueUpdateCmd.Flags().Int("done-ratio", 0, "進捗率 (0-100)")
	issueUpdateCmd.Flags().Float64("estimated-hours", 0, "予定工数")
	issueUpdateCmd.Flags().Bool("is-private", false, "プライベート設定")
	issueUpdateCmd.Flags().String("notes", "", "更新コメント")
	issueUpdateCmd.Flags().Bool("private-notes", false, "コメントをプライベートにする")
}
