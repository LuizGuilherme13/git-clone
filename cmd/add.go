package cmd

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
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

	path := filepath.Join(AbsDir, ".backup/objects", "index.json")

	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		_, err = file.WriteString("{}")
		if err != nil {
			fmt.Println("Erro ao escrever no arquivo:", err)
			return
		}
	}

	for _, file := range args {

		// Lê o arquivo passado por parâmetro
		content, err := os.ReadFile(file)
		if err != nil {
			fmt.Printf("Erro ao ler %s: %v\n", file, err)
			return
		}

		fileName := uuid.New().String() + "_" + filepath.Base(file) + ".gz"
		destPath := filepath.Join(AbsDir, ".backup/objects", fileName)

		// Compactando
		compactContent, err := compressContent(content)
		if err != nil {
			fmt.Printf("Erro ao comprimir %s: %v", file, err)
			return
		}

		// Copiando o arquivo para o destino
		err = os.WriteFile(destPath, compactContent, 0644)
		if err != nil {
			fmt.Printf("Erro ao salvar %s: %v\n", file, err)
			return
		}

		index := Index{}

		// Abrir arquivo
		jsonFile, err := os.Open(path)
		if err != nil {
			fmt.Printf("Erro ao abrir .json: %v\n", err)
			return
		}
		defer jsonFile.Close()

		// Ler arquivo
		bytes, err := io.ReadAll(jsonFile)
		if err != nil {
			fmt.Printf("Erro ao ler .json: %v\n", err)
			return
		}

		// Deserializar
		err = json.Unmarshal(bytes, &index)
		if err != nil {
			fmt.Println("Erro ao desserializar JSON:", err)
			return
		}

		match := false
		for i := range index.Objects {
			obj := &index.Objects[i]

			if obj.Path == file {
				obj.Id = fileName
				match = true
				break
			}
		}

		if !match {
			index.Objects = append(index.Objects, Object{
				Path: file,
				Id:   fileName,
			})
		}

		updatedBytes, err := json.MarshalIndent(index, "", "  ")
		if err != nil {
			fmt.Println("Erro ao serializar JSON:", err)
			return
		}

		err = os.WriteFile(path, updatedBytes, 0644)
		if err != nil {
			fmt.Println("Erro ao gravar o arquivo:", err)
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

type Object struct {
	Path string `json:"path"`
	Id   string `json:"id"`
}

type Index struct {
	Objects []Object `json:"objects"`
}
