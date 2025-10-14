package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/cli/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
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
		format, _ := cmd.Flags().GetString("format")

		result, err := client.ListCustomFields(context.Background())
		if err != nil {
			return fmt.Errorf("カスタムフィールドの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatCustomFieldsTable(result.CustomFields)
		case formatText:
			return formatCustomFieldsText(result.CustomFields)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

// formatCustomFieldsTable formats custom fields in table format.
func formatCustomFieldsTable(fields []redmine.CustomFieldDefinition) error {
	if len(fields) == 0 {
		fmt.Println("カスタムフィールドが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Name", "Type", "Format", "Required", "Visible"}
	rows := make([][]string, 0, len(fields))

	for _, f := range fields {
		rows = append(rows, []string{
			strconv.Itoa(f.ID),
			formatter.TruncateString(f.Name, 25),
			formatter.TruncateString(f.CustomizedType, 15),
			formatter.TruncateString(f.FieldFormat, 10),
			strconv.FormatBool(f.IsRequired),
			strconv.FormatBool(f.Visible),
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatCustomFieldsText formats custom fields in simple text format.
func formatCustomFieldsText(fields []redmine.CustomFieldDefinition) error {
	if len(fields) == 0 {
		fmt.Println("カスタムフィールドが見つかりませんでした。")
		return nil
	}

	for _, f := range fields {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(f.ID)))
		fmt.Println(formatter.FormatKeyValue("Name", f.Name))
		fmt.Println(formatter.FormatKeyValue("Customized Type", f.CustomizedType))
		fmt.Println(formatter.FormatKeyValue("Field Format", f.FieldFormat))
		fmt.Println(formatter.FormatKeyValue("Required", strconv.FormatBool(f.IsRequired)))
		fmt.Println(formatter.FormatKeyValue("Visible", strconv.FormatBool(f.Visible)))
		if f.DefaultValue != "" {
			fmt.Println(formatter.FormatKeyValue("Default Value", f.DefaultValue))
		}
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(customFieldCmd)

	// Subcommands
	customFieldCmd.AddCommand(customFieldListCmd)

	// Flags for list command
	customFieldListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")
}
