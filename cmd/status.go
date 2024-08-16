package cmd

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"slices"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(statusCmd)
}

var statusCmd = &cobra.Command{
	Use: "status",
	Run: status,
}

var files = []string{}

func status(cmd *cobra.Command, args []string) {
	err := filepath.Walk(AbsDir, walk)
	if err != nil {
		fmt.Println(err)
	}

	index := Index{Path: filepath.Join(AbsDir, ".backup", "index.json")}
	if err := index.Unmarshal(index.Path); err != nil {
		log.Fatalln(err)
	}

	untracked := []string{}
	objects := []string{}

	for _, obj := range index.Objects {
		objects = append(objects, obj.Path)
	}

	for _, file := range files {
		if contain := slices.Contains(objects, file); !contain {
			untracked = append(untracked, file)
		}
	}

	fmt.Println("Untracked files:")
	for i := range untracked {
		fmt.Printf("\t%s\n", untracked[i])
	}

}

func walk(path string, info fs.FileInfo, err error) error {
	if err != nil {
		fmt.Println(err)
		return err
	}

	relPath, err := filepath.Rel(AbsDir, path)
	if err != nil {
		fmt.Println(err)
		return err
	}

	if !info.IsDir() && relPath != "." && !(len(relPath) > 1 && relPath[0] == '.') {
		files = append(files, relPath)
	}

	return nil
}
