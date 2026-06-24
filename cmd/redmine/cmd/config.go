package cmd

import (
	"errors"
	"fmt"
	"sort"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/internal/config"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Redmine CLI configuration",
	Long:  `Redmine CLIの設定を管理します。`,
}

// loadOrEmpty loads the config, returning an empty config if the file does not
// exist yet (so profiles can be added before the first save).
func loadOrEmpty() (*config.Config, error) {
	exists, err := config.Exists()
	if err != nil {
		return nil, fmt.Errorf("設定ファイルの確認に失敗しました: %w", err)
	}
	if !exists {
		return &config.Config{Profiles: map[string]config.Profile{}}, nil
	}
	return config.Load()
}

var configInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize configuration interactively",
	Long:  `対話的に設定ファイルを初期化します。`,
	RunE: func(_ *cobra.Command, _ []string) error {
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

		profile, err := config.PromptProfile()
		if err != nil {
			return err
		}

		cfg := &config.Config{
			CurrentProfile: config.DefaultProfileName,
			Profiles:       map[string]config.Profile{config.DefaultProfileName: profile},
		}
		if err := config.Save(cfg); err != nil {
			return err
		}

		configPath, _ := config.GetConfigPath()
		fmt.Println()
		fmt.Printf("プロファイル %q を保存しました: %s\n", config.DefaultProfileName, configPath)
		return nil
	},
}

var configAddCmd = &cobra.Command{
	Use:   "add [name]",
	Short: "Add a new profile",
	Long:  `新しいプロファイルを対話的に追加します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := loadOrEmpty()
		if err != nil {
			return err
		}

		if _, ok := cfg.Profiles[name]; ok {
			return fmt.Errorf("プロファイル %q は既に存在します", name)
		}

		profile, err := config.PromptProfile()
		if err != nil {
			return err
		}

		cfg.Profiles[name] = profile
		// 最初のプロファイルなら current に設定する
		if cfg.CurrentProfile == "" {
			cfg.CurrentProfile = name
		}

		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("プロファイル %q を追加しました。\n", name)
		if cfg.CurrentProfile == name {
			fmt.Printf("現在のプロファイルは %q です。\n", name)
		}
		return nil
	},
}

var configUseCmd = &cobra.Command{
	Use:   "use [name]",
	Short: "Switch the current profile",
	Long:  `現在のプロファイルを切り替えます。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if _, ok := cfg.Profiles[name]; !ok {
			return fmt.Errorf("プロファイル %q が見つかりません", name)
		}

		cfg.CurrentProfile = name
		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("現在のプロファイルを %q に切り替えました。\n", name)
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List profiles",
	Long:  `登録済みのプロファイル一覧を表示します。`,
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(cfg.Profiles) == 0 {
			fmt.Println("プロファイルが登録されていません。'redmine config add <name>' で追加してください。")
			return nil
		}

		names := make([]string, 0, len(cfg.Profiles))
		for name := range cfg.Profiles {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			profile := cfg.Profiles[name]
			marker := " "
			if name == cfg.CurrentProfile {
				marker = "*"
			}
			fmt.Printf("%s %s\t%s\t%s\n", marker, name, profile.APIURL, config.MaskAPIKey(profile.APIKey))
		}
		return nil
	},
}

var configRemoveCmd = &cobra.Command{
	Use:   "remove [name]",
	Short: "Remove a profile",
	Long:  `プロファイルを削除します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(_ *cobra.Command, args []string) error {
		name := args[0]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if _, ok := cfg.Profiles[name]; !ok {
			return fmt.Errorf("プロファイル %q が見つかりません", name)
		}

		delete(cfg.Profiles, name)
		// current を削除した場合は選択を解除する
		if cfg.CurrentProfile == name {
			cfg.CurrentProfile = ""
		}

		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("プロファイル %q を削除しました。\n", name)
		if cfg.CurrentProfile == "" && len(cfg.Profiles) > 0 {
			fmt.Println("現在のプロファイルが未選択です。'redmine config use <name>' で選択してください。")
		}
		return nil
	},
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current profile",
	Long:  `現在のプロファイルの設定内容を表示します（API キーはマスクされます）。`,
	RunE: func(_ *cobra.Command, _ []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		profile, err := cfg.ResolveProfile("")
		if err != nil {
			return err
		}

		fmt.Printf("プロファイル: %s\n", cfg.CurrentProfile)
		fmt.Printf("  URL:     %s\n", profile.APIURL)
		fmt.Printf("  API Key: %s\n", config.MaskAPIKey(profile.APIKey))
		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a value on the current profile",
	Long:  `現在のプロファイルの設定値を更新します。使用可能なキー: api_url, api_key`,
	Args:  cobra.ExactArgs(2),
	RunE: func(_ *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]

		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if cfg.CurrentProfile == "" {
			return errors.New("プロファイルが選択されていません。'redmine config use <name>' で選択してください")
		}
		profile := cfg.Profiles[cfg.CurrentProfile]

		switch key {
		case "api_url":
			profile.APIURL = value
			fmt.Printf("URLを更新しました: %s\n", value)
		case "api_key":
			profile.APIKey = value
			fmt.Println("API Keyを更新しました")
		default:
			return errors.New("無効なキーです。使用可能なキー: api_url, api_key")
		}

		cfg.Profiles[cfg.CurrentProfile] = profile
		if err := config.Save(cfg); err != nil {
			return err
		}

		fmt.Printf("プロファイル %q を保存しました。\n", cfg.CurrentProfile)
		return nil
	},
}

var configPathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show configuration file path",
	Long:  `設定ファイルのパスを表示します。`,
	RunE: func(_ *cobra.Command, _ []string) error {
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

	configCmd.AddCommand(configInitCmd)
	configCmd.AddCommand(configAddCmd)
	configCmd.AddCommand(configUseCmd)
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configRemoveCmd)
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configPathCmd)
}
