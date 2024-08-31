package cmd

import (
	"fmt"

	"github.com/LuizGuilherme13/git-clone/models"
	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use: "commit",
	Run: commit,
}

func commit(cmd *cobra.Command, args []string) {
	commit := models.Commit{Message: args[0]}

	if err := commit.Commit(); err != nil {
		fmt.Println(err)
		return
	}
}
