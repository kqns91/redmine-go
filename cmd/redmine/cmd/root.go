package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cliconfig "github.com/kqns91/redmine-go/internal/config"
	"github.com/kqns91/redmine-go/internal/version"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

const (
	formatJSON  = "json"
	formatTable = "table"
	formatText  = "text"
)

var (
	client      *redmine.Client
	profileFlag string
	configFlag  string
)

// rootCmd はCLIのルートコマンドを表します
var rootCmd = &cobra.Command{
	Use:     "redmine",
	Short:   "Redmine API client CLI",
	Version: version.GetVersion(),
	Long: `redmine は Redmine の REST API を操作するための CLI ツールです。
すべての Redmine API 操作を CLI から実行できます。`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// --config はすべてのコマンド（config サブコマンド含む）に適用する
		cliconfig.SetConfigPath(configFlag)

		// Skip config initialization for config commands
		if cmd.Parent() != nil && cmd.Parent().Name() == "config" {
			return nil
		}
		if cmd.Name() == "config" {
			return nil
		}

		// 優先順位: 1. 環境変数, 2. 設定ファイル
		apiURL := cliconfig.URLFromEnv()
		apiKey := cliconfig.APIKeyFromEnv()

		// 環境変数で揃っていなければ、設定ファイルの選択プロファイルから補完する
		if apiURL == "" || apiKey == "" {
			cfg, err := cliconfig.Load()
			if err == nil {
				// プロファイル選択: --profile フラグ > REDMINE_PROFILE > current_profile
				name := profileFlag
				if name == "" {
					name = cliconfig.ProfileFromEnv()
				}
				profile, err := cfg.ResolveProfile(name)
				if err != nil {
					return err
				}
				if apiURL == "" {
					apiURL = profile.APIURL
				}
				if apiKey == "" {
					apiKey = profile.APIKey
				}
			}
		}

		if apiURL == "" {
			return errors.New("REDMINE_URL が設定されていません。以下のいずれかの方法で設定してください:\n  1. 'redmine config init' で設定ファイルを作成\n  2. REDMINE_URL 環境変数を設定")
		}
		if apiKey == "" {
			return errors.New("REDMINE_API_KEY が設定されていません。以下のいずれかの方法で設定してください:\n  1. 'redmine config init' で設定ファイルを作成\n  2. REDMINE_API_KEY 環境変数を設定")
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
	rootCmd.SetVersionTemplate(
		fmt.Sprintf("redmine version %s (commit: %s, built: %s)\n",
			version.GetVersion(), version.GetCommit(), version.GetDate()))

	rootCmd.PersistentFlags().StringVar(&profileFlag, "profile", "",
		"使用するプロファイル名 (デフォルト: 設定ファイルの current_profile)")
	rootCmd.PersistentFlags().StringVar(&configFlag, "config", "",
		"設定ファイルのパス (デフォルト: ~/.config/redmine/config, REDMINE_CONFIG でも指定可)")
}
