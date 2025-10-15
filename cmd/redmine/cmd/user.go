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

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Manage Redmine users",
	Long:  `ユーザーの作成、取得、更新、削除などの操作を行います（多くの操作には管理者権限が必要です）。`,
}

var userListCmd = &cobra.Command{
	Use:   "list",
	Short: "List users",
	Long:  `ユーザーをリスト表示します（管理者権限が必要です）。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		status, _ := cmd.Flags().GetString("status")
		name, _ := cmd.Flags().GetString("name")
		groupID, _ := cmd.Flags().GetInt("group-id")
		include, _ := cmd.Flags().GetString("include")
		limit, _ := cmd.Flags().GetInt("limit")
		offset, _ := cmd.Flags().GetInt("offset")
		format, _ := cmd.Flags().GetString("format")

		opts := &redmine.ListUsersOptions{
			Status:  status,
			Name:    name,
			GroupID: groupID,
			Include: include,
			Limit:   limit,
			Offset:  offset,
		}

		result, err := client.ListUsers(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("ユーザーの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatTable:
			return formatUsersTable(result.Users)
		case formatText:
			return formatUsersText(result.Users)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s", format)
		}
	},
}

var userGetCmd = &cobra.Command{
	Use:   "get [user_id]",
	Short: "Get a user by ID",
	Long:  `指定したIDのユーザーを取得します。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なuser_id: %w", err)
		}

		include, _ := cmd.Flags().GetString("include")
		format, _ := cmd.Flags().GetString("format")

		opts := &redmine.ShowUserOptions{
			Include: include,
		}

		result, err := client.ShowUser(context.Background(), id, opts)
		if err != nil {
			return fmt.Errorf("ユーザーの取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatUserDetail(&result.User)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
	},
}

