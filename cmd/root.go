package cmd

import (
	"os"

	"github.com/LuizGuilherme13/git-clone/common"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "backup",
	Short: "backup is a minimalistc git cli tool",
	Long:  "backup is a minimalistc git cli tool",
}

func Execute() {
	common.GetPaths()

	rootCmd.AddCommand(initRepoCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(logCmd)

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
