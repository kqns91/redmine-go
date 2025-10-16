package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var wikiCmd = &cobra.Command{
	Use:   "wiki",
	Short: "Manage Redmine wiki pages",
	Long:  `Wikiページの作成、取得、更新、削除などの操作を行います。`,
}

var wikiListCmd = &cobra.Command{
	Use:   "list [project_id_or_identifier]",
	Short: "List wiki pages for a project",
	Long:  `指定したプロジェクトのWikiページ一覧を取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")

		result, err := client.ListWikiPages(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("wikiページ一覧の取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatWikiPagesTable(result.WikiPages)
		case formatText:
			return formatWikiPagesText(result.WikiPages)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

var wikiGetCmd = &cobra.Command{
	Use:   "get [project_id_or_identifier] [page_name]",
	Short: "Get a wiki page",
	Long:  `指定したプロジェクトのWikiページを取得します。`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		include, _ := cmd.Flags().GetString("include")
		version, _ := cmd.Flags().GetInt("version")
		format, _ := cmd.Flags().GetString("format")

		opts := &redmine.GetWikiPageOptions{
			Include: include,
			Version: version,
		}

		result, err := client.GetWikiPage(context.Background(), args[0], args[1], opts)
		if err != nil {
			return fmt.Errorf("wikiページの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatWikiPageDetail(&result.WikiPage)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
	},
}

var wikiCreateOrUpdateCmd = &cobra.Command{
	Use:   "create-or-update [project_id_or_identifier] [page_name]",
	Short: "Create or update a wiki page",
	Long:  `Wikiページを作成または更新します。`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		text, _ := cmd.Flags().GetString("text")
		comments, _ := cmd.Flags().GetString("comments")
		version, _ := cmd.Flags().GetInt("version")
		uploadsJSON, _ := cmd.Flags().GetString("uploads")

		if text == "" {
			return errors.New("--text フラグは必須です")
		}

		page := redmine.WikiPageUpdate{
			Text:     text,
			Comments: comments,
			Version:  version,
		}

		// Parse uploads if provided
		if uploadsJSON != "" {
			var uploads []redmine.Upload
			if err := json.Unmarshal([]byte(uploadsJSON), &uploads); err != nil {
				return fmt.Errorf("無効なuploads JSON: %w", err)
			}
			page.Uploads = uploads
		}

		err := client.CreateOrUpdateWikiPage(context.Background(), args[0], args[1], page)
		if err != nil {
			return fmt.Errorf("wikiページの作成/更新に失敗しました: %w", err)
		}

		fmt.Println("Wikiページを作成/更新しました")
		return nil
	},
}

var wikiDeleteCmd = &cobra.Command{
	Use:   "delete [project_id_or_identifier] [page_name]",
	Short: "Delete a wiki page",
	Long:  `Wikiページを削除します。`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := client.DeleteWikiPage(context.Background(), args[0], args[1])
		if err != nil {
			return fmt.Errorf("wikiページの削除に失敗しました: %w", err)
		}

		fmt.Println("Wikiページを削除しました")
		return nil
	},
}

// formatWikiPageDetail formats a single wiki page in detailed text format.
func formatWikiPageDetail(w *redmine.WikiPage) error {
	// Title
	fmt.Println(formatter.FormatTitle("Wiki Page: " + w.Title))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("Title", w.Title))
	if w.Version > 0 {
		fmt.Println(formatter.FormatKeyValue("Version", strconv.Itoa(w.Version)))
	}
	if w.Author.Name != "" {
		fmt.Println(formatter.FormatKeyValue("Author", w.Author.Name))
	}
	if w.Comments != "" {
		fmt.Println(formatter.FormatKeyValue("Comments", w.Comments))
	}

	// Content
	if w.Text != "" {
		fmt.Println()
		fmt.Println(formatter.FormatSection("内容"))
		fmt.Println(formatter.TruncateString(w.Text, 200))
	}

	// Timestamps
	fmt.Println()
	fmt.Println(formatter.FormatSection("タイムスタンプ"))
	if w.CreatedOn != "" {
		fmt.Println(formatter.FormatKeyValue("Created", w.CreatedOn))
	}
	if w.UpdatedOn != "" {
		fmt.Println(formatter.FormatKeyValue("Updated", w.UpdatedOn))
	}

	return nil
}

// formatWikiPagesTable formats wiki pages in table format.
func formatWikiPagesTable(pages []redmine.WikiPageIndex) error {
	if len(pages) == 0 {
		fmt.Println("Wikiページが見つかりませんでした。")
		return nil
	}

	headers := []string{"Title", "Version", "Parent", "Created", "Updated"}
	rows := make([][]string, 0, len(pages))

	for _, p := range pages {
		parent := "-"
		if p.Parent.Name != "" {
			parent = p.Parent.Name
		}

		rows = append(rows, []string{
			formatter.TruncateString(p.Title, 30),
			strconv.Itoa(p.Version),
			formatter.TruncateString(parent, 20),
			p.CreatedOn,
			p.UpdatedOn,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatWikiPagesText formats wiki pages in simple text format.
func formatWikiPagesText(pages []redmine.WikiPageIndex) error {
	if len(pages) == 0 {
		fmt.Println("Wikiページが見つかりませんでした。")
		return nil
	}

	for _, p := range pages {
		fmt.Println(formatter.FormatKeyValue("Title", p.Title))
		fmt.Println(formatter.FormatKeyValue("Version", strconv.Itoa(p.Version)))
		if p.Parent.Name != "" {
			fmt.Println(formatter.FormatKeyValue("Parent", p.Parent.Name))
		}
		fmt.Println(formatter.FormatKeyValue("Created", p.CreatedOn))
		fmt.Println(formatter.FormatKeyValue("Updated", p.UpdatedOn))
		fmt.Println()
	}

	return nil
}

// includeOptionsForWiki returns valid include options for wiki commands
func includeOptionsForWiki() []string {
	return []string{"attachments"}
}

func init() {
	rootCmd.AddCommand(wikiCmd)

	// Subcommands
	wikiCmd.AddCommand(wikiListCmd)
	wikiCmd.AddCommand(wikiGetCmd)
	wikiCmd.AddCommand(wikiCreateOrUpdateCmd)
	wikiCmd.AddCommand(wikiDeleteCmd)

	// Flags for list command
	wikiListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")

	// Flags for get command
	wikiGetCmd.Flags().String("include", "", "追加で取得する情報 (attachments)")
	wikiGetCmd.Flags().Int("version", 0, "取得するバージョン番号")
	wikiGetCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")

	// Register flag completion for get command
	_ = wikiGetCmd.RegisterFlagCompletionFunc("include", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return includeOptionsForWiki(), cobra.ShellCompDirectiveNoFileComp
	})

	// Flags for create-or-update command
	wikiCreateOrUpdateCmd.Flags().String("text", "", "Wikiページの本文 (必須)")
	wikiCreateOrUpdateCmd.Flags().String("comments", "", "更新コメント")
	wikiCreateOrUpdateCmd.Flags().Int("version", 0, "更新するバージョン番号（競合チェック用）")
	wikiCreateOrUpdateCmd.Flags().String("uploads", "", "アップロードファイル情報 (JSON形式, 例: '[{\"token\":\"xxx\",\"filename\":\"file.pdf\"}]')")
}
