package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/LuizGuilherme13/git-clone/common"
	"github.com/LuizGuilherme13/git-clone/models"
	"github.com/spf13/cobra"
)

var initRepoCmd = &cobra.Command{
	Use: "init",
	Run: initRepo,
}

func initRepo(cmd *cobra.Command, args []string) {
	if err := os.MkdirAll(".backup/objects", os.ModePerm); err != nil {
		fmt.Println(models.CheckError(err.Error()))
		return
	}

	_, err := models.OpenIndexFile()
	if err != nil {
		fmt.Println(err)

		os.RemoveAll(filepath.Join(common.RootPath, ".backup"))

		return
	}

	_, err = models.OpenHeadFile()
	if err != nil {
		fmt.Println(err)

		os.RemoveAll(filepath.Join(common.RootPath, ".backup"))

		return
	}

	models.PrintOk(fmt.Sprintf("Repositório iniciado em: %s", common.RootPath))
}
