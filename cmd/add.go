package cmd

import (
	"fmt"

	"github.com/LuizGuilherme13/git-clone/models"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use: "add",
	Run: add,
}

func add(cmd *cobra.Command, args []string) {
	index, err := models.OpenIndexFile()
	if err != nil {
		fmt.Println(err)
		return
	}

	if err := index.Read(); err != nil {
		fmt.Println(err)
		return
	}

Next:
	for _, path := range args {
		blob, err := models.CreateBlob(path)
		if err != nil {
			fmt.Println(err)
			return
		}

		found := false
		for i := range index.Objects {
			obj := &index.Objects[i]

			if obj.Path == blob.Path {
				found = true

				if obj.Hash != blob.Hash {
					obj.Hash = blob.Hash

					if err := blob.Stage(); err != nil {
						fmt.Println(err)
						return
					}
				}

				continue Next
			}
		}

		if !found {
			if err := blob.Stage(); err != nil {
				fmt.Println(err)
				return
			}

			index.Objects = append(index.Objects, *blob)
		}
	}

	if err := index.Write(); err != nil {
		fmt.Println(err)
		return
	}
}
