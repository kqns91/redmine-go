package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
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
		result, err := client.ListIssueStatuses(context.Background())
		if err != nil {
			return fmt.Errorf("チケットステータスの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(issueStatusCmd)

	// Subcommands
	issueStatusCmd.AddCommand(issueStatusListCmd)
}
