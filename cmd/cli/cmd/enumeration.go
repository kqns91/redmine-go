package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/cli/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
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
		format, _ := cmd.Flags().GetString("format")
		var result *redmine.EnumerationsResponse
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

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatEnumerationsTable(result.Enumerations, args[0])
		case formatText:
			return formatEnumerationsText(result.Enumerations, args[0])
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

// formatEnumerationsTable formats enumerations in table format.
func formatEnumerationsTable(enums []redmine.Enumeration, _ string) error {
	if len(enums) == 0 {
		fmt.Println("列挙値が見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Name", "Is Default"}
	rows := make([][]string, 0, len(enums))

	for _, e := range enums {
		rows = append(rows, []string{
			strconv.Itoa(e.ID),
			formatter.TruncateString(e.Name, 40),
			strconv.FormatBool(e.IsDefault),
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatEnumerationsText formats enumerations in simple text format.
func formatEnumerationsText(enums []redmine.Enumeration, enumType string) error {
	if len(enums) == 0 {
		fmt.Println("列挙値が見つかりませんでした。")
		return nil
	}

	fmt.Println(formatter.FormatTitle("Enumerations: " + enumType))
	fmt.Println()

	for _, e := range enums {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(e.ID)))
		fmt.Println(formatter.FormatKeyValue("Name", e.Name))
		fmt.Println(formatter.FormatKeyValue("Is Default", strconv.FormatBool(e.IsDefault)))
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(enumerationCmd)

	// Subcommands
	enumerationCmd.AddCommand(enumerationListCmd)

	// Flags for list command
	enumerationListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")
}
