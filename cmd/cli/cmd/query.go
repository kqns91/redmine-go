package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/cli/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
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
		format, _ := cmd.Flags().GetString("format")

		result, err := client.ListQueries(context.Background())
		if err != nil {
			return fmt.Errorf("クエリの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatQueriesTable(result.Queries)
		case formatText:
			return formatQueriesText(result.Queries)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

// formatQueriesTable formats queries in table format.
func formatQueriesTable(queries []redmine.Query) error {
	if len(queries) == 0 {
		fmt.Println("クエリが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Name", "Is Public", "Project ID"}
	rows := make([][]string, 0, len(queries))

	for _, q := range queries {
		projectID := "-"
		if q.ProjectID > 0 {
			projectID = strconv.Itoa(q.ProjectID)
		}

		rows = append(rows, []string{
			strconv.Itoa(q.ID),
			formatter.TruncateString(q.Name, 40),
			strconv.FormatBool(q.IsPublic),
			projectID,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatQueriesText formats queries in simple text format.
func formatQueriesText(queries []redmine.Query) error {
	if len(queries) == 0 {
		fmt.Println("クエリが見つかりませんでした。")
		return nil
	}

	for _, q := range queries {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(q.ID)))
		fmt.Println(formatter.FormatKeyValue("Name", q.Name))
		fmt.Println(formatter.FormatKeyValue("Is Public", strconv.FormatBool(q.IsPublic)))
		if q.ProjectID > 0 {
			fmt.Println(formatter.FormatKeyValue("Project ID", strconv.Itoa(q.ProjectID)))
		}
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(queryCmd)

	// Subcommands
	queryCmd.AddCommand(queryListCmd)

	// Flags for list command
	queryListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")
}
