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

var issueRelationCmd = &cobra.Command{
	Use:   "issue-relation",
	Short: "Manage Redmine issue relations",
	Long:  `チケット間の関連の作成、取得、削除などの操作を行います。`,
}

var issueRelationListCmd = &cobra.Command{
	Use:   "list [issue_id]",
	Short: "List relations for an issue",
	Long:  `指定したチケットの関連をリスト表示します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なissue_id: %w", err)
		}

		result, err := client.ListIssueRelations(context.Background(), issueID)
		if err != nil {
			return fmt.Errorf("チケット関連の取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var issueRelationGetCmd = &cobra.Command{
	Use:   "get [relation_id]",
	Short: "Get an issue relation by ID",
	Long:  `指定したIDのチケット関連を取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		relationID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なrelation_id: %w", err)
		}

		result, err := client.ShowIssueRelation(context.Background(), relationID)
		if err != nil {
			return fmt.Errorf("チケット関連の取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var issueRelationCreateCmd = &cobra.Command{
	Use:   "create [issue_id]",
	Short: "Create a new issue relation",
	Long:  `指定したチケットに新しい関連を作成します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		issueID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なissue_id: %w", err)
		}

		issueToID, _ := cmd.Flags().GetInt("issue-to-id")
		relationType, _ := cmd.Flags().GetString("relation-type")
		delay, _ := cmd.Flags().GetInt("delay")

		if issueToID == 0 {
			return errors.New("--issue-to-id フラグは必須です")
		}
		if relationType == "" {
			return errors.New("--relation-type フラグは必須です")
		}

		relation := redmine.IssueRelation{
			IssueToID:    issueToID,
			RelationType: relationType,
			Delay:        delay,
		}

		result, err := client.CreateIssueRelation(context.Background(), issueID, relation)
		if err != nil {
			return fmt.Errorf("チケット関連の作成に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var issueRelationDeleteCmd = &cobra.Command{
	Use:   "delete [relation_id]",
	Short: "Delete an issue relation",
	Long:  `チケット関連を削除します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		relationID, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なrelation_id: %w", err)
		}

		err = client.DeleteIssueRelation(context.Background(), relationID)
		if err != nil {
			return fmt.Errorf("チケット関連の削除に失敗しました: %w", err)
		}

		fmt.Println("チケット関連を削除しました")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(issueRelationCmd)

	// Subcommands
	issueRelationCmd.AddCommand(issueRelationListCmd)
	issueRelationCmd.AddCommand(issueRelationGetCmd)
	issueRelationCmd.AddCommand(issueRelationCreateCmd)
	issueRelationCmd.AddCommand(issueRelationDeleteCmd)

	// Flags for create command
	issueRelationCreateCmd.Flags().Int("issue-to-id", 0, "関連先チケットID (必須)")
	issueRelationCreateCmd.Flags().String("relation-type", "", "関連タイプ (必須: relates, duplicates, duplicated, blocks, blocked, precedes, follows, copied_to, copied_from)")
	issueRelationCreateCmd.Flags().Int("delay", 0, "遅延日数 (precedesまたはfollowsの場合)")
}
