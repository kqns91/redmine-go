package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

const (
	statusActive   = "Active"
	statusClosed   = "Closed"
	statusArchived = "Archived"
)

var projectCmd = &cobra.Command{
	Use:   "project",
	Short: "Manage Redmine projects",
	Long:  `プロジェクトの作成、取得、更新、削除などの操作を行います。`,
}

var projectListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all projects",
	Long:  `すべてのプロジェクトをリスト表示します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		include, _ := cmd.Flags().GetString("include")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		format, _ := cmd.Flags().GetString("format")

		opts := &redmine.ListProjectsOptions{
			Include: include,
			Limit:   limit,
			Offset:  offset,
		}

		result, err := client.ListProjects(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("プロジェクトの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatProjectsTable(result.Projects)
		case formatText:
			return formatProjectsText(result.Projects)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

var projectGetCmd = &cobra.Command{
	Use:   "get [project_id_or_identifier]",
	Short: "Get a project by ID or identifier",
	Long:  `指定したIDまたは識別子のプロジェクトを取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		include, _ := cmd.Flags().GetString("include")
		format, _ := cmd.Flags().GetString("format")

		opts := &redmine.ShowProjectOptions{
			Include: include,
		}

		result, err := client.ShowProject(context.Background(), args[0], opts)
		if err != nil {
			return fmt.Errorf("プロジェクトの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatProjectDetail(&result.Project)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
	},
}

var projectCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new project",
	Long:  `新しいプロジェクトを作成します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		identifier, _ := cmd.Flags().GetString("identifier")
		description, _ := cmd.Flags().GetString("description")
		homepage, _ := cmd.Flags().GetString("homepage")
		isPublic, _ := cmd.Flags().GetBool("public")
		inheritMembers, _ := cmd.Flags().GetBool("inherit-members")
		parentID, _ := cmd.Flags().GetInt("parent-id")

		if name == "" {
			return errors.New("--name フラグは必須です")
		}
		if identifier == "" {
			return errors.New("--identifier フラグは必須です")
		}

		req := redmine.ProjectCreateRequest{
			Name:           name,
			Identifier:     identifier,
			Description:    description,
			Homepage:       homepage,
			IsPublic:       isPublic,
			InheritMembers: inheritMembers,
			ParentID:       parentID,
		}

		result, err := client.CreateProject(context.Background(), req)
		if err != nil {
			return fmt.Errorf("プロジェクトの作成に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var projectUpdateCmd = &cobra.Command{
	Use:   "update [project_id_or_identifier]",
	Short: "Update an existing project",
	Long:  `既存のプロジェクトを更新します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, _ := cmd.Flags().GetString("name")
		description, _ := cmd.Flags().GetString("description")
		homepage, _ := cmd.Flags().GetString("homepage")
		isPublic, _ := cmd.Flags().GetBool("public")
		inheritMembers, _ := cmd.Flags().GetBool("inherit-members")
		parentID, _ := cmd.Flags().GetInt("parent-id")

		req := redmine.ProjectUpdateRequest{
			Name:           name,
			Description:    description,
			Homepage:       homepage,
			IsPublic:       isPublic,
			InheritMembers: inheritMembers,
			ParentID:       parentID,
		}

		err := client.UpdateProject(context.Background(), args[0], req)
		if err != nil {
			return fmt.Errorf("プロジェクトの更新に失敗しました: %w", err)
		}

		fmt.Println("プロジェクトを更新しました")
		return nil
	},
}

var projectDeleteCmd = &cobra.Command{
	Use:   "delete [project_id_or_identifier]",
	Short: "Delete a project",
	Long:  `プロジェクトを削除します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := client.DeleteProject(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("プロジェクトの削除に失敗しました: %w", err)
		}

		fmt.Println("プロジェクトを削除しました")
		return nil
	},
}

var projectArchiveCmd = &cobra.Command{
	Use:   "archive [project_id_or_identifier]",
	Short: "Archive a project",
	Long:  `プロジェクトをアーカイブします (Redmine 5.0以降)。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := client.ArchiveProject(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("プロジェクトのアーカイブに失敗しました: %w", err)
		}

		fmt.Println("プロジェクトをアーカイブしました")
		return nil
	},
}

