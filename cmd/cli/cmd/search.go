package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/cli/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search in Redmine",
	Long:  `Redmine全体で検索を行います。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scope, _ := cmd.Flags().GetString("scope")
		issues, _ := cmd.Flags().GetBool("issues")
		wikiPages, _ := cmd.Flags().GetBool("wiki-pages")
		attachments, _ := cmd.Flags().GetBool("attachments")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		format, _ := cmd.Flags().GetString("format")

		if args[0] == "" {
			return errors.New("検索クエリは必須です")
		}

		opts := &redmine.SearchOptions{
			Query:       strings.Split(args[0], " "),
			Scope:       scope,
			Issues:      issues,
			WikiPages:   wikiPages,
			Attachments: attachments,
			Limit:       limit,
			Offset:      offset,
		}

		result, err := client.Search(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("検索に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatSearchResultsTable(result.Results)
		case formatText:
			return formatSearchResultsText(result.Results)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

// formatSearchResultsTable formats search results in table format.
func formatSearchResultsTable(results []redmine.SearchResult) error {
	if len(results) == 0 {
		fmt.Println("検索結果が見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Type", "Title", "Datetime"}
	rows := make([][]string, 0, len(results))

	for _, r := range results {
		rows = append(rows, []string{
			strconv.Itoa(r.ID),
			r.Type,
			formatter.TruncateString(r.Title, 50),
			r.Datetime,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatSearchResultsText formats search results in simple text format.
func formatSearchResultsText(results []redmine.SearchResult) error {
	if len(results) == 0 {
		fmt.Println("検索結果が見つかりませんでした。")
		return nil
	}

	for _, r := range results {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(r.ID)))
		fmt.Println(formatter.FormatKeyValue("Type", r.Type))
		fmt.Println(formatter.FormatKeyValue("Title", r.Title))
		if r.Description != "" {
			fmt.Println(formatter.FormatKeyValue("Description", formatter.TruncateString(r.Description, 100)))
		}
		if r.URL != "" {
			fmt.Println(formatter.FormatKeyValue("URL", r.URL))
		}
		fmt.Println(formatter.FormatKeyValue("Datetime", r.Datetime))
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(searchCmd)

	// Flags
	searchCmd.Flags().String("scope", "", "検索範囲 (all, my_projects, subprojects)")
	searchCmd.Flags().Bool("issues", false, "チケットを検索対象に含める")
	searchCmd.Flags().Bool("wiki-pages", false, "Wikiページを検索対象に含める")
	searchCmd.Flags().Bool("attachments", false, "添付ファイルを検索対象に含める")
	searchCmd.Flags().Int("limit", 0, "取得する最大件数")
	searchCmd.Flags().Int("offset", 0, "取得開始位置")
	searchCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")
}
