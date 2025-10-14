package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var customFieldCmd = &cobra.Command{
	Use:   "custom-field",
	Short: "Manage Redmine custom fields",
	Long:  `カスタムフィールドの取得などの操作を行います。`,
}

var customFieldListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all custom fields",
	Long:  `すべてのカスタムフィールド定義をリスト表示します（管理者権限が必要）。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.ListCustomFields(context.Background())
		if err != nil {
			return fmt.Errorf("カスタムフィールドの取得に失敗しました: %w", err)
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
	rootCmd.AddCommand(customFieldCmd)

	// Subcommands
	customFieldCmd.AddCommand(customFieldListCmd)
}
