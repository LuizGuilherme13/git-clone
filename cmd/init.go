package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"
)

var AbsDir string

func init() {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Erro ao obter destino do arquivo")
		return
	}

	AbsDir = wd

	rootCmd.AddCommand(initRepoCmd)
}

var initRepoCmd = &cobra.Command{
	Use: "init",
	Run: initRepo,
}

func initRepo(cmd *cobra.Command, args []string) {
	err := os.MkdirAll(".backup/objects", os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

}
