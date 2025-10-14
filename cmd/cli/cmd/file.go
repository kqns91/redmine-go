package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Manage Redmine files",
	Long:  `ファイルの取得などの操作を行います。`,
}

var fileListCmd = &cobra.Command{
	Use:   "list [project_id_or_identifier]",
	Short: "List files for a project",
	Long:  `指定したプロジェクトのファイルをリスト表示します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.ListFiles(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("ファイルの取得に失敗しました: %w", err)
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
	rootCmd.AddCommand(fileCmd)

	// Subcommands
	fileCmd.AddCommand(fileListCmd)
}
