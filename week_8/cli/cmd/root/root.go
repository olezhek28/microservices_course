package root

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "my-app",
	Short: "Моё cli приложение",
}

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Что-то создает",
}

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Что-то удаляет",
}

var createUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Создает нового пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		usernamesStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get usernames: %s\n", err.Error())
		}

		log.Printf("user %s created\n", usernamesStr)
	},
}

var deleteUserCmd = &cobra.Command{
	Use:   "user",
	Short: "Удаляет пользователя",
	Run: func(cmd *cobra.Command, args []string) {
		usernamesStr, err := cmd.Flags().GetString("username")
		if err != nil {
			log.Fatalf("failed to get usernames: %s\n", err.Error())
		}

		log.Printf("user %s deleted\n", usernamesStr)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)

	createCmd.AddCommand(createUserCmd)
	deleteCmd.AddCommand(deleteUserCmd)

	createUserCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	err := createUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}

	deleteUserCmd.Flags().StringP("username", "u", "", "Имя пользователя")
	err = deleteUserCmd.MarkFlagRequired("username")
	if err != nil {
		log.Fatalf("failed to mark username flag as required: %s\n", err.Error())
	}
}
