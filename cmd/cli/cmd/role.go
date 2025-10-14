package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"
)

var roleCmd = &cobra.Command{
	Use:   "role",
	Short: "Manage Redmine roles",
	Long:  `ロールの取得などの操作を行います。`,
}

var roleListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all roles",
	Long:  `すべてのロールをリスト表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.ListRoles(context.Background())
		if err != nil {
			return fmt.Errorf("ロールの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var roleGetCmd = &cobra.Command{
	Use:   "get [role_id]",
	Short: "Get a role by ID",
	Long:  `指定したIDのロール（権限情報を含む）を取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なrole_id: %w", err)
		}

		result, err := client.ShowRole(context.Background(), id)
		if err != nil {
			return fmt.Errorf("ロールの取得に失敗しました: %w", err)
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
	rootCmd.AddCommand(roleCmd)

	// Subcommands
	roleCmd.AddCommand(roleListCmd)
	roleCmd.AddCommand(roleGetCmd)
}
