package cmd

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

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
	commit := Commit{}
	index := Index{Path: filepath.Join(AbsDir, ".backup", "index.json")}

	err := index.Unmarshal(index.Path)
	if err != nil {
		log.Fatalln(err)
	}
	commit.Index = index.Objects

	content, err := os.ReadFile(index.Path)
	if err != nil {
		log.Fatalln(err)
	}
	commit.Hash = fmt.Sprintf("%x", sha1.Sum(content))

	file, err := os.OpenFile(filepath.Join(AbsDir, ".backup", "HEAD.txt"), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	parentHash, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	commit.Parent = string(parentHash)

	data, err := json.MarshalIndent(commit, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(filepath.Join(AbsDir, ".backup/objects", commit.Hash), data, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	file.Seek(0, io.SeekStart)
	file.Truncate(0)
	_, err = file.WriteString(commit.Hash)
	if err != nil {
		log.Fatalln(err)
	}
}

type Commit struct {
	Hash   string   `json:"hash"`
	Index  []Object `json:"index"`
	Parent string   `json:"parent"`
}
