package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var groupCmd = &cobra.Command{
	Use:   "group",
	Short: "Manage Redmine user groups",
	Long:  `ユーザーグループの作成、取得、更新、削除などの操作を行います。`,
}

var groupListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all groups",
	Long:  `すべてのユーザーグループをリスト表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		include, _ := cmd.Flags().GetString("include")
		format, _ := cmd.Flags().GetString("format")

		opts := &redmine.ListGroupsOptions{
			Include: include,
		}

		result, err := client.ListGroups(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("グループの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatGroupsTable(result.Groups)
		case formatText:
			return formatGroupsText(result.Groups)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

var groupGetCmd = &cobra.Command{
	Use:   "get [group_id]",
	Short: "Get a group by ID",
	Long:  `指定したIDのグループを取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なgroup_id: %w", err)
		}

		include, _ := cmd.Flags().GetString("include")
		format, _ := cmd.Flags().GetString("format")

		opts := &redmine.ShowGroupOptions{
			Include: include,
		}

		result, err := client.ShowGroup(context.Background(), id, opts)
		if err != nil {
			return fmt.Errorf("グループの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatGroupDetail(&result.Group)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
	},
}

var groupCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new group",
	Long:  `新しいグループを作成します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		userIDsStr, _ := cmd.Flags().GetString("user-ids")

		if name == "" {
			return errors.New("--name フラグは必須です")
		}

		group := redmine.Group{
			Name: name,
		}

		if userIDsStr != "" {
			// Parse user IDs from comma-separated string
			userIDsStrs := strings.Split(userIDsStr, ",")
			userIDs := make([]int, 0, len(userIDsStrs))
			for _, idStr := range userIDsStrs {
				id, err := strconv.Atoi(strings.TrimSpace(idStr))
				if err != nil {
					return fmt.Errorf("無効なuser_id: %s", idStr)
				}
				userIDs = append(userIDs, id)
			}
			group.UserIDs = userIDs
		}

		result, err := client.CreateGroup(context.Background(), group)
		if err != nil {
			return fmt.Errorf("グループの作成に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var groupUpdateCmd = &cobra.Command{
	Use:   "update [group_id]",
	Short: "Update an existing group",
	Long:  `既存のグループを更新します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なgroup_id: %w", err)
		}

		name, _ := cmd.Flags().GetString("name")
		userIDsStr, _ := cmd.Flags().GetString("user-ids")

		group := redmine.Group{
			Name: name,
		}

		if userIDsStr != "" {
			// Parse user IDs from comma-separated string
			userIDsStrs := strings.Split(userIDsStr, ",")
			userIDs := make([]int, 0, len(userIDsStrs))
			for _, idStr := range userIDsStrs {
				userID, err := strconv.Atoi(strings.TrimSpace(idStr))
				if err != nil {
					return fmt.Errorf("無効なuser_id: %s", idStr)
				}
				userIDs = append(userIDs, userID)
			}
			group.UserIDs = userIDs
		}

		err = client.UpdateGroup(context.Background(), id, group)
		if err != nil {
			return fmt.Errorf("グループの更新に失敗しました: %w", err)
		}

		fmt.Println("グループを更新しました")
		return nil
	},
}

var groupDeleteCmd = &cobra.Command{
	Use:   "delete [group_id]",
	Short: "Delete a group",
	Long:  `グループを削除します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なgroup_id: %w", err)
		}

		err = client.DeleteGroup(context.Background(), id)
		if err != nil {
			return fmt.Errorf("グループの削除に失敗しました: %w", err)
		}

		fmt.Println("グループを削除しました")
		return nil
	},
}

var groupAddUserCmd = &cobra.Command{
	Use:   "add-user [group_id] [user_id]",
	Short: "Add a user to a group",
	Long:  `グループにユーザーを追加します。`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なgroup_id: %w", err)
		}

		userID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("無効なuser_id: %w", err)
		}

		err = client.AddUserToGroup(context.Background(), groupID, userID)
		if err != nil {
			return fmt.Errorf("グループへのユーザー追加に失敗しました: %w", err)
		}

		fmt.Println("グループにユーザーを追加しました")
		return nil
	},
}

var groupRemoveUserCmd = &cobra.Command{
	Use:   "remove-user [group_id] [user_id]",
	Short: "Remove a user from a group",
	Long:  `グループからユーザーを削除します。`,
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		groupID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なgroup_id: %w", err)
		}

		userID, err := strconv.Atoi(args[1])
		if err != nil {
			return fmt.Errorf("無効なuser_id: %w", err)
		}

		err = client.RemoveUserFromGroup(context.Background(), groupID, userID)
		if err != nil {
			return fmt.Errorf("グループからのユーザー削除に失敗しました: %w", err)
		}

		fmt.Println("グループからユーザーを削除しました")
		return nil
	},
}

// formatGroupDetail formats a single group in detailed text format.
func formatGroupDetail(g *redmine.Group) error {
	// Title
	fmt.Println(formatter.FormatTitle("Group: " + g.Name))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(g.ID)))
	fmt.Println(formatter.FormatKeyValue("Name", g.Name))

	return nil
}

// formatGroupsTable formats groups in table format.
func formatGroupsTable(groups []redmine.Group) error {
	if len(groups) == 0 {
		fmt.Println("グループが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Name"}
	rows := make([][]string, 0, len(groups))

	for _, g := range groups {
		rows = append(rows, []string{
			strconv.Itoa(g.ID),
			formatter.TruncateString(g.Name, 60),
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatGroupsText formats groups in simple text format.
func formatGroupsText(groups []redmine.Group) error {
	if len(groups) == 0 {
		fmt.Println("グループが見つかりませんでした。")
		return nil
	}

	for _, g := range groups {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(g.ID)))
		fmt.Println(formatter.FormatKeyValue("Name", g.Name))
		fmt.Println()
	}

	return nil
}

// includeOptionsForGroup returns valid include options for group commands
func includeOptionsForGroup() []string {
	return []string{"users", "memberships"}
}

func init() {
	rootCmd.AddCommand(groupCmd)

	// Subcommands
	groupCmd.AddCommand(groupListCmd)
	groupCmd.AddCommand(groupGetCmd)
	groupCmd.AddCommand(groupCreateCmd)
	groupCmd.AddCommand(groupUpdateCmd)
	groupCmd.AddCommand(groupDeleteCmd)
	groupCmd.AddCommand(groupAddUserCmd)
	groupCmd.AddCommand(groupRemoveUserCmd)

	// Flags for list command
	groupListCmd.Flags().String("include", "", "追加で取得する情報 (users, memberships)")
	groupListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")

	// Register flag completion for list command
	_ = groupListCmd.RegisterFlagCompletionFunc("include", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return includeOptionsForGroup(), cobra.ShellCompDirectiveNoFileComp
	})

	// Flags for get command
	groupGetCmd.Flags().String("include", "", "追加で取得する情報 (users, memberships)")
	groupGetCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")

	// Register flag completion for get command
	_ = groupGetCmd.RegisterFlagCompletionFunc("include", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return includeOptionsForGroup(), cobra.ShellCompDirectiveNoFileComp
	})

	// Flags for create command
	groupCreateCmd.Flags().String("name", "", "グループ名 (必須)")
	groupCreateCmd.Flags().String("user-ids", "", "ユーザーIDのカンマ区切りリスト (例: 1,2,3)")

	// Flags for update command
	groupUpdateCmd.Flags().String("name", "", "グループ名")
	groupUpdateCmd.Flags().String("user-ids", "", "ユーザーIDのカンマ区切りリスト (例: 1,2,3)")
}
