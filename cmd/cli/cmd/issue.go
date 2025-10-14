package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

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
		include, _ := cmd.Flags().GetString("include")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		sort, _ := cmd.Flags().GetString("sort")

		opts := &redmine.ListIssuesOptions{
			ProjectID:    projectID,
			SubprojectID: subprojectID,
			TrackerID:    trackerID,
			StatusID:     statusID,
			AssignedToID: assignedToID,
			Include:      include,
			Limit:        limit,
			Offset:       offset,
			Sort:         sort,
		}

		result, err := client.ListIssues(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("チケットの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
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

		opts := &redmine.ShowIssueOptions{
			Include: include,
		}

		result, err := client.ShowIssue(context.Background(), id, opts)
		if err != nil {
			return fmt.Errorf("チケットの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
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
		subject, _ := cmd.Flags().GetString("subject")
		description, _ := cmd.Flags().GetString("description")
		assignedToID, _ := cmd.Flags().GetInt("assigned-to-id")
		startDate, _ := cmd.Flags().GetString("start-date")
		dueDate, _ := cmd.Flags().GetString("due-date")

		if projectID == 0 {
			return errors.New("--project-id フラグは必須です")
		}
		if trackerID == 0 {
			return errors.New("--tracker-id フラグは必須です")
		}
		if subject == "" {
			return errors.New("--subject フラグは必須です")
		}

		issue := redmine.Issue{
			Project:     redmine.Resource{ID: projectID},
			Tracker:     redmine.Resource{ID: trackerID},
			Subject:     subject,
			Description: description,
			StartDate:   startDate,
			DueDate:     dueDate,
		}

		if statusID > 0 {
			issue.Status = redmine.Resource{ID: statusID}
		}
		if priorityID > 0 {
			issue.Priority = redmine.Resource{ID: priorityID}
		}
		if assignedToID > 0 {
			issue.AssignedTo = redmine.Resource{ID: assignedToID}
		}

		result, err := client.CreateIssue(context.Background(), issue)
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
		assignedToID, _ := cmd.Flags().GetInt("assigned-to-id")
		startDate, _ := cmd.Flags().GetString("start-date")
		dueDate, _ := cmd.Flags().GetString("due-date")

		issue := redmine.Issue{
			Subject:     subject,
			Description: description,
			StartDate:   startDate,
			DueDate:     dueDate,
		}

		if statusID > 0 {
			issue.Status = redmine.Resource{ID: statusID}
		}
		if priorityID > 0 {
			issue.Priority = redmine.Resource{ID: priorityID}
		}
		if assignedToID > 0 {
			issue.AssignedTo = redmine.Resource{ID: assignedToID}
		}

		err = client.UpdateIssue(context.Background(), id, issue)
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
	issueListCmd.Flags().String("include", "", "追加で取得する情報")
	issueListCmd.Flags().Int("limit", 0, "取得する最大件数")
	issueListCmd.Flags().Int("offset", 0, "取得開始位置のオフセット")
	issueListCmd.Flags().String("sort", "", "ソート順")

	// Flags for get command
	issueGetCmd.Flags().String("include", "", "追加で取得する情報")

	// Flags for create command
	issueCreateCmd.Flags().Int("project-id", 0, "プロジェクトID (必須)")
	issueCreateCmd.Flags().Int("tracker-id", 0, "トラッカーID (必須)")
	issueCreateCmd.Flags().Int("status-id", 0, "ステータスID")
	issueCreateCmd.Flags().Int("priority-id", 0, "優先度ID")
	issueCreateCmd.Flags().String("subject", "", "件名 (必須)")
	issueCreateCmd.Flags().String("description", "", "説明")
	issueCreateCmd.Flags().Int("assigned-to-id", 0, "担当者ID")
	issueCreateCmd.Flags().String("start-date", "", "開始日 (YYYY-MM-DD)")
	issueCreateCmd.Flags().String("due-date", "", "期日 (YYYY-MM-DD)")

	// Flags for update command
	issueUpdateCmd.Flags().String("subject", "", "件名")
	issueUpdateCmd.Flags().String("description", "", "説明")
	issueUpdateCmd.Flags().Int("status-id", 0, "ステータスID")
	issueUpdateCmd.Flags().Int("priority-id", 0, "優先度ID")
	issueUpdateCmd.Flags().Int("assigned-to-id", 0, "担当者ID")
	issueUpdateCmd.Flags().String("start-date", "", "開始日 (YYYY-MM-DD)")
	issueUpdateCmd.Flags().String("due-date", "", "期日 (YYYY-MM-DD)")
}
