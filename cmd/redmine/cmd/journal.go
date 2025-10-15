package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var journalCmd = &cobra.Command{
	Use:   "journal",
	Short: "Manage Redmine journals",
	Long:  `ジャーナル（チケット履歴）の取得などの操作を行います。`,
}

var journalGetCmd = &cobra.Command{
	Use:   "get [journal_id]",
	Short: "Get a journal by ID",
	Long:  `指定したIDのジャーナルを取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なjournal_id: %w", err)
		}

		format, _ := cmd.Flags().GetString("format")

		result, err := client.ShowJournal(context.Background(), id)
		if err != nil {
			return fmt.Errorf("ジャーナルの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatJournalDetail(&result.Journal)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
	},
}

// formatJournalDetail formats a single journal in detailed text format.
func formatJournalDetail(j *redmine.Journal) error {
	// Title
	fmt.Println(formatter.FormatTitle("Journal #" + strconv.Itoa(j.ID)))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(j.ID)))
	if j.User.Name != "" {
		fmt.Println(formatter.FormatKeyValue("User", j.User.Name))
	}
	if j.Notes != "" {
		fmt.Println(formatter.FormatKeyValue("Notes", j.Notes))
	}
	fmt.Println(formatter.FormatKeyValue("Created", j.CreatedOn))

	// Details
	if len(j.Details) > 0 {
		fmt.Println()
		fmt.Println(formatter.FormatSection("変更内容"))
		for _, detail := range j.Details {
			if detail.Property != "" {
				fmt.Printf("  %s.%s: %s -> %s\n", detail.Property, detail.Name, detail.OldValue, detail.NewValue)
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(journalCmd)

	// Subcommands
	journalCmd.AddCommand(journalGetCmd)

	// Flags for get command
	journalGetCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")
}
