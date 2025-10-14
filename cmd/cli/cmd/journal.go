package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var journalCmd = &cobra.Command{
	Use:   "journal",
	Short: "Manage Redmine journals",
	Long:  `ジャーナル（チケット履歴）の取得などの操作を行います。`,
}

var journalGetCmd = &cobra.Command{
	Use:   "get [journal_id]",
	Short: "Get a journal by ID",
	Long:  `指定したIDのジャーナルを取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なjournal_id: %w", err)
		}

		result, err := client.ShowJournal(context.Background(), id)
		if err != nil {
			return fmt.Errorf("ジャーナルの取得に失敗しました: %w", err)
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
	rootCmd.AddCommand(journalCmd)

	// Subcommands
	journalCmd.AddCommand(journalGetCmd)
}
