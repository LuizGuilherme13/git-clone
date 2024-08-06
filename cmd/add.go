package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(addCmd)
}

var addCmd = &cobra.Command{
	Use: "add",
	Run: add,
}

func add(cmd *cobra.Command, args []string) {
	for _, file := range args {
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Erro ao ler %s: %v\n", file, err)
			return
		}

		dest, err := os.Getwd()
		if err != nil {
			fmt.Println("Erro ao obter destino do arquivo")
			return
		}

		fileName := filepath.Base(file) + ".gz"
		destPath := filepath.Join(dest, ".got/objects", fileName)

		err = os.WriteFile(destPath, content, 0644)
		if err != nil {
			fmt.Printf("Erro ao salvar %s: %v\n", file, err)
			return
		}
	}

}
