package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/pkg/redmine"
)

var issueCategoryCmd = &cobra.Command{
	Use:   "issue-category",
	Short: "Manage Redmine issue categories",
	Long:  `チケットカテゴリの作成、取得、更新、削除などの操作を行います。`,
}

var issueCategoryListCmd = &cobra.Command{
	Use:   "list [project_id_or_identifier]",
	Short: "List issue categories for a project",
	Long:  `指定したプロジェクトのチケットカテゴリをリスト表示します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		result, err := client.ListIssueCategories(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("チケットカテゴリの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var issueCategoryGetCmd = &cobra.Command{
	Use:   "get [category_id]",
	Short: "Get an issue category by ID",
	Long:  `指定したIDのチケットカテゴリを取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なcategory_id: %w", err)
		}

		result, err := client.ShowIssueCategory(context.Background(), id)
		if err != nil {
			return fmt.Errorf("チケットカテゴリの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var issueCategoryCreateCmd = &cobra.Command{
	Use:   "create [project_id_or_identifier]",
	Short: "Create a new issue category",
	Long:  `指定したプロジェクトに新しいチケットカテゴリを作成します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		assignedToID, _ := cmd.Flags().GetInt("assigned-to-id")

		if name == "" {
			return errors.New("--name フラグは必須です")
		}

		req := redmine.IssueCategoryCreateRequest{
			Name:         name,
			AssignedToID: assignedToID,
		}

		result, err := client.CreateIssueCategory(context.Background(), args[0], req)
		if err != nil {
			return fmt.Errorf("チケットカテゴリの作成に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var issueCategoryUpdateCmd = &cobra.Command{
	Use:   "update [category_id]",
	Short: "Update an existing issue category",
	Long:  `既存のチケットカテゴリを更新します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なcategory_id: %w", err)
		}

		name, _ := cmd.Flags().GetString("name")
		assignedToID, _ := cmd.Flags().GetInt("assigned-to-id")

		req := redmine.IssueCategoryUpdateRequest{
			Name:         name,
			AssignedToID: assignedToID,
		}

		err = client.UpdateIssueCategory(context.Background(), id, req)
		if err != nil {
			return fmt.Errorf("チケットカテゴリの更新に失敗しました: %w", err)
		}

		fmt.Println("チケットカテゴリを更新しました")
		return nil
	},
}

var issueCategoryDeleteCmd = &cobra.Command{
	Use:   "delete [category_id]",
	Short: "Delete an issue category",
	Long:  `チケットカテゴリを削除します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なcategory_id: %w", err)
		}

		reassignToID, _ := cmd.Flags().GetInt("reassign-to-id")

		var opts *redmine.DeleteIssueCategoryOptions
		if reassignToID > 0 {
			opts = &redmine.DeleteIssueCategoryOptions{
				ReassignToID: reassignToID,
			}
		}

		err = client.DeleteIssueCategory(context.Background(), id, opts)
		if err != nil {
			return fmt.Errorf("チケットカテゴリの削除に失敗しました: %w", err)
		}

		fmt.Println("チケットカテゴリを削除しました")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(issueCategoryCmd)

	// Subcommands
	issueCategoryCmd.AddCommand(issueCategoryListCmd)
	issueCategoryCmd.AddCommand(issueCategoryGetCmd)
	issueCategoryCmd.AddCommand(issueCategoryCreateCmd)
	issueCategoryCmd.AddCommand(issueCategoryUpdateCmd)
	issueCategoryCmd.AddCommand(issueCategoryDeleteCmd)

	// Flags for create command
	issueCategoryCreateCmd.Flags().String("name", "", "カテゴリ名 (必須)")
	issueCategoryCreateCmd.Flags().Int("assigned-to-id", 0, "デフォルト担当者ID")

	// Flags for update command
	issueCategoryUpdateCmd.Flags().String("name", "", "カテゴリ名")
	issueCategoryUpdateCmd.Flags().Int("assigned-to-id", 0, "デフォルト担当者ID")

	// Flags for delete command
	issueCategoryDeleteCmd.Flags().Int("reassign-to-id", 0, "削除前にチケットを別カテゴリに再割り当てするカテゴリID")
}
