package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"log"
	"os"
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

	objects := []string{}
	untracked := []string{}
	modified := []string{}
	toBeCommited := []string{}

	//* Buscando o último commit
	headCommit, err := os.ReadFile(filepath.Join(AbsDir, ".backup", "HEAD.txt"))
	if err != nil {
		log.Fatalln(err)
	}

	content, err := os.ReadFile(filepath.Join(pathToCommits, string(headCommit)))
	if err != nil {
		log.Fatalln(err)
	}

	lastCommit := Commit{}
	if err := json.Unmarshal(content, &lastCommit); err != nil {
		log.Fatalln(err)
	}
	//*

	for _, obj := range index.Objects {
		objects = append(objects, obj.Path)
	}

	for _, file := range files {
		if contain := slices.Contains(objects, file); !contain {
			untracked = append(untracked, file)
		} else {
			object := Object{Path: file}
			if err := object.compress(); err != nil {
				log.Fatalln(err)
			}

			isModified := slices.ContainsFunc(index.Objects, func(obj Object) bool {
				if obj.Path != object.Path {
					return false
				}

				return obj.Id != object.Id
			})

			if isModified {
				modified = append(modified, file)
			}
		}
	}

Main:
	for _, obj := range index.Objects {
		neverCommited := true

		for i := range lastCommit.Index {
			if obj.Path == lastCommit.Index[i].Path {
				neverCommited = false

				if obj.Id != lastCommit.Index[i].Id {
					toBeCommited = append(toBeCommited, obj.Path)
				}

				continue Main
			}
		}

		if neverCommited {
			toBeCommited = append(toBeCommited, obj.Path)
		}
	}

	if len(toBeCommited) > 0 {
		fmt.Println("Changes to be commited:")
		for i := range toBeCommited {
			fmt.Printf("\t%smodified:   %s\n", Green, toBeCommited[i])
		}
		fmt.Println(Reset)
	}

	if len(modified) > 0 {
		fmt.Println("Changes not staged:")
		for i := range modified {
			fmt.Printf("\t%smodified:   %s\n", Yellow, modified[i])
		}
		fmt.Println(Reset)
	}

	if len(untracked) > 0 {
		fmt.Println("Untracked files:")
		for i := range untracked {
			fmt.Printf("\t%s%s\n", Red, untracked[i])
		}
		fmt.Println(Reset)
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
