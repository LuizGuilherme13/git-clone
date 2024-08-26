package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use: "log",
	Run: logF,
}

func logF(cmd *cobra.Command, args []string) {
	fmt.Println("LOG")
}
