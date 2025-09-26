package cmd

import (
	"fmt"
	"os"

	"github.com/ants-1/github-user-activity/service"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "github-user-activity",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		if name == "" || len(name) < 1 {
			fmt.Println("Please enter a name")
			return
		}

		activities, err := service.GetActivity(name)
		if err != nil {
			fmt.Println(err)
			return
		}

		if len(activities) == 0 {
			fmt.Println("No recent activity found.")
			return
		}

		fmt.Println("Recent GitHub activity:")
		for _, act := range activities {
			fmt.Println("-", service.FormatEvent(act))
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
