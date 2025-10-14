package config

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	configDirName  = ".config/redmine"
	configFileName = "config"
)

// Config represents the CLI configuration
type Config struct {
	APIURL string `json:"api_url"`
	APIKey string `json:"api_key"`
}

// GetConfigPath returns the full path to the config file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ホームディレクトリの取得に失敗しました: %w", err)
	}

	return filepath.Join(homeDir, configDirName, configFileName), nil
}

// Load loads the configuration from the config file
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	//nolint:gosec // Config file path is constructed internally, not from user input
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("設定ファイルが見つかりません。'redmine-cli config init' を実行して初期設定を行ってください")
		}
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("設定ファイルのパースに失敗しました: %w", err)
	}

	return &cfg, nil
}

// Save saves the configuration to the config file
func Save(cfg *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	//nolint:gosec // 0755 is appropriate for config directory
	if err := os.MkdirAll(configDir, 0o755); err != nil {
		return fmt.Errorf("設定ディレクトリの作成に失敗しました: %w", err)
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("設定のシリアライズに失敗しました: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0o600); err != nil {
		return fmt.Errorf("設定ファイルの書き込みに失敗しました: %w", err)
	}

	return nil
}

// Exists checks if the config file exists
func Exists() (bool, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return false, err
	}

	_, err = os.Stat(configPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// InitInteractive interactively initializes the configuration
func InitInteractive() (*Config, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Redmine CLI 初期設定")
	fmt.Println("==================")
	fmt.Println()

	// Get API URL
	fmt.Print("Redmine API URL (例: https://redmine.example.com): ")
	apiURL, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("入力の読み込みに失敗しました: %w", err)
	}
	apiURL = strings.TrimSpace(apiURL)
	if apiURL == "" {
		return nil, errors.New("API URLは必須です")
	}

	// Get API Key
	fmt.Print("Redmine API Key: ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("入力の読み込みに失敗しました: %w", err)
	}
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return nil, errors.New("API Keyは必須です")
	}

	cfg := &Config{
		APIURL: apiURL,
		APIKey: apiKey,
	}

	// Confirm
	fmt.Println()
	fmt.Println("設定内容:")
	fmt.Printf("  API URL: %s\n", cfg.APIURL)
	fmt.Printf("  API Key: %s\n", maskAPIKey(cfg.APIKey))
	fmt.Println()
	fmt.Print("この設定で保存しますか？ (y/N): ")

	confirm, err := reader.ReadString('\n')
	if err != nil {
		return nil, fmt.Errorf("入力の読み込みに失敗しました: %w", err)
	}
	confirm = strings.ToLower(strings.TrimSpace(confirm))

	if confirm != "y" && confirm != "yes" {
		return nil, errors.New("設定の保存がキャンセルされました")
	}

	return cfg, nil
}

// maskAPIKey masks the API key for display
func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "********"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
