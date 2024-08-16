package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use: "status",
	Run: status,
}

func status(cmd *cobra.Command, args []string) {
	fmt.Println("STATUS")
}
