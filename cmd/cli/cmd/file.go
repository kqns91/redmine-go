package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/cli/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
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
		format, _ := cmd.Flags().GetString("format")

		result, err := client.ListFiles(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("ファイルの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatFilesTable(result.Files)
		case formatText:
			return formatFilesText(result.Files)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

// formatFilesTable formats files in table format.
func formatFilesTable(files []redmine.File) error {
	if len(files) == 0 {
		fmt.Println("ファイルが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Filename", "Size", "Downloads", "Created"}
	rows := make([][]string, 0, len(files))

	for _, f := range files {
		rows = append(rows, []string{
			strconv.Itoa(f.ID),
			formatter.TruncateString(f.Filename, 40),
			strconv.Itoa(f.Filesize),
			strconv.Itoa(f.Downloads),
			f.CreatedOn,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatFilesText formats files in simple text format.
func formatFilesText(files []redmine.File) error {
	if len(files) == 0 {
		fmt.Println("ファイルが見つかりませんでした。")
		return nil
	}

	for _, f := range files {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(f.ID)))
		fmt.Println(formatter.FormatKeyValue("Filename", f.Filename))
		fmt.Println(formatter.FormatKeyValue("Size", strconv.Itoa(f.Filesize)))
		fmt.Println(formatter.FormatKeyValue("Content-Type", f.ContentType))
		if f.Description != "" {
			fmt.Println(formatter.FormatKeyValue("Description", formatter.TruncateString(f.Description, 80)))
		}
		fmt.Println(formatter.FormatKeyValue("Downloads", strconv.Itoa(f.Downloads)))
		fmt.Println(formatter.FormatKeyValue("Created", f.CreatedOn))
		if f.Author.Name != "" {
			fmt.Println(formatter.FormatKeyValue("Author", f.Author.Name))
		}
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(fileCmd)

	// Subcommands
	fileCmd.AddCommand(fileListCmd)

	// Flags for list command
	fileListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")
}
