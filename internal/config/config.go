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

	// DefaultProfileName is the name given to the profile created when migrating
	// a legacy single-profile config, and the default for the first profile.
	DefaultProfileName = "default"
)

// Profile holds the connection settings for a single Redmine instance.
type Profile struct {
	APIURL string `json:"api_url"`
	APIKey string `json:"api_key"`
}

// Config represents the CLI configuration. It holds one or more named profiles
// and a pointer to the currently selected one.
type Config struct {
	CurrentProfile string             `json:"current_profile"`
	Profiles       map[string]Profile `json:"profiles"`
}

// GetConfigPath returns the full path to the config file.
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("ホームディレクトリの取得に失敗しました: %w", err)
	}

	return filepath.Join(homeDir, configDirName, configFileName), nil
}

// Load loads the configuration from the config file.
//
// A legacy single-profile config ({"api_url": ..., "api_key": ...}) is migrated
// in memory to a profile named DefaultProfileName. The migration is
// non-destructive: the file on disk is only rewritten when a command saves it.
func Load() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	//nolint:gosec // Config file path is constructed internally, not from user input
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("設定ファイルが見つかりません。'redmine config init' を実行して初期設定を行ってください")
		}
		return nil, fmt.Errorf("設定ファイルの読み込みに失敗しました: %w", err)
	}

	// raw accepts both the current profile-based format and the legacy
	// single-profile format so older config files keep working.
	var raw struct {
		CurrentProfile string             `json:"current_profile"`
		Profiles       map[string]Profile `json:"profiles"`
		// Legacy single-profile fields.
		APIURL string `json:"api_url"`
		APIKey string `json:"api_key"`
	}
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("設定ファイルのパースに失敗しました: %w", err)
	}

	cfg := &Config{
		CurrentProfile: raw.CurrentProfile,
		Profiles:       raw.Profiles,
	}

	// Migrate a legacy single-profile config into a named "default" profile.
	if len(cfg.Profiles) == 0 && (raw.APIURL != "" || raw.APIKey != "") {
		cfg.Profiles = map[string]Profile{
			DefaultProfileName: {APIURL: raw.APIURL, APIKey: raw.APIKey},
		}
		cfg.CurrentProfile = DefaultProfileName
	}

	if cfg.Profiles == nil {
		cfg.Profiles = map[string]Profile{}
	}

	return cfg, nil
}

// Save saves the configuration to the config file.
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

// Exists checks if the config file exists.
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

// ResolveProfile returns the profile selected by the given name. If name is
// empty, the current profile is used. It returns an error if no profile is
// selected or the selected profile does not exist.
func (c *Config) ResolveProfile(name string) (Profile, error) {
	if name == "" {
		name = c.CurrentProfile
	}
	if name == "" {
		return Profile{}, errors.New("プロファイルが選択されていません。'redmine config use <name>' で選択してください")
	}

	profile, ok := c.Profiles[name]
	if !ok {
		return Profile{}, fmt.Errorf("プロファイル %q が見つかりません", name)
	}

	return profile, nil
}

// PromptProfile interactively reads a Redmine URL and API key from stdin,
// shows a masked summary, and asks for confirmation. It returns an error if
// the input is invalid or the user cancels.
func PromptProfile() (Profile, error) {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Redmine URL (例: https://redmine.example.com): ")
	apiURL, err := reader.ReadString('\n')
	if err != nil {
		return Profile{}, fmt.Errorf("入力の読み込みに失敗しました: %w", err)
	}
	apiURL = strings.TrimSpace(apiURL)
	if apiURL == "" {
		return Profile{}, errors.New("URLは必須です")
	}

	fmt.Print("Redmine API Key: ")
	apiKey, err := reader.ReadString('\n')
	if err != nil {
		return Profile{}, fmt.Errorf("入力の読み込みに失敗しました: %w", err)
	}
	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return Profile{}, errors.New("API Keyは必須です")
	}

	profile := Profile{APIURL: apiURL, APIKey: apiKey}

	fmt.Println()
	fmt.Println("設定内容:")
	fmt.Printf("  URL:     %s\n", profile.APIURL)
	fmt.Printf("  API Key: %s\n", MaskAPIKey(profile.APIKey))
	fmt.Println()
	fmt.Print("この設定で保存しますか？ (y/N): ")

	confirm, err := reader.ReadString('\n')
	if err != nil {
		return Profile{}, fmt.Errorf("入力の読み込みに失敗しました: %w", err)
	}
	confirm = strings.ToLower(strings.TrimSpace(confirm))
	if confirm != "y" && confirm != "yes" {
		return Profile{}, errors.New("設定の保存がキャンセルされました")
	}

	return profile, nil
}

// MaskAPIKey masks the API key for display.
func MaskAPIKey(key string) string {
	if len(key) <= 8 {
		return "********"
	}
	return key[:4] + "..." + key[len(key)-4:]
}
