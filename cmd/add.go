package cmd

import (
	"bytes"
	"compress/gzip"
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

		compactContent, err := compressContent(content)
		if err != nil {
			fmt.Printf("Erro ao comprimir %s: %v", file, err)
			return
		}

		err = os.WriteFile(destPath, compactContent, 0644)
		if err != nil {
			fmt.Printf("Erro ao salvar %s: %v\n", file, err)
			return
		}
	}

}

func compressContent(content []byte) ([]byte, error) {
	buf := bytes.Buffer{}

	writer := gzip.NewWriter(&buf)

	_, err := writer.Write(content)
	if err != nil {
		return nil, err
	}

	if err := writer.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
