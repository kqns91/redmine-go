package cmd

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Redmine CLI configuration",
	Long:  `Redmine CLIの設定を管理します。`,
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration interactively",
	Long:  `対話的に設定ファイルを初期化します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Check if config already exists
		exists, err := config.Exists()
		if err != nil {
			return fmt.Errorf("設定ファイルの確認に失敗しました: %w", err)
		}

		if exists {
			fmt.Println("警告: 設定ファイルが既に存在します。")
			fmt.Print("上書きしますか？ (y/N): ")
			var confirm string
			_, _ = fmt.Scanln(&confirm)
			if confirm != "y" && confirm != "yes" {
				fmt.Println("初期化をキャンセルしました。")
				return nil
			}
		}

		// Interactive initialization
		cfg, err := config.InitInteractive()
		if err != nil {
			return err
		}

		// Save configuration
		if err := config.Save(cfg); err != nil {
			return err
		}

		configPath, _ := config.GetConfigPath()
		fmt.Println()
		fmt.Printf("設定を保存しました: %s\n", configPath)
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `現在の設定内容を表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		configPath, _ := config.GetConfigPath()
		fmt.Printf("設定ファイル: %s\n", configPath)
		fmt.Println()

		output, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			return fmt.Errorf("設定のシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Long:  `設定値を更新します。使用可能なキー: api_url, api_key`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		// Load existing config
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		// Update the specified key
		switch key {
		case "api_url":
			cfg.APIURL = value
			fmt.Printf("API URLを更新しました: %s\n", value)
		case "api_key":
			cfg.APIKey = value
			fmt.Println("API Keyを更新しました")
		default:
			return errors.New("無効なキーです。使用可能なキー: api_url, api_key")
		}

		// Save configuration
		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Println("設定を保存しました。")
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Long:  `設定ファイルのパスを表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		configPath, err := config.GetConfigPath()
		if err != nil {
			return err
		}

		fmt.Println(configPath)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Subcommands
	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configPathCmd)
}
