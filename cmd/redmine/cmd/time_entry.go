package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
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
		format, _ := cmd.Flags().GetString("format")

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

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatTimeEntriesTable(result.TimeEntries)
		case formatText:
			return formatTimeEntriesText(result.TimeEntries)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
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
		format, _ := cmd.Flags().GetString("format")

		result, err := client.ShowTimeEntry(context.Background(), id)
		if err != nil {
			return fmt.Errorf("作業時間の取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatTimeEntryDetail(&result.TimeEntry)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
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
		customFieldsJSON, _ := cmd.Flags().GetString("custom-fields")

		if hours <= 0 {
			return errors.New("--hours フラグは必須で、0より大きい値を指定してください")
		}
		if activityID == 0 {
			return errors.New("--activity-id フラグは必須です")
		}
		if issueID == 0 && projectID == 0 {
			return errors.New("--issue-id または --project-id のいずれかは必須です")
		}

		customFields, err := parseCustomFieldsForTimeEntry(customFieldsJSON)
		if err != nil {
			return fmt.Errorf("カスタムフィールドのパースに失敗しました: %w", err)
		}

		req := redmine.TimeEntryCreateRequest{
			Hours:        hours,
			ActivityID:   activityID,
			Comments:     comments,
			SpentOn:      spentOn,
			CustomFields: customFields,
		}

		if issueID > 0 {
			req.IssueID = issueID
		}
		if projectID > 0 {
			req.ProjectID = projectID
		}

		result, err := client.CreateTimeEntry(context.Background(), req)
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
		customFieldsJSON, _ := cmd.Flags().GetString("custom-fields")

		customFields, err := parseCustomFieldsForTimeEntry(customFieldsJSON)
		if err != nil {
			return fmt.Errorf("カスタムフィールドのパースに失敗しました: %w", err)
		}

		req := redmine.TimeEntryUpdateRequest{
			Comments:     comments,
			SpentOn:      spentOn,
			CustomFields: customFields,
		}

		if hours > 0 {
			req.Hours = hours
		}
		if activityID > 0 {
			req.ActivityID = activityID
		}

		err = client.UpdateTimeEntry(context.Background(), id, req)
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

// formatTimeEntryDetail formats a single time entry in detailed text format.
func formatTimeEntryDetail(t *redmine.TimeEntry) error {
	// Title
	fmt.Println(formatter.FormatTitle(fmt.Sprintf("Time Entry #%d", t.ID)))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(t.ID)))
	if t.User.Name != "" {
		fmt.Println(formatter.FormatKeyValue("User", t.User.Name))
	}
	if t.Project.Name != "" {
		fmt.Println(formatter.FormatKeyValue("Project", t.Project.Name))
	}
	if t.Issue.ID != 0 {
		fmt.Println(formatter.FormatKeyValue("Issue", fmt.Sprintf("#%d", t.Issue.ID)))
	}
	if t.Activity.Name != "" {
		fmt.Println(formatter.FormatKeyValue("Activity", t.Activity.Name))
	}
	fmt.Println(formatter.FormatKeyValue("Hours", fmt.Sprintf("%.2f", t.Hours)))
	fmt.Println(formatter.FormatKeyValue("Spent On", t.SpentOn))

	// Comments
	if t.Comments != "" {
		fmt.Println()
		fmt.Println(formatter.FormatSection("コメント"))
		fmt.Println(t.Comments)
	}

	// Timestamps
	fmt.Println()
	fmt.Println(formatter.FormatSection("タイムスタンプ"))
	fmt.Println(formatter.FormatKeyValue("Created", t.CreatedOn))
	fmt.Println(formatter.FormatKeyValue("Updated", t.UpdatedOn))

	return nil
}

// formatTimeEntriesTable formats time entries in table format.
func formatTimeEntriesTable(entries []redmine.TimeEntry) error {
	if len(entries) == 0 {
		fmt.Println("作業時間が見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "User", "Project", "Issue", "Activity", "Hours", "Spent On"}
	rows := make([][]string, 0, len(entries))

	for _, t := range entries {
		issueStr := "-"
		if t.Issue.ID != 0 {
			issueStr = fmt.Sprintf("#%d", t.Issue.ID)
		}

		rows = append(rows, []string{
			strconv.Itoa(t.ID),
			formatter.TruncateString(t.User.Name, 15),
			formatter.TruncateString(t.Project.Name, 20),
			issueStr,
			formatter.TruncateString(t.Activity.Name, 15),
			fmt.Sprintf("%.2f", t.Hours),
			t.SpentOn,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatTimeEntriesText formats time entries in simple text format.
func formatTimeEntriesText(entries []redmine.TimeEntry) error {
	if len(entries) == 0 {
		fmt.Println("作業時間が見つかりませんでした。")
		return nil
	}

	for _, t := range entries {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(t.ID)))
		fmt.Println(formatter.FormatKeyValue("User", t.User.Name))
		fmt.Println(formatter.FormatKeyValue("Project", t.Project.Name))
		if t.Issue.ID != 0 {
			fmt.Println(formatter.FormatKeyValue("Issue", fmt.Sprintf("#%d", t.Issue.ID)))
		}
		fmt.Println(formatter.FormatKeyValue("Activity", t.Activity.Name))
		fmt.Println(formatter.FormatKeyValue("Hours", fmt.Sprintf("%.2f", t.Hours)))
		fmt.Println(formatter.FormatKeyValue("Spent On", t.SpentOn))
		if t.Comments != "" {
			fmt.Println(formatter.FormatKeyValue("Comments", formatter.TruncateString(t.Comments, 80)))
		}
		fmt.Println()
	}

	return nil
}

// parseCustomFieldsForTimeEntry parses custom fields from JSON string to []CustomField
func parseCustomFieldsForTimeEntry(s string) ([]redmine.CustomField, error) {
	if s == "" {
		return nil, nil
	}
	var result []redmine.CustomField
	if err := json.Unmarshal([]byte(s), &result); err != nil {
		return nil, fmt.Errorf("無効なJSON: %w", err)
	}
	return result, nil
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
	timeEntryListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")

	// Flags for get command
	timeEntryGetCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")

	// Flags for create command
	timeEntryCreateCmd.Flags().Int("issue-id", 0, "チケットID")
	timeEntryCreateCmd.Flags().Int("project-id", 0, "プロジェクトID")
	timeEntryCreateCmd.Flags().Float64("hours", 0, "作業時間（時間単位、必須）")
	timeEntryCreateCmd.Flags().Int("activity-id", 0, "作業分類ID (必須)")
	timeEntryCreateCmd.Flags().String("comments", "", "コメント")
	timeEntryCreateCmd.Flags().String("spent-on", "", "作業日 (YYYY-MM-DD)")
	timeEntryCreateCmd.Flags().String("custom-fields", "", "カスタムフィールド (JSON形式, 例: '[{\"id\":1,\"value\":\"foo\"}]')")

	// Flags for update command
	timeEntryUpdateCmd.Flags().Float64("hours", 0, "作業時間（時間単位）")
	timeEntryUpdateCmd.Flags().Int("activity-id", 0, "作業分類ID")
	timeEntryUpdateCmd.Flags().String("comments", "", "コメント")
	timeEntryUpdateCmd.Flags().String("spent-on", "", "作業日 (YYYY-MM-DD)")
	timeEntryUpdateCmd.Flags().String("custom-fields", "", "カスタムフィールド (JSON形式, 例: '[{\"id\":1,\"value\":\"foo\"}]')")
}
