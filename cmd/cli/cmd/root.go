package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cliconfig "github.com/kqns91/redmine-go/cmd/cli/internal/config"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var (
	apiURL       string
	apiKey       string
	outputFormat string
	client       *redmine.Client
)

// rootCmd はCLIのルートコマンドを表します
var rootCmd = &cobra.Command{
	Use:   "redmine-cli",
	Short: "Redmine API client CLI",
	Long: `redmine-cli は Redmine の REST API を操作するための CLI ツールです。
すべての Redmine API 操作を CLI から実行できます。`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// Skip config initialization for config commands
		if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
			return nil
		}
		if cmd.Name() == "config" {
			return nil
		}

		// 優先順位: 1. フラグ, 2. 環境変数, 3. 設定ファイル
		if apiURL == "" {
			apiURL = os.Getenv("REDMINE_API_URL")
		}
		if apiKey == "" {
			apiKey = os.Getenv("REDMINE_API_KEY")
		}

		// 設定ファイルから読み込み（フラグと環境変数が未設定の場合）
		if apiURL == "" || apiKey == "" {
			cfg, err := cliconfig.Load()
			if err == nil {
				if apiURL == "" {
					apiURL = cfg.APIURL
				}
				if apiKey == "" {
					apiKey = cfg.APIKey
				}
			}
		}

		if apiURL == "" {
			return errors.New("REDMINE_API_URL が設定されていません。以下のいずれかの方法で設定してください:\n  1. 'redmine-cli config init' で設定ファイルを作成\n  2. --url フラグを指定\n  3. REDMINE_API_URL 環境変数を設定")
		}
		if apiKey == "" {
			return errors.New("REDMINE_API_KEY が設定されていません。以下のいずれかの方法で設定してください:\n  1. 'redmine-cli config init' で設定ファイルを作成\n  2. --key フラグを指定\n  3. REDMINE_API_KEY 環境変数を設定")
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
	rootCmd.PersistentFlags().StringVar(&apiURL, "url", "", "Redmine API URL (優先順位: フラグ > 環境変数 > 設定ファイル)")
	rootCmd.PersistentFlags().StringVar(&apiKey, "key", "", "Redmine API Key (優先順位: フラグ > 環境変数 > 設定ファイル)")
	rootCmd.PersistentFlags().StringVarP(&outputFormat, "format", "f", "table", "出力フォーマット (json, table, text)")
}

// GetOutputFormat returns the current output format.
func GetOutputFormat() string {
	return outputFormat
}
