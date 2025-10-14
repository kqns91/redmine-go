package cmd

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/cli/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
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
		format, _ := cmd.Flags().GetString("format")

		result, err := client.ListRoles(context.Background())
		if err != nil {
			return fmt.Errorf("ロールの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatRolesTable(result.Roles)
		case formatText:
			return formatRolesText(result.Roles)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
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

		format, _ := cmd.Flags().GetString("format")

		result, err := client.ShowRole(context.Background(), id)
		if err != nil {
			return fmt.Errorf("ロールの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatRoleDetail(&result.Role)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
	},
}

// formatRoleDetail formats a single role in detailed text format.
func formatRoleDetail(r *redmine.Role) error {
	// Title
	fmt.Println(formatter.FormatTitle("Role: " + r.Name))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(r.ID)))
	fmt.Println(formatter.FormatKeyValue("Name", r.Name))
	fmt.Println(formatter.FormatKeyValue("Assignable", strconv.FormatBool(r.Assignable)))

	// Permissions
	if len(r.Permissions) > 0 {
		fmt.Println()
		fmt.Println(formatter.FormatSection("権限"))
		for _, perm := range r.Permissions {
			fmt.Println("  - " + perm)
		}
	}

	return nil
}

// formatRolesTable formats roles in table format.
func formatRolesTable(roles []redmine.Role) error {
	if len(roles) == 0 {
		fmt.Println("ロールが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Name", "Assignable", "Permissions"}
	rows := make([][]string, 0, len(roles))

	for _, r := range roles {
		permCount := strconv.Itoa(len(r.Permissions))
		if len(r.Permissions) == 0 {
			permCount = "-"
		}

		rows = append(rows, []string{
			strconv.Itoa(r.ID),
			formatter.TruncateString(r.Name, 30),
			strconv.FormatBool(r.Assignable),
			permCount,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatRolesText formats roles in simple text format.
func formatRolesText(roles []redmine.Role) error {
	if len(roles) == 0 {
		fmt.Println("ロールが見つかりませんでした。")
		return nil
	}

	for _, r := range roles {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(r.ID)))
		fmt.Println(formatter.FormatKeyValue("Name", r.Name))
		fmt.Println(formatter.FormatKeyValue("Assignable", strconv.FormatBool(r.Assignable)))
		if len(r.Permissions) > 0 {
			perms := strings.Join(r.Permissions, ", ")
			fmt.Println(formatter.FormatKeyValue("Permissions", formatter.TruncateString(perms, 80)))
		}
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(roleCmd)

	// Subcommands
	roleCmd.AddCommand(roleListCmd)
	roleCmd.AddCommand(roleGetCmd)

	// Flags for list command
	roleListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")

	// Flags for get command
	roleGetCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")
}
