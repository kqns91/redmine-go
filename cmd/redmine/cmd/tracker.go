package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var trackerCmd = &cobra.Command{
	Use:   "tracker",
	Short: "Manage Redmine trackers",
	Long:  `トラッカーの取得などの操作を行います。`,
}

var trackerListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all trackers",
	Long:  `すべてのトラッカーをリスト表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")

		result, err := client.ListTrackers(context.Background())
		if err != nil {
			return fmt.Errorf("トラッカーの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatTrackersTable(result.Trackers)
		case formatText:
			return formatTrackersText(result.Trackers)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

// formatTrackersTable formats trackers in table format.
func formatTrackersTable(trackers []redmine.Tracker) error {
	if len(trackers) == 0 {
		fmt.Println("トラッカーが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Name", "Default Status", "Description"}
	rows := make([][]string, 0, len(trackers))

	for _, t := range trackers {
		defaultStatus := "-"
		if t.DefaultStatus.Name != "" {
			defaultStatus = t.DefaultStatus.Name
		}

		rows = append(rows, []string{
			strconv.Itoa(t.ID),
			formatter.TruncateString(t.Name, 20),
			formatter.TruncateString(defaultStatus, 20),
			formatter.TruncateString(t.Description, 40),
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatTrackersText formats trackers in simple text format.
func formatTrackersText(trackers []redmine.Tracker) error {
	if len(trackers) == 0 {
		fmt.Println("トラッカーが見つかりませんでした。")
		return nil
	}

	for _, t := range trackers {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(t.ID)))
		fmt.Println(formatter.FormatKeyValue("Name", t.Name))
		if t.DefaultStatus.Name != "" {
			fmt.Println(formatter.FormatKeyValue("Default Status", t.DefaultStatus.Name))
		}
		if t.Description != "" {
			fmt.Println(formatter.FormatKeyValue("Description", formatter.TruncateString(t.Description, 80)))
		}
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(trackerCmd)

	// Subcommands
	trackerCmd.AddCommand(trackerListCmd)

	// Flags for list command
	trackerListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")
}
