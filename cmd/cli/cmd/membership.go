package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/cli/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var membershipCmd = &cobra.Command{
	Use:   "membership",
	Short: "Manage Redmine project memberships",
	Long:  `プロジェクトメンバーシップの作成、取得、更新、削除などの操作を行います。`,
}

var membershipListCmd = &cobra.Command{
	Use:   "list [project_id_or_identifier]",
	Short: "List memberships for a project",
	Long:  `指定したプロジェクトのメンバーシップをリスト表示します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")

		result, err := client.ListMemberships(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("メンバーシップの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatMembershipsTable(result.Memberships)
		case formatText:
			return formatMembershipsText(result.Memberships)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

var membershipGetCmd = &cobra.Command{
	Use:   "get [membership_id]",
	Short: "Get a membership by ID",
	Long:  `指定したIDのメンバーシップを取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なmembership_id: %w", err)
		}

		result, err := client.ShowMembership(context.Background(), id)
		if err != nil {
			return fmt.Errorf("メンバーシップの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var membershipCreateCmd = &cobra.Command{
	Use:   "create [project_id_or_identifier]",
	Short: "Create a new membership",
	Long:  `指定したプロジェクトに新しいメンバーシップを作成します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		userID, _ := cmd.Flags().GetInt("user-id")
		roleIDsStr, _ := cmd.Flags().GetString("role-ids")

		if userID == 0 {
			return errors.New("--user-id フラグは必須です")
		}
		if roleIDsStr == "" {
			return errors.New("--role-ids フラグは必須です")
		}

		// Parse role IDs from comma-separated string
		roleIDsStrs := strings.Split(roleIDsStr, ",")
		roleIDs := make([]int, 0, len(roleIDsStrs))
		for _, idStr := range roleIDsStrs {
			id, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				return fmt.Errorf("無効なrole_id: %s", idStr)
			}
			roleIDs = append(roleIDs, id)
		}

		membership := redmine.MembershipCreateUpdate{
			UserID:  userID,
			RoleIDs: roleIDs,
		}

		result, err := client.CreateMembership(context.Background(), args[0], membership)
		if err != nil {
			return fmt.Errorf("メンバーシップの作成に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var membershipUpdateCmd = &cobra.Command{
	Use:   "update [membership_id]",
	Short: "Update an existing membership",
	Long:  `既存のメンバーシップのロールを更新します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なmembership_id: %w", err)
		}

		roleIDsStr, _ := cmd.Flags().GetString("role-ids")
		if roleIDsStr == "" {
			return errors.New("--role-ids フラグは必須です")
		}

		// Parse role IDs from comma-separated string
		roleIDsStrs := strings.Split(roleIDsStr, ",")
		roleIDs := make([]int, 0, len(roleIDsStrs))
		for _, idStr := range roleIDsStrs {
			roleID, err := strconv.Atoi(strings.TrimSpace(idStr))
			if err != nil {
				return fmt.Errorf("無効なrole_id: %s", idStr)
			}
			roleIDs = append(roleIDs, roleID)
		}

		err = client.UpdateMembership(context.Background(), id, roleIDs)
		if err != nil {
			return fmt.Errorf("メンバーシップの更新に失敗しました: %w", err)
		}

		fmt.Println("メンバーシップを更新しました")
		return nil
	},
}

var membershipDeleteCmd = &cobra.Command{
	Use:   "delete [membership_id]",
	Short: "Delete a membership",
	Long:  `メンバーシップを削除します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なmembership_id: %w", err)
		}

		err = client.DeleteMembership(context.Background(), id)
		if err != nil {
			return fmt.Errorf("メンバーシップの削除に失敗しました: %w", err)
		}

		fmt.Println("メンバーシップを削除しました")
		return nil
	},
}

// formatMembershipsTable formats memberships in table format.
func formatMembershipsTable(memberships []redmine.Membership) error {
	if len(memberships) == 0 {
		fmt.Println("メンバーシップが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "User/Group", "Project", "Roles"}
	rows := make([][]string, 0, len(memberships))

	for _, m := range memberships {
		userOrGroup := "-"
		if m.User.Name != "" {
			userOrGroup = m.User.Name
		} else if m.Group.Name != "" {
			userOrGroup = m.Group.Name + " (G)"
		}

		roleNames := make([]string, 0, len(m.Roles))
		for _, r := range m.Roles {
			roleNames = append(roleNames, r.Name)
		}
		roles := strings.Join(roleNames, ", ")

		rows = append(rows, []string{
			strconv.Itoa(m.ID),
			formatter.TruncateString(userOrGroup, 25),
			formatter.TruncateString(m.Project.Name, 25),
			formatter.TruncateString(roles, 30),
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatMembershipsText formats memberships in simple text format.
func formatMembershipsText(memberships []redmine.Membership) error {
	if len(memberships) == 0 {
		fmt.Println("メンバーシップが見つかりませんでした。")
		return nil
	}

	for _, m := range memberships {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(m.ID)))
		if m.User.Name != "" {
			fmt.Println(formatter.FormatKeyValue("User", m.User.Name))
		} else if m.Group.Name != "" {
			fmt.Println(formatter.FormatKeyValue("Group", m.Group.Name))
		}
		fmt.Println(formatter.FormatKeyValue("Project", m.Project.Name))

		roleNames := make([]string, 0, len(m.Roles))
		for _, r := range m.Roles {
			roleNames = append(roleNames, r.Name)
		}
		fmt.Println(formatter.FormatKeyValue("Roles", strings.Join(roleNames, ", ")))
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(membershipCmd)

	// Subcommands
	membershipCmd.AddCommand(membershipListCmd)
	membershipCmd.AddCommand(membershipGetCmd)
	membershipCmd.AddCommand(membershipCreateCmd)
	membershipCmd.AddCommand(membershipUpdateCmd)
	membershipCmd.AddCommand(membershipDeleteCmd)

	// Flags for list command
	membershipListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")

	// Flags for create command
	membershipCreateCmd.Flags().Int("user-id", 0, "ユーザーID (必須)")
	membershipCreateCmd.Flags().String("role-ids", "", "ロールIDのカンマ区切りリスト (必須、例: 1,2,3)")

	// Flags for update command
	membershipUpdateCmd.Flags().String("role-ids", "", "ロールIDのカンマ区切りリスト (必須、例: 1,2,3)")
}
