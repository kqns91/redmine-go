package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

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

		result, err := client.ShowAttachment(context.Background(), id)
		if err != nil {
			return fmt.Errorf("添付ファイルの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
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

func init() {
	rootCmd.AddCommand(attachmentCmd)

	// Subcommands
	attachmentCmd.AddCommand(attachmentShowCmd)
	attachmentCmd.AddCommand(attachmentUpdateCmd)
	attachmentCmd.AddCommand(attachmentDeleteCmd)

	// Flags for update command
	attachmentUpdateCmd.Flags().String("filename", "", "ファイル名")
	attachmentUpdateCmd.Flags().String("description", "", "説明")
}
