package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var trackerCmd = &cobra.Command{
	Use:   "tracker",
	Short: "Manage Redmine trackers",
	Long:  `トラッカーの取得などの操作を行います。`,
}

var trackerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all trackers",
	Long:  `すべてのトラッカーをリスト表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.ListTrackers(context.Background())
		if err != nil {
			return fmt.Errorf("トラッカーの取得に失敗しました: %w", err)
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
	rootCmd.AddCommand(trackerCmd)

	// Subcommands
	trackerCmd.AddCommand(trackerListCmd)
}
