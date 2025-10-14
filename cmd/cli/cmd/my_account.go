package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"

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
		result, err := client.GetMyAccount(context.Background())
		if err != nil {
			return fmt.Errorf("アカウント情報の取得に失敗しました: %w", err)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("JSONのシリアライズに失敗しました: %w", err)
		}

		fmt.Println(string(output))
		return nil
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

func init() {
	rootCmd.AddCommand(myAccountCmd)

	// Subcommands
	myAccountCmd.AddCommand(myAccountGetCmd)
	myAccountCmd.AddCommand(myAccountUpdateCmd)

	// Flags for update command
	myAccountUpdateCmd.Flags().String("firstname", "", "名")
	myAccountUpdateCmd.Flags().String("lastname", "", "姓")
	myAccountUpdateCmd.Flags().String("mail", "", "メールアドレス")
}
