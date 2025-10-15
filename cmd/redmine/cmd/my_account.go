package cmd

import (
	"context"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/kqns91/redmine-go/cmd/redmine/internal/formatter"
	"github.com/kqns91/redmine-go/pkg/redmine"
)

var myAccountCmd = &cobra.Command{
	Use:   "my-account",
	Short: "Manage current user's account",
	Long:  `現在のユーザーのアカウント情報を取得・更新します。`,
}

var myAccountGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current user's account",
	Long:  `現在のユーザーのアカウント情報を取得します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")

		result, err := client.GetMyAccount(context.Background())
		if err != nil {
			return fmt.Errorf("アカウント情報の取得に失敗しました: %w", err)
		}

		// Format output based on --format flag
		switch format {
		case formatJSON:
			return formatter.OutputJSON(result)
		case formatText:
			return formatMyAccountDetail(&result.User)
		default:
			return fmt.Errorf("不明な出力フォーマット: %s (利用可能: json, text)", format)
		}
	},
}

var myAccountUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update current user's account",
	Long:  `現在のユーザーのアカウント情報を更新します。`,
	RunE: func(cmd *cobra.Command, args []string) error {
		firstname, _ := cmd.Flags().GetString("firstname")
		lastname, _ := cmd.Flags().GetString("lastname")
		mail, _ := cmd.Flags().GetString("mail")

		user := redmine.User{
			Firstname: firstname,
			Lastname:  lastname,
			Mail:      mail,
		}

		err := client.UpdateMyAccount(context.Background(), user)
		if err != nil {
			return fmt.Errorf("アカウント情報の更新に失敗しました: %w", err)
		}

		fmt.Println("アカウント情報を更新しました")
		return nil
	},
}

// formatMyAccountDetail formats current user account in detailed text format.
func formatMyAccountDetail(u *redmine.User) error {
	// Title
	fullName := u.Firstname + " " + u.Lastname
	fmt.Println(formatter.FormatTitle("My Account: " + fullName))
	fmt.Println()

	// Basic Info
	fmt.Println(formatter.FormatSection("基本情報"))
	fmt.Println(formatter.FormatKeyValue("ID", strconv.Itoa(u.ID)))
	fmt.Println(formatter.FormatKeyValue("Login", u.Login))
	fmt.Println(formatter.FormatKeyValue("Firstname", u.Firstname))
	fmt.Println(formatter.FormatKeyValue("Lastname", u.Lastname))
	fmt.Println(formatter.FormatKeyValue("Mail", u.Mail))
	fmt.Println(formatter.FormatKeyValue("Admin", strconv.FormatBool(u.Admin)))

	if u.CreatedOn != "" {
		fmt.Println(formatter.FormatKeyValue("Created", u.CreatedOn))
	}
	if u.LastLoginOn != "" {
		fmt.Println(formatter.FormatKeyValue("Last Login", u.LastLoginOn))
	}

	return nil
}

func init() {
	rootCmd.AddCommand(myAccountCmd)

	// Subcommands
	myAccountCmd.AddCommand(myAccountGetCmd)
	myAccountCmd.AddCommand(myAccountUpdateCmd)

	// Flags for get command
	myAccountGetCmd.Flags().StringP("format", "f", formatText, "出力フォーマット (json, text)")

	// Flags for update command
	myAccountUpdateCmd.Flags().String("firstname", "", "名")
	myAccountUpdateCmd.Flags().String("lastname", "", "姓")
	myAccountUpdateCmd.Flags().String("mail", "", "メールアドレス")
}
