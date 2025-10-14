package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

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

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(newsCmd)

	// Subcommands
	newsCmd.AddCommand(newsListCmd)

	// Flags for list command
	newsListCmd.Flags().Int("limit", 0, "取得する最大件数")
	newsListCmd.Flags().Int("offset", 0, "取得開始位置")
}