var userCurrentCmd = &cobra.Command{
	Use:   "current",
	Short: "Get the current user",
	Long:  `現在認証されているユーザーの情報を取得します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		include, _ := cmd.Flags().GetString("include")

		opts := &redmine.ShowUserOptions{
			Include: include,
		}

		result, err := client.GetCurrentUser(context.Background(), opts)
		if err != nil {
			return fmt.Errorf("現在のユーザーの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var userCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new user",
	Long:  `新しいユーザーを作成します（管理者権限が必要です）。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		login, _ := cmd.Flags().GetString("login")
		firstname, _ := cmd.Flags().GetString("firstname")
		lastname, _ := cmd.Flags().GetString("lastname")
		mail, _ := cmd.Flags().GetString("mail")
		password, _ := cmd.Flags().GetString("password")
		authSourceID, _ := cmd.Flags().GetInt("auth-source-id")
		mailNotification, _ := cmd.Flags().GetString("mail-notification")
		mustChangePasswd, _ := cmd.Flags().GetBool("must-change-passwd")
		generatePassword, _ := cmd.Flags().GetBool("generate-password")

		if login == "" {
			return errors.New("--login フラグは必須です")
		}
		if firstname == "" {
			return errors.New("--firstname フラグは必須です")
		}
		if lastname == "" {
			return errors.New("--lastname フラグは必須です")
		}
		if mail == "" {
			return errors.New("--mail フラグは必須です")
		}

		user := redmine.User{
			Login:             login,
			Firstname:         firstname,
			Lastname:          lastname,
			Mail:              mail,
			Password:          password,
			AuthSourceID:      authSourceID,
			MailNotification:  mailNotification,
			MustChangePasswd:  mustChangePasswd,
			GeneratePassword:  generatePassword,
		}

		result, err := client.CreateUser(context.Background(), user)
		if err != nil {
			return fmt.Errorf("ユーザーの作成に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var userUpdateCmd = &cobra.Command{
	Use:   "update [user_id]",
	Short: "Update an existing user",
	Long:  `既存のユーザーを更新します（管理者権限が必要です）。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なuser_id: %w", err)
		}

		login, _ := cmd.Flags().GetString("login")
		firstname, _ := cmd.Flags().GetString("firstname")
		lastname, _ := cmd.Flags().GetString("lastname")
		mail, _ := cmd.Flags().GetString("mail")
		password, _ := cmd.Flags().GetString("password")
		authSourceID, _ := cmd.Flags().GetInt("auth-source-id")
		mailNotification, _ := cmd.Flags().GetString("mail-notification")
		mustChangePasswd, _ := cmd.Flags().GetBool("must-change-passwd")
		generatePassword, _ := cmd.Flags().GetBool("generate-password")

		user := redmine.User{
			Login:             login,
			Firstname:         firstname,
			Lastname:          lastname,
			Mail:              mail,
			Password:          password,
			AuthSourceID:      authSourceID,
			MailNotification:  mailNotification,
			MustChangePasswd:  mustChangePasswd,
			GeneratePassword:  generatePassword,
		}

		err = client.UpdateUser(context.Background(), id, user)
		if err != nil {
			return fmt.Errorf("ユーザーの更新に失敗しました: %w", err)
		}

		fmt.Println("ユーザーを更新しました")
		return nil
	},
}

var userDeleteCmd = &cobra.Command{
	Use:   "delete [user_id]",
	Short: "Delete a user",
	Long:  `ユーザーを削除します（管理者権限が必要です）。`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("無効なuser_id: %w", err)
		}

		err = client.DeleteUser(context.Background(), id)
		if err != nil {
			return fmt.Errorf("ユーザーの削除に失敗しました: %w", err)
		}

		fmt.Println("ユーザーを削除しました")
		return nil
	},
}

// getUserStatus returns a human-readable status string.
func getUserStatus(status int) string {
	switch status {
	case 1:
		return "Active"
	case 2:
		return "Registered"
	case 3:
		return "Locked"
	default:
		return "Unknown"
	}
}

// formatUserDetail formats a single user in detailed text format.
func formatUserDetail(u *redmine.User) error {
	// Title
	fmt.Println(formatter.FormatTitle("User: " + u.Login))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(u.ID)))
	fmt.Println(formatter.FormatKeyValue("Login", u.Login))
	fmt.Println(formatter.FormatKeyValue("Firstname", u.Firstname))
	fmt.Println(formatter.FormatKeyValue("Lastname", u.Lastname))
	fmt.Println(formatter.FormatKeyValue("Mail", u.Mail))
	if u.Status > 0 {
		fmt.Println(formatter.FormatKeyValue("Status", getUserStatus(u.Status)))
	}
	fmt.Println(formatter.FormatKeyValue("Admin", strconv.FormatBool(u.Admin)))

	// Timestamps
	fmt.Println()
	fmt.Println(formatter.FormatSection("タイムスタンプ"))
	if u.CreatedOn != "" {
		fmt.Println(formatter.FormatKeyValue("Created", u.CreatedOn))
	}
	if u.UpdatedOn != "" {
		fmt.Println(formatter.FormatKeyValue("Updated", u.UpdatedOn))
	}
	if u.LastLoginOn != "" {
		fmt.Println(formatter.FormatKeyValue("Last Login", u.LastLoginOn))
	}

	return nil
}

// formatUsersTable formats users in table format.
func formatUsersTable(users []redmine.User) error {
	if len(users) == 0 {
		fmt.Println("ユーザーが見つかりませんでした。")
		return nil
	}

	headers := []string{"ID", "Login", "Name", "Mail", "Status", "Admin"}
	rows := make([][]string, 0, len(users))

	for _, u := range users {
		name := u.Firstname + " " + u.Lastname
		statusStr := "-"
		if u.Status > 0 {
			statusStr = getUserStatus(u.Status)
		}

		rows = append(rows, []string{
			strconv.Itoa(u.ID),
			formatter.TruncateString(u.Login, 15),
			formatter.TruncateString(name, 25),
			formatter.TruncateString(u.Mail, 30),
			statusStr,
			strconv.FormatBool(u.Admin),
		})
	}

	formatter.RenderTable(headers, rows)
	return nil
}

// formatUsersText formats users in simple text format.
func formatUsersText(users []redmine.User) error {
	if len(users) == 0 {
		fmt.Println("ユーザーが見つかりませんでした。")
		return nil
	}

	for _, u := range users {
		fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(u.ID)))
		fmt.Println(formatter.FormatKeyValue("Login", u.Login))
		fmt.Println(formatter.FormatKeyValue("Name", u.Firstname+" "+u.Lastname))
		fmt.Println(formatter.FormatKeyValue("Mail", u.Mail))
		if u.Status > 0 {
			fmt.Println(formatter.FormatKeyValue("Status", getUserStatus(u.Status)))
		}
		fmt.Println(formatter.FormatKeyValue("Admin", strconv.FormatBool(u.Admin)))
		fmt.Println()
	}

	return nil
}


// includeOptionsForUser returns valid include options for user commands
func includeOptionsForUser() []string {
	return []string{"memberships", "groups"}
}

func init() {
	rootCmd.AddCommand(userCmd)

	// Subcommands
	userCmd.AddCommand(userListCmd)
	userCmd.AddCommand(userGetCmd)
	userCmd.AddCommand(userCurrentCmd)
	userCmd.AddCommand(userCreateCmd)
	userCmd.AddCommand(userUpdateCmd)
	userCmd.AddCommand(userDeleteCmd)

	// Flags for list command
	userListCmd.Flags().String("status", "", "ステータス (active, locked, registered)")
	userListCmd.Flags().String("name", "", "ユーザー名")
	userListCmd.Flags().Int("group-id", 0, "グループID")
	userListCmd.Flags().String("include", "", "追加で取得する情報 (memberships, groups)")
	userListCmd.Flags().Int("limit", 0, "取得する最大件数")
	userListCmd.Flags().Int("offset", 0, "取得開始位置のオフセット")
	userListCmd.Flags().StringP("format", "f", formatTable, "出力フォーマット (json, table, text)")

	// Register flag completion for list command
	_ = userListCmd.RegisterFlagCompletionFunc("include", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return includeOptionsForUser(), cobra.ShellCompDirectiveNoFileComp
	})

	// Flags for get command
	userGetCmd.Flags().String("include", "", "追加で取得する情報 (memberships, groups)")
	userGetCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")

	// Register flag completion for get command
	_ = userGetCmd.RegisterFlagCompletionFunc("include", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return includeOptionsForUser(), cobra.ShellCompDirectiveNoFileComp
	})

	// Flags for current command
	userCurrentCmd.Flags().String("include", "", "追加で取得する情報 (memberships, groups)")

	// Register flag completion for current command
	_ = userCurrentCmd.RegisterFlagCompletionFunc("include", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return includeOptionsForUser(), cobra.ShellCompDirectiveNoFileComp
	})

	// Flags for create command
	userCreateCmd.Flags().String("login", "", "ログインID (必須)")
	userCreateCmd.Flags().String("firstname", "", "名 (必須)")
	userCreateCmd.Flags().String("lastname", "", "姓 (必須)")
	userCreateCmd.Flags().String("mail", "", "メールアドレス (必須)")
	userCreateCmd.Flags().String("password", "", "パスワード")
	userCreateCmd.Flags().Int("auth-source-id", 0, "認証ソースID")
	userCreateCmd.Flags().String("mail-notification", "", "メール通知設定 (all, selected, only_my_events, only_assigned, only_owner, none)")
	userCreateCmd.Flags().Bool("must-change-passwd", false, "初回ログイン時にパスワード変更を強制")
	userCreateCmd.Flags().Bool("generate-password", false, "パスワードを自動生成")

	// Flags for update command
	userUpdateCmd.Flags().String("login", "", "ログインID")
	userUpdateCmd.Flags().String("firstname", "", "名")
	userUpdateCmd.Flags().String("lastname", "", "姓")
	userUpdateCmd.Flags().String("mail", "", "メールアドレス")
	userUpdateCmd.Flags().String("password", "", "パスワード")
	userUpdateCmd.Flags().Int("auth-source-id", 0, "認証ソースID")
	userUpdateCmd.Flags().String("mail-notification", "", "メール通知設定 (all, selected, only_my_events, only_assigned, only_owner, none)")
	userUpdateCmd.Flags().Bool("must-change-passwd", false, "次回ログイン時にパスワード変更を強制")
	userUpdateCmd.Flags().Bool("generate-password", false, "パスワードを自動生成")
}
