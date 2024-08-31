package cmd

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/LuizGuilherme13/git-clone/common"
	"github.com/LuizGuilherme13/git-clone/models"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use: "status",
	Run: status,
}

var files = []string{}

func status(cmd *cobra.Command, args []string) {
	err := filepath.Walk(common.RootPath, walk)
	if err != nil {
		fmt.Println(err)
		return
	}

	index, err := models.OpenIndexFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := index.Read(); err != nil {
		fmt.Println(err)
		return
	}

	objects := []string{}
	untracked := []string{}
	modified := []string{}
	toBeCommited := []string{}

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

	lastCommit := models.Commit{}
	if err := json.Unmarshal(content, &lastCommit); err != nil {
		fmt.Println(models.CheckError(err.Error()))
		return
	}

	for _, obj := range index.Objects {
		objects = append(objects, obj.Path)
	}

	for _, file := range files {
		if contain := slices.Contains(objects, file); !contain {
			untracked = append(untracked, file)
		} else {
			blob, err := models.CreateBlob(file)
			if err != nil {
				fmt.Println(err)
				return
			}

			isModified := slices.ContainsFunc(index.Objects, func(obj models.Blob) bool {
				if obj.Path != blob.Path {
					return false
				}

				return obj.Hash != blob.Hash
			})

			if isModified {
				modified = append(modified, file)
			}
		}
	}

Main:
	for _, indexedObj := range index.Objects {
		neverCommited := true

		for _, commitedObj := range lastCommit.Index.Objects {
			if indexedObj.Path == commitedObj.Path {
				neverCommited = false

				if indexedObj.Hash != commitedObj.Hash {
					toBeCommited = append(toBeCommited, indexedObj.Path)
				}

				continue Main
			}
		}

		if neverCommited {
			toBeCommited = append(toBeCommited, indexedObj.Path)
		}
	}

	if len(toBeCommited) > 0 {
		fmt.Println("Changes to be commited:")
		for i := range toBeCommited {
			fmt.Printf("\t%s%s\n", common.ColorGreen, toBeCommited[i])
		}
		fmt.Println(common.ColorReset)
	}

	if len(modified) > 0 {
		fmt.Println("Changes not staged:")
		for i := range modified {
			fmt.Printf("\t%s%s\n", common.ColorYellow, modified[i])
		}
		fmt.Println(common.ColorReset)
	}

	if len(untracked) > 0 {
		fmt.Println("Untracked files:")
		for i := range untracked {
			fmt.Printf("\t%s%s\n", common.ColorRed, untracked[i])
		}
		fmt.Println(common.ColorReset)
	}

}

func walk(path string, info fs.FileInfo, err error) error {
	if err != nil {
		return models.CheckError(err.Error())
	}

	relPath, err := filepath.Rel(common.RootPath, path)
	if err != nil {
		return models.CheckError(err.Error())
	}

	if !info.IsDir() && relPath != "." && !(len(relPath) > 1 && relPath[0] == '.') {
		files = append(files, relPath)
	}

	return nil
}
