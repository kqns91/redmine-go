package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var attachmentCmd = &cobra.Command{
	Use:   "attachment",
	Short: "Manage Redmine attachments",
	Long:  `添付ファイルの取得、更新、削除などの操作を行います。`,
}

var attachmentShowCmd = &cobra.Command{
	Use:   "show [attachment_id]",
	Short: "Show an attachment by ID",
	Long:  `指定したIDの添付ファイルを取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なattachment_id: %w", err)
		}

		format, _ := cmd.Flags().GetString("format")

		result, err := client.ShowAttachment(context.Background(), id)
		if err != nil {
			return fmt.Errorf("添付ファイルの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatAttachmentDetail(&result.Attachment)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
	},
}

var attachmentUpdateCmd = &cobra.Command{
	Use:   "update [attachment_id]",
	Short: "Update an attachment",
	Long:  `添付ファイルの情報を更新します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なattachment_id: %w", err)
		}

		filename, _ := cmd.Flags().GetString("filename")
		description, _ := cmd.Flags().GetString("description")

		attachment := redmine.Attachment{
			Filename:    filename,
			Description: description,
		}

		err = client.UpdateAttachment(context.Background(), id, attachment)
		if err != nil {
			return fmt.Errorf("添付ファイルの更新に失敗しました: %w", err)
		}

		fmt.Println("添付ファイルを更新しました")
		return nil
	},
}

var attachmentDeleteCmd = &cobra.Command{
	Use:   "delete [attachment_id]",
	Short: "Delete an attachment",
	Long:  `添付ファイルを削除します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なattachment_id: %w", err)
		}

		err = client.DeleteAttachment(context.Background(), id)
		if err != nil {
			return fmt.Errorf("添付ファイルの削除に失敗しました: %w", err)
		}

		fmt.Println("添付ファイルを削除しました")
		return nil
	},
}

// formatAttachmentDetail formats a single attachment in detailed text format.
func formatAttachmentDetail(a *redmine.Attachment) error {
	// Title
	fmt.Println(formatter.FormatTitle("Attachment: " + a.Filename))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(a.ID)))
	fmt.Println(formatter.FormatKeyValue("Filename", a.Filename))
	fmt.Println(formatter.FormatKeyValue("Filesize", strconv.Itoa(a.Filesize)))
	fmt.Println(formatter.FormatKeyValue("Content-Type", a.ContentType))
	if a.Description != "" {
		fmt.Println(formatter.FormatKeyValue("Description", a.Description))
	}
	if a.ContentURL != "" {
		fmt.Println(formatter.FormatKeyValue("Content URL", a.ContentURL))
	}
	fmt.Println(formatter.FormatKeyValue("Created", a.CreatedOn))

	// Author Info
	if a.Author.ID != 0 {
		fmt.Println()
		fmt.Println(formatter.FormatSection("作成者"))
		fmt.Println(formatter.FormatKeyValue("Name", a.Author.Name))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(attachmentCmd)

	// Subcommands
	attachmentCmd.AddCommand(attachmentShowCmd)
	attachmentCmd.AddCommand(attachmentUpdateCmd)
	attachmentCmd.AddCommand(attachmentDeleteCmd)

	// Flags for show command
	attachmentShowCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")

	// Flags for update command
	attachmentUpdateCmd.Flags().String("filename", "", "ファイル名")
	attachmentUpdateCmd.Flags().String("description", "", "説明")
}
