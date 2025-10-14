package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

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

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
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
}
