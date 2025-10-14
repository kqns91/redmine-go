package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/cli/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Manage Redmine versions",
	Long:  `プロジェクトバージョンの作成、取得、更新、削除などの操作を行います。`,
}

var versionListCmd = &cobra.Command{
	Use:   "list [project_id_or_identifier]",
	Short: "List versions for a project",
	Long:  `指定したプロジェクトのバージョンをリスト表示します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")

		result, err := client.ListVersions(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("バージョンの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatVersionsTable(result.Versions)
		case formatText:
			return formatVersionsText(result.Versions)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

var versionGetCmd = &cobra.Command{
	Use:   "get [version_id]",
	Short: "Get a version by ID",
	Long:  `指定したIDのバージョンを取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なversion_id: %w", err)
		}

		result, err := client.ShowVersion(context.Background(), id)
		if err != nil {
			return fmt.Errorf("バージョンの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var versionCreateCmd = &cobra.Command{
	Use:   "create [project_id_or_identifier]",
	Short: "Create a new version",
	Long:  `指定したプロジェクトに新しいバージョンを作成します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")
		dueDate, _ := cmd.Flags().GetString("due-date")
		sharing, _ := cmd.Flags().GetString("sharing")
		wikiPageTitle, _ := cmd.Flags().GetString("wiki-page-title")

		if name == "" {
			return errors.New("--name フラグは必須です")
		}

		version := redmine.Version{
			Name:          name,
			Description:   description,
			Status:        status,
			DueDate:       dueDate,
			Sharing:       sharing,
			WikiPageTitle: wikiPageTitle,
		}

		result, err := client.CreateVersion(context.Background(), args[0], version)
		if err != nil {
			return fmt.Errorf("バージョンの作成に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var versionUpdateCmd = &cobra.Command{
	Use:   "update [version_id]",
	Short: "Update an existing version",
	Long:  `既存のバージョンを更新します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なversion_id: %w", err)
		}

		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		status, _ := cmd.Flags().GetString("status")
		dueDate, _ := cmd.Flags().GetString("due-date")
		sharing, _ := cmd.Flags().GetString("sharing")
		wikiPageTitle, _ := cmd.Flags().GetString("wiki-page-title")

		version := redmine.Version{
			Name:          name,
			Description:   description,
			Status:        status,
			DueDate:       dueDate,
			Sharing:       sharing,
			WikiPageTitle: wikiPageTitle,
		}

		err = client.UpdateVersion(context.Background(), id, version)
		if err != nil {
			return fmt.Errorf("バージョンの更新に失敗しました: %w", err)
		}

		fmt.Println("バージョンを更新しました")
		return nil
	},
}

var versionDeleteCmd = &cobra.Command{
	Use:   "delete [version_id]",
	Short: "Delete a version",
	Long:  `バージョンを削除します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なversion_id: %w", err)
		}

		err = client.DeleteVersion(context.Background(), id)
		if err != nil {
			return fmt.Errorf("バージョンの削除に失敗しました: %w", err)
		}

		fmt.Println("バージョンを削除しました")
		return nil
	},
}

// formatVersionsTable formats versions in table format.
func formatVersionsTable(versions []redmine.Version) error {
	if len(versions) == 0 {
		fmt.Println("バージョンが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Name", "Status", "Due Date", "Created"}
	rows := make([][]string, 0, len(versions))

	for _, v := range versions {
		dueDate := "-"
		if v.DueDate != "" {
			dueDate = v.DueDate
		}
		status := v.Status
		if status == "" {
			status = "open"
		}

		rows = append(rows, []string{
			strconv.Itoa(v.ID),
			formatter.TruncateString(v.Name, 30),
			status,
			dueDate,
			v.CreatedOn,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatVersionsText formats versions in simple text format.
func formatVersionsText(versions []redmine.Version) error {
	if len(versions) == 0 {
		fmt.Println("バージョンが見つかりませんでした。")
		return nil
	}

	for _, v := range versions {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(v.ID)))
		fmt.Println(formatter.FormatKeyValue("Name", v.Name))
		status := v.Status
		if status == "" {
			status = "open"
		}
		fmt.Println(formatter.FormatKeyValue("Status", status))
		if v.Description != "" {
			fmt.Println(formatter.FormatKeyValue("Description", formatter.TruncateString(v.Description, 80)))
		}
		if v.DueDate != "" {
			fmt.Println(formatter.FormatKeyValue("Due Date", v.DueDate))
		}
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(versionCmd)

	// Subcommands
	versionCmd.AddCommand(versionListCmd)
	versionCmd.AddCommand(versionGetCmd)
	versionCmd.AddCommand(versionCreateCmd)
	versionCmd.AddCommand(versionUpdateCmd)
	versionCmd.AddCommand(versionDeleteCmd)

	// Flags for list command
	versionListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")

	// Flags for create command
	versionCreateCmd.Flags().String("name", "", "バージョン名 (必須)")
	versionCreateCmd.Flags().String("description", "", "説明")
	versionCreateCmd.Flags().String("status", "", "ステータス (open, locked, closed)")
	versionCreateCmd.Flags().String("due-date", "", "期日 (YYYY-MM-DD)")
	versionCreateCmd.Flags().String("sharing", "", "共有設定 (none, descendants, hierarchy, tree, system)")
	versionCreateCmd.Flags().String("wiki-page-title", "", "Wikiページタイトル")

	// Flags for update command
	versionUpdateCmd.Flags().String("name", "", "バージョン名")
	versionUpdateCmd.Flags().String("description", "", "説明")
	versionUpdateCmd.Flags().String("status", "", "ステータス (open, locked, closed)")
	versionUpdateCmd.Flags().String("due-date", "", "期日 (YYYY-MM-DD)")
	versionUpdateCmd.Flags().String("sharing", "", "共有設定 (none, descendants, hierarchy, tree, system)")
	versionUpdateCmd.Flags().String("wiki-page-title", "", "Wikiページタイトル")
}
