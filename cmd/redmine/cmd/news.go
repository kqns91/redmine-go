package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var newsCmd = &cobra.Command{
	Use:   "news",
	Short: "Manage Redmine news",
	Long:  `ニュースの取得などの操作を行います。`,
}

var newsListCmd = &cobra.Command{
	Use:   "list [project_id_or_identifier]",
	Short: "List news",
	Long:  `ニュースをリスト表示します。プロジェクトIDを指定した場合は、そのプロジェクトのニュースのみを表示します。`,
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		format, _ := cmd.Flags().GetString("format")

		opts := &redmine.ListNewsOptions{
			Limit:  limit,
			Offset: offset,
		}

		var result *redmine.NewsResponse
		var err error

		if len(args) > 0 {
			// Project-specific news
			result, err = client.ListProjectNews(context.Background(), args[0], opts)
		} else {
			// All news
			result, err = client.ListNews(context.Background(), opts)
		}

		if err != nil {
			return fmt.Errorf("ニュースの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatNewsTable(result.News)
		case formatText:
			return formatNewsText(result.News)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

// formatNewsTable formats news in table format.
func formatNewsTable(news []redmine.News) error {
	if len(news) == 0 {
		fmt.Println("ニュースが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Title", "Project", "Author", "Created"}
	rows := make([][]string, 0, len(news))

	for _, n := range news {
		rows = append(rows, []string{
			strconv.Itoa(n.ID),
			formatter.TruncateString(n.Title, 35),
			formatter.TruncateString(n.Project.Name, 20),
			formatter.TruncateString(n.Author.Name, 15),
			n.CreatedOn,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatNewsText formats news in simple text format.
func formatNewsText(news []redmine.News) error {
	if len(news) == 0 {
		fmt.Println("ニュースが見つかりませんでした。")
		return nil
	}

	for _, n := range news {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(n.ID)))
		fmt.Println(formatter.FormatKeyValue("Title", n.Title))
		fmt.Println(formatter.FormatKeyValue("Project", n.Project.Name))
		fmt.Println(formatter.FormatKeyValue("Author", n.Author.Name))
		if n.Summary != "" {
			fmt.Println(formatter.FormatKeyValue("Summary", formatter.TruncateString(n.Summary, 80)))
		}
		if n.Description != "" {
			fmt.Println(formatter.FormatKeyValue("Description", formatter.TruncateString(n.Description, 80)))
		}
		fmt.Println(formatter.FormatKeyValue("Created", n.CreatedOn))
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(newsCmd)

	// Subcommands
	newsCmd.AddCommand(newsListCmd)

	// Flags for list command
	newsListCmd.Flags().Int("limit", 0, "取得する最大件数")
	newsListCmd.Flags().Int("offset", 0, "取得開始位置")
	newsListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")
}
