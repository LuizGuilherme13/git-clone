package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LuizGuilherme13/git-clone/common"
	"github.com/LuizGuilherme13/git-clone/models"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use: "log",
	Run: logFunc,
}

func logFunc(cmd *cobra.Command, args []string) {
	head, err := models.OpenHeadFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	content, err := os.ReadFile(filepath.Join(common.ObjPath, head.Hash))
	if err != nil {
		fmt.Println(models.CheckError(err.Error()))
		return
	}

	commit := models.Commit{}
	if err := json.Unmarshal(content, &commit); err != nil {
		fmt.Println(models.CheckError(err.Error()))
		return
	}

	displayCommit(commit)
}

func displayCommit(commit models.Commit) {
	fmt.Printf("%scommit %s\n", common.ColorYellow, commit.Hash)
	fmt.Println(common.ColorReset)
	fmt.Printf("   %s\n", commit.Message)

	if commit.Parent == "" {
		fmt.Println("(END)")
		return
	}
	fmt.Println()

	getParent(commit)
}

func getParent(commit models.Commit) {
	content, err := os.ReadFile(filepath.Join(common.ObjPath, commit.Parent))
	if err != nil {
		println(models.CheckError(err.Error()))
		return
	}

	parent := models.Commit{}
	if err := json.Unmarshal(content, &parent); err != nil {
		println(models.CheckError(err.Error()))
		return
	}

	displayCommit(parent)
}
