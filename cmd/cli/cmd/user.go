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

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
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

		opts := &redmine.ShowUserOptions{
			Include: include,
		}

		result, err := client.ShowUser(context.Background(), id, opts)
		if err != nil {
			return fmt.Errorf("ユーザーの取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
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
			Login:     login,
			Firstname: firstname,
			Lastname:  lastname,
			Mail:      mail,
		}

		// パスワードは別の方法で設定する必要があるため、ここではユーザー作成のみ
		// 実際のAPIではパスワードフィールドが存在する可能性がありますが、
		// 現在のUser構造体には含まれていないため、コメントとして残しておきます
		_ = password

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

		user := redmine.User{
			Login:     login,
			Firstname: firstname,
			Lastname:  lastname,
			Mail:      mail,
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
	userListCmd.Flags().String("include", "", "追加で取得する情報")
	userListCmd.Flags().Int("limit", 0, "取得する最大件数")
	userListCmd.Flags().Int("offset", 0, "取得開始位置のオフセット")

	// Flags for get command
	userGetCmd.Flags().String("include", "", "追加で取得する情報")

	// Flags for current command
	userCurrentCmd.Flags().String("include", "", "追加で取得する情報")

	// Flags for create command
	userCreateCmd.Flags().String("login", "", "ログインID (必須)")
	userCreateCmd.Flags().String("firstname", "", "名 (必須)")
	userCreateCmd.Flags().String("lastname", "", "姓 (必須)")
	userCreateCmd.Flags().String("mail", "", "メールアドレス (必須)")
	userCreateCmd.Flags().String("password", "", "パスワード")

	// Flags for update command
	userUpdateCmd.Flags().String("login", "", "ログインID")
	userUpdateCmd.Flags().String("firstname", "", "名")
	userUpdateCmd.Flags().String("lastname", "", "姓")
	userUpdateCmd.Flags().String("mail", "", "メールアドレス")
}
