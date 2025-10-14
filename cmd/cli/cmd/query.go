package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "Manage Redmine queries",
	Long:  `カスタムクエリの取得などの操作を行います。`,
}

var queryListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all queries",
	Long:  `ユーザーが閲覧可能なすべてのカスタムクエリをリスト表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.ListQueries(context.Background())
		if err != nil {
			return fmt.Errorf("クエリの取得に失敗しました: %w", err)
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
	rootCmd.AddCommand(queryCmd)

	// Subcommands
	queryCmd.AddCommand(queryListCmd)
}
