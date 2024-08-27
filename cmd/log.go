package cmd

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(logCmd)
}

var logCmd = &cobra.Command{
	Use: "log",
	Run: logFunc,
}

var pathToCommits = filepath.Join(AbsDir, ".backup", "objects")

func logFunc(cmd *cobra.Command, args []string) {
	headCommit, err := os.ReadFile(filepath.Join(AbsDir, ".backup", "HEAD.txt"))
	if err != nil {
		log.Fatalln(err)
	}

	content, err := os.ReadFile(filepath.Join(pathToCommits, string(headCommit)))
	if err != nil {
		log.Fatalln(err)
	}

	commit := Commit{}
	if err := json.Unmarshal(content, &commit); err != nil {
		log.Fatalln(err)
	}

	getParent(pathToCommits, commit)
}

func getParent(path string, commit Commit) {

	fmt.Printf("commit %s\n", commit.Hash)
	fmt.Println()
	fmt.Printf("   %s\n", commit.Message)

	if commit.Parent == "" {
		fmt.Println("(END)")
		return
	}
	fmt.Println()

	content, err := os.ReadFile(filepath.Join(path, commit.Parent))
	if err != nil {
		log.Fatalln(err)
	}

	parentCommit := Commit{}
	if err := json.Unmarshal(content, &parentCommit); err != nil {
		log.Fatalln(err)
	}

	getParent(pathToCommits, parentCommit)
}
