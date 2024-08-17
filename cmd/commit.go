package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(commitCmd)
}

var commitCmd = &cobra.Command{
	Use: "commit",
	Run: commit,
}

func commit(cmd *cobra.Command, args []string) {
	fmt.Println("COMMIT")
}
