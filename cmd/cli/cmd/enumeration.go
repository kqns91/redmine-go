package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

var enumerationCmd = &cobra.Command{
	Use:   "enumeration",
	Short: "Manage Redmine enumerations",
	Long:  `列挙値（チケット優先度、作業分類、文書カテゴリ）の取得などの操作を行います。`,
}

var enumerationListCmd = &cobra.Command{
	Use:   "list [type]",
	Short: "List enumerations by type",
	Long:  `指定したタイプの列挙値をリスト表示します。タイプ: issue-priorities, time-entry-activities, document-categories`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var result interface{}
		var err error

		switch args[0] {
		case "issue-priorities":
			result, err = client.ListIssuePriorities(context.Background())
		case "time-entry-activities":
			result, err = client.ListTimeEntryActivities(context.Background())
		case "document-categories":
			result, err = client.ListDocumentCategories(context.Background())
		default:
			return errors.New("無効なタイプです。使用可能なタイプ: issue-priorities, time-entry-activities, document-categories")
		}

		if err != nil {
			return fmt.Errorf("列挙値の取得に失敗しました: %w", err)
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
	rootCmd.AddCommand(enumerationCmd)

	// Subcommands
	enumerationCmd.AddCommand(enumerationListCmd)
}
