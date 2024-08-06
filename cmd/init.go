package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(initRepoCmd)
}

var initRepoCmd = &cobra.Command{
	Use: "init",
	Run: initRepo,
}

func initRepo(cmd *cobra.Command, args []string) {
	err := os.MkdirAll(".got/objects", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