var projectUnarchiveCmd = &cobra.Command{
	Use:   "unarchive [project_id_or_identifier]",
	Short: "Unarchive a project",
	Long:  `プロジェクトのアーカイブを解除します (Redmine 5.0以降)。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		err := client.UnarchiveProject(context.Background(), args[0])
		if err != nil {
			return fmt.Errorf("プロジェクトのアーカイブ解除に失敗しました: %w", err)
		}

		fmt.Println("プロジェクトのアーカイブを解除しました")
		return nil
	},
}

// getProjectStatus returns a human-readable status string.
func getProjectStatus(status int) string {
	switch status {
	case 5:
		return statusClosed
	case 9:
		return statusArchived
	default:
		return statusActive
	}
}

// formatProjectDetail formats a single project in detailed text format.
func formatProjectDetail(p *redmine.Project) error {
	// Title
	fmt.Println(formatter.FormatTitle("Project: " + p.Name))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(p.ID)))
	fmt.Println(formatter.FormatKeyValue("Identifier", p.Identifier))
	fmt.Println(formatter.FormatKeyValue("Name", p.Name))

	status := getProjectStatus(p.Status)
	fmt.Println(formatter.FormatKeyValue("Status", status))

	if p.Description != "" {
		fmt.Println(formatter.FormatKeyValue("Description", p.Description))
	}
	if p.Homepage != "" {
		fmt.Println(formatter.FormatKeyValue("Homepage", p.Homepage))
	}
	fmt.Println(formatter.FormatKeyValue("Public", strconv.FormatBool(p.IsPublic)))
	fmt.Println(formatter.FormatKeyValue("Created", p.CreatedOn))
	fmt.Println(formatter.FormatKeyValue("Updated", p.UpdatedOn))

	// Parent Project
	if p.Parent.ID != 0 {
		fmt.Println()
		fmt.Println(formatter.FormatSection("親プロジェクト"))
		fmt.Println(formatter.FormatKeyValue("Name", p.Parent.Name))
	}

	return nil
}

// formatProjectsTable formats projects in table format.
func formatProjectsTable(projects []redmine.Project) error {
	if len(projects) == 0 {
		fmt.Println("プロジェクトが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Identifier", "Name", "Status", "Created"}
	rows := make([][]string, 0, len(projects))

	for _, p := range projects {
		rows = append(rows, []string{
			strconv.Itoa(p.ID),
			p.Identifier,
			formatter.TruncateString(p.Name, 40),
			getProjectStatus(p.Status),
			p.CreatedOn,
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatProjectsText formats projects in simple text format.
func formatProjectsText(projects []redmine.Project) error {
	if len(projects) == 0 {
		fmt.Println("プロジェクトが見つかりませんでした。")
		return nil
	}

	for _, p := range projects {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(p.ID)))
		fmt.Println(formatter.FormatKeyValue("Identifier", p.Identifier))
		fmt.Println(formatter.FormatKeyValue("Name", p.Name))
		fmt.Println(formatter.FormatKeyValue("Status", getProjectStatus(p.Status)))
		if p.Description != "" {
			fmt.Println(formatter.FormatKeyValue("Description", formatter.TruncateString(p.Description, 80)))
		}
		fmt.Println()
	}

	return nil
}

func init() {
	rootCmd.AddCommand(projectCmd)

	// Subcommands
	projectCmd.AddCommand(projectListCmd)
	projectCmd.AddCommand(projectGetCmd)
	projectCmd.AddCommand(projectCreateCmd)
	projectCmd.AddCommand(projectUpdateCmd)
	projectCmd.AddCommand(projectDeleteCmd)
	projectCmd.AddCommand(projectArchiveCmd)
	projectCmd.AddCommand(projectUnarchiveCmd)

	// Flags for list command
	projectListCmd.Flags().String("include", "", "追加で取得する情報 (例: trackers,issue_categories)")
	projectListCmd.Flags().Int("limit", 0, "取得する最大件数")
	projectListCmd.Flags().Int("offset", 0, "取得開始位置のオフセット")
	projectListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")

	// Flags for get command
	projectGetCmd.Flags().String("include", "", "追加で取得する情報 (例: trackers,issue_categories)")
	projectGetCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")

	// Flags for create command
	projectCreateCmd.Flags().String("name", "", "プロジェクト名 (必須)")
	projectCreateCmd.Flags().String("identifier", "", "プロジェクト識別子 (必須)")
	projectCreateCmd.Flags().String("description", "", "プロジェクトの説明")
	projectCreateCmd.Flags().String("homepage", "", "ホームページURL")
	projectCreateCmd.Flags().Bool("public", true, "公開プロジェクトにするかどうか")
	projectCreateCmd.Flags().Bool("inherit-members", false, "親プロジェクトのメンバーを継承するかどうか")
	projectCreateCmd.Flags().Int("parent-id", 0, "親プロジェクトID")

	// Flags for update command
	projectUpdateCmd.Flags().String("name", "", "プロジェクト名")
	projectUpdateCmd.Flags().String("description", "", "プロジェクトの説明")
	projectUpdateCmd.Flags().String("homepage", "", "ホームページURL")
	projectUpdateCmd.Flags().Bool("public", false, "公開プロジェクトにするかどうか")
	projectUpdateCmd.Flags().Bool("inherit-members", false, "親プロジェクトのメンバーを継承するかどうか")
	projectUpdateCmd.Flags().Int("parent-id", 0, "親プロジェクトID")
}
