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

var timeEntryCmd = &cobra.Command{
	Use:     "time-entry",
	Aliases: []string{"time"},
	Short:   "Manage Redmine time entries",
	Long:    `作業時間の記録、取得、更新、削除などの操作を行います。`,
}

var timeEntryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List time entries",
	Long:  `作業時間をリスト表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetInt("user-id")
		projectID, _ := cmd.Flags().GetString("project-id")
		spentOn, _ := cmd.Flags().GetString("spent-on")
		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")

		opts := &redmine.ListTimeEntriesOptions{
			UserID:    userID,
			ProjectID: projectID,
			SpentOn:   spentOn,
			From:      from,
			To:        to,
			Limit:     limit,
			Offset:    offset,
		}

		result, err := client.ListTimeEntries(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("作業時間の取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var timeEntryGetCmd = &cobra.Command{
	Use:   "get [time_entry_id]",
	Short: "Get a time entry by ID",
	Long:  `指定したIDの作業時間を取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なtime_entry_id: %w", err)
		}

		result, err := client.ShowTimeEntry(context.Background(), id)
		if err != nil {
			return fmt.Errorf("作業時間の取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var timeEntryCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new time entry",
	Long:  `新しい作業時間を記録します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID, _ := cmd.Flags().GetInt("issue-id")
		projectID, _ := cmd.Flags().GetInt("project-id")
		hours, _ := cmd.Flags().GetFloat64("hours")
		activityID, _ := cmd.Flags().GetInt("activity-id")
		comments, _ := cmd.Flags().GetString("comments")
		spentOn, _ := cmd.Flags().GetString("spent-on")

		if hours <= 0 {
			return errors.New("--hours フラグは必須で、0より大きい値を指定してください")
		}
		if activityID == 0 {
			return errors.New("--activity-id フラグは必須です")
		}
		if issueID == 0 && projectID == 0 {
			return errors.New("--issue-id または --project-id のいずれかは必須です")
		}

		timeEntry := redmine.TimeEntry{
			Hours:    hours,
			Activity: redmine.Resource{ID: activityID},
			Comments: comments,
			SpentOn:  spentOn,
		}

		if issueID > 0 {
			timeEntry.Issue = redmine.Resource{ID: issueID}
		}
		if projectID > 0 {
			timeEntry.Project = redmine.Resource{ID: projectID}
		}

		result, err := client.CreateTimeEntry(context.Background(), timeEntry)
		if err != nil {
			return fmt.Errorf("作業時間の記録に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var timeEntryUpdateCmd = &cobra.Command{
	Use:   "update [time_entry_id]",
	Short: "Update an existing time entry",
	Long:  `既存の作業時間を更新します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なtime_entry_id: %w", err)
		}

		hours, _ := cmd.Flags().GetFloat64("hours")
		activityID, _ := cmd.Flags().GetInt("activity-id")
		comments, _ := cmd.Flags().GetString("comments")
		spentOn, _ := cmd.Flags().GetString("spent-on")

		timeEntry := redmine.TimeEntry{
			Comments: comments,
			SpentOn:  spentOn,
		}

		if hours > 0 {
			timeEntry.Hours = hours
		}
		if activityID > 0 {
			timeEntry.Activity = redmine.Resource{ID: activityID}
		}

		err = client.UpdateTimeEntry(context.Background(), id, timeEntry)
		if err != nil {
			return fmt.Errorf("作業時間の更新に失敗しました: %w", err)
		}

		fmt.Println("作業時間を更新しました")
		return nil
	},
}

var timeEntryDeleteCmd = &cobra.Command{
	Use:   "delete [time_entry_id]",
	Short: "Delete a time entry",
	Long:  `作業時間を削除します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なtime_entry_id: %w", err)
		}

		err = client.DeleteTimeEntry(context.Background(), id)
		if err != nil {
			return fmt.Errorf("作業時間の削除に失敗しました: %w", err)
		}

		fmt.Println("作業時間を削除しました")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(timeEntryCmd)

	// Subcommands
	timeEntryCmd.AddCommand(timeEntryListCmd)
	timeEntryCmd.AddCommand(timeEntryGetCmd)
	timeEntryCmd.AddCommand(timeEntryCreateCmd)
	timeEntryCmd.AddCommand(timeEntryUpdateCmd)
	timeEntryCmd.AddCommand(timeEntryDeleteCmd)

	// Flags for list command
	timeEntryListCmd.Flags().Int("user-id", 0, "ユーザーID")
	timeEntryListCmd.Flags().String("project-id", "", "プロジェクトID")
	timeEntryListCmd.Flags().String("spent-on", "", "作業日 (YYYY-MM-DD)")
	timeEntryListCmd.Flags().String("from", "", "開始日 (YYYY-MM-DD)")
	timeEntryListCmd.Flags().String("to", "", "終了日 (YYYY-MM-DD)")
	timeEntryListCmd.Flags().Int("limit", 0, "取得する最大件数")
	timeEntryListCmd.Flags().Int("offset", 0, "取得開始位置のオフセット")

	// Flags for create command
	timeEntryCreateCmd.Flags().Int("issue-id", 0, "チケットID")
	timeEntryCreateCmd.Flags().Int("project-id", 0, "プロジェクトID")
	timeEntryCreateCmd.Flags().Float64("hours", 0, "作業時間（時間単位、必須）")
	timeEntryCreateCmd.Flags().Int("activity-id", 0, "作業分類ID (必須)")
	timeEntryCreateCmd.Flags().String("comments", "", "コメント")
	timeEntryCreateCmd.Flags().String("spent-on", "", "作業日 (YYYY-MM-DD)")

	// Flags for update command
	timeEntryUpdateCmd.Flags().Float64("hours", 0, "作業時間（時間単位）")
	timeEntryUpdateCmd.Flags().Int("activity-id", 0, "作業分類ID")
	timeEntryUpdateCmd.Flags().String("comments", "", "コメント")
	timeEntryUpdateCmd.Flags().String("spent-on", "", "作業日 (YYYY-MM-DD)")
}
