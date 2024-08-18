package cmd

import (
	"crypto/sha1"
	"encoding/json"
	"errors"
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
	commit := Commit{Index: make(map[string]string)}
	index := Index{Path: filepath.Join(AbsDir, ".backup", "index.json")}

	err := index.Unmarshal(index.Path)
	if err != nil {
		log.Fatalln(err)
	}

	for _, obj := range index.Objects {
		commit.Index[obj.Path] = obj.Id
	}

	content, err := os.ReadFile(index.Path)
	if err != nil {
		log.Fatalln(err)
	}
	commit.Hash = fmt.Sprintf("%x", sha1.Sum(content))

	head := filepath.Join(AbsDir, ".backup", "HEAD.txt")

	file, err := os.OpenFile(head, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	parent, err := io.ReadAll(file)
	if err != nil {
		log.Fatalln(err)
	}
	commit.Parent = string(parent)

	dest := filepath.Join(AbsDir, ".backup/objects", commit.Hash)

	data, err := json.MarshalIndent(commit, "", "  ")
	if err != nil {
		log.Fatalln(err)
	}

	err = os.WriteFile(dest, data, 0644)
	if err != nil {
		log.Fatalln(err)
	}

	if _, err := os.Stat(head); errors.Is(err, os.ErrNotExist) {
		_, err := os.Create(head)
		if err != nil {
			log.Fatalln(err)
		}
	}

	file.Seek(0, io.SeekStart)
	file.Truncate(0)
	_, err = file.WriteString(commit.Hash)
	if err != nil {
		log.Fatalln(err)
	}
}

type Commit struct {
	Hash   string            `json:"hash"`
	Index  map[string]string `json:"index"`
	Parent string            `json:"parent"`
}
