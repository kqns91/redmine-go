package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

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
		result, err := client.ListWikiPages(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("wikiページ一覧の取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
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

		opts := &redmine.GetWikiPageOptions{
			Include: include,
			Version: version,
		}

		result, err := client.GetWikiPage(context.Background(), args[0], args[1], opts)
		if err != nil {
			return fmt.Errorf("wikiページの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
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

		if text == "" {
			return errors.New("--text フラグは必須です")
		}

		page := redmine.WikiPageUpdate{
			Text:     text,
			Comments: comments,
			Version:  version,
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

func init() {
	rootCmd.AddCommand(wikiCmd)

	// Subcommands
	wikiCmd.AddCommand(wikiListCmd)
	wikiCmd.AddCommand(wikiGetCmd)
	wikiCmd.AddCommand(wikiCreateOrUpdateCmd)
	wikiCmd.AddCommand(wikiDeleteCmd)

	// Flags for get command
	wikiGetCmd.Flags().String("include", "", "追加で取得する情報 (例: attachments)")
	wikiGetCmd.Flags().Int("version", 0, "取得するバージョン番号")

	// Flags for create-or-update command
	wikiCreateOrUpdateCmd.Flags().String("text", "", "Wikiページの本文 (必須)")
	wikiCreateOrUpdateCmd.Flags().String("comments", "", "更新コメント")
	wikiCreateOrUpdateCmd.Flags().Int("version", 0, "更新するバージョン番号（競合チェック用）")
}
