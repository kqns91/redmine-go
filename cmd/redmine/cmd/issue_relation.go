package cmd

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
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

		format, _ := cmd.Flags().GetString("format")

		result, err := client.ListIssueRelations(context.Background(), issueID)
		if err != nil {
			return fmt.Errorf("チケット関連の取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatIssueRelationsTable(result.Relations)
		case formatText:
			return formatIssueRelationsText(result.Relations)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
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

		format, _ := cmd.Flags().GetString("format")

		result, err := client.ShowIssueRelation(context.Background(), relationID)
		if err != nil {
			return fmt.Errorf("チケット関連の取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatIssueRelationDetail(&result.Relation)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
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

		return formatter.OutputJSON(result)
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

// formatIssueRelationDetail formats a single issue relation in detailed text format.
func formatIssueRelationDetail(r *redmine.IssueRelation) error {
	// Title
	fmt.Println(formatter.FormatTitle("Issue Relation #" + strconv.Itoa(r.ID)))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(r.ID)))
	fmt.Println(formatter.FormatKeyValue("Issue ID", strconv.Itoa(r.IssueID)))
	fmt.Println(formatter.FormatKeyValue("Issue To ID", strconv.Itoa(r.IssueToID)))
	fmt.Println(formatter.FormatKeyValue("Relation Type", r.RelationType))
	if r.Delay != 0 {
		fmt.Println(formatter.FormatKeyValue("Delay", strconv.Itoa(r.Delay)))
	}

	return nil
}

// formatIssueRelationsTable formats issue relations in table format.
func formatIssueRelationsTable(relations []redmine.IssueRelation) error {
	if len(relations) == 0 {
		fmt.Println("チケット関連が見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Issue ID", "Issue To ID", "Type", "Delay"}
	rows := make([][]string, 0, len(relations))

	for _, r := range relations {
		delayStr := ""
		if r.Delay != 0 {
			delayStr = strconv.Itoa(r.Delay)
		}
		rows = append(rows, []string{
			strconv.Itoa(r.ID),
			strconv.Itoa(r.IssueID),
			strconv.Itoa(r.IssueToID),
			r.RelationType,
			delayStr,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatIssueRelationsText formats issue relations in simple text format.
func formatIssueRelationsText(relations []redmine.IssueRelation) error {
	if len(relations) == 0 {
		fmt.Println("チケット関連が見つかりませんでした。")
		return nil
	}

	for _, r := range relations {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(r.ID)))
		fmt.Println(formatter.FormatKeyValue("Issue ID", strconv.Itoa(r.IssueID)))
		fmt.Println(formatter.FormatKeyValue("Issue To ID", strconv.Itoa(r.IssueToID)))
		fmt.Println(formatter.FormatKeyValue("Type", r.RelationType))
		if r.Delay != 0 {
			fmt.Println(formatter.FormatKeyValue("Delay", strconv.Itoa(r.Delay)))
		}
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(issueRelationCmd)

	// Subcommands
	issueRelationCmd.AddCommand(issueRelationListCmd)
	issueRelationCmd.AddCommand(issueRelationGetCmd)
	issueRelationCmd.AddCommand(issueRelationCreateCmd)
	issueRelationCmd.AddCommand(issueRelationDeleteCmd)

	// Flags for list command
	issueRelationListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")

	// Flags for get command
	issueRelationGetCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")

	// Flags for create command
	issueRelationCreateCmd.Flags().Int("issue-to-id", 0, "関連先チケットID (必須)")
	issueRelationCreateCmd.Flags().String("relation-type", "", "関連タイプ (必須: relates, duplicates, duplicated, blocks, blocked, precedes, follows, copied_to, copied_from)")
	issueRelationCreateCmd.Flags().Int("delay", 0, "遅延日数 (precedesまたはfollowsの場合)")
}
