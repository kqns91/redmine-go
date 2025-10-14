package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/cli/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var issueStatusCmd = &cobra.Command{
	Use:   "issue-status",
	Short: "Manage Redmine issue statuses",
	Long:  `チケットステータスの取得などの操作を行います。`,
}

var issueStatusListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all issue statuses",
	Long:  `すべてのチケットステータスをリスト表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")

		result, err := client.ListIssueStatuses(context.Background())
		if err != nil {
			return fmt.Errorf("チケットステータスの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatIssueStatusesTable(result.IssueStatuses)
		case formatText:
			return formatIssueStatusesText(result.IssueStatuses)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

// formatIssueStatusesTable formats issue statuses in table format.
func formatIssueStatusesTable(statuses []redmine.IssueStatus) error {
	if len(statuses) == 0 {
		fmt.Println("チケットステータスが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Name", "Is Closed"}
	rows := make([][]string, 0, len(statuses))

	for _, s := range statuses {
		rows = append(rows, []string{
			strconv.Itoa(s.ID),
			formatter.TruncateString(s.Name, 30),
			strconv.FormatBool(s.IsClosed),
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatIssueStatusesText formats issue statuses in simple text format.
func formatIssueStatusesText(statuses []redmine.IssueStatus) error {
	if len(statuses) == 0 {
		fmt.Println("チケットステータスが見つかりませんでした。")
		return nil
	}

	for _, s := range statuses {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(s.ID)))
		fmt.Println(formatter.FormatKeyValue("Name", s.Name))
		fmt.Println(formatter.FormatKeyValue("Is Closed", strconv.FormatBool(s.IsClosed)))
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(issueStatusCmd)

	// Subcommands
	issueStatusCmd.AddCommand(issueStatusListCmd)

	// Flags for list command
	issueStatusListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")
}
