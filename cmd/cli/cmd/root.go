package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

var (
	apiURL string
	apiKey string
	client *redmine.Client
)

// rootCmd はCLIのルートコマンドを表します
var rootCmd = &cobra.Command{
	Use:   "redmine-cli",
	Short: "Redmine API client CLI",
	Long: `redmine-cli は Redmine の REST API を操作するための CLI ツールです。
すべての Redmine API 操作を CLI から実行できます。`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// 環境変数からAPIのURLとキーを取得
		if apiURL == "" {
			apiURL = os.Getenv("REDMINE_API_URL")
		}
		if apiKey == "" {
			apiKey = os.Getenv("REDMINE_API_KEY")
		}

		if apiURL == "" {
			return errors.New("REDMINE_API_URL が設定されていません。--url フラグまたは環境変数を設定してください")
		}
		if apiKey == "" {
			return errors.New("REDMINE_API_KEY が設定されていません。--key フラグまたは環境変数を設定してください")
		}

		// Redmine クライアントを初期化
		client = redmine.New(apiURL, apiKey)
		return nil
	},
}

// Execute はルートコマンドを実行します
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// グローバルフラグの定義
	rootCmd.PersistentFlags().StringVar(&apiURL, "url", "", "Redmine API URL (環境変数: REDMINE_API_URL)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "key", "", "Redmine API Key (環境変数: REDMINE_API_KEY)")
}
