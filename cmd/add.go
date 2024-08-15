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
	pathToIndex := filepath.Join(AbsDir, ".backup/objects", "index.json")

	if _, err := os.Stat(pathToIndex); errors.Is(err, os.ErrNotExist) {
		if err := createIndexFile(pathToIndex); err != nil {
			log.Fatalln(err)
		}
	}

	index := Index{}

	if err := index.Unmarshal(pathToIndex); err != nil {
		log.Fatalln(err)
	}

	for _, file := range args {
		object := Object{
			Path: file,
			Id:   uuid.New().String() + "_" + filepath.Base(file) + ".gz",
		}

		if err := object.compress(); err != nil {
			log.Fatalln(err)
		}

		if err := object.staging(); err != nil {
			log.Fatalln(err)
		}

		found := false
		for i := range index.Objects {
			obj := &index.Objects[i]

			if obj.Path == object.Path {
				found = true
				obj.Id = object.Id
				break
			}
		}

		if !found {
			index.Objects = append(index.Objects, object)
		}
	}

	data, err := json.MarshalIndent(index, "", "  ")
	if err != nil {
		fmt.Println("Erro ao serializar JSON:", err)
		return
	}

	err = os.WriteFile(pathToIndex, data, 0644)
	if err != nil {
		fmt.Println("Erro ao gravar o arquivo:", err)
		return
	}
}

func createIndexFile(pathToIndex string) error {
	file, err := os.Create(pathToIndex)
	if err != nil {
		return fmt.Errorf("erro ao criar index.json: %w", err)
	}
	defer file.Close()

	_, err = file.WriteString("{}")
	if err != nil {
		return fmt.Errorf("erro ao escrever no arquivo: %w", err)
	}

	return nil
}

type Object struct {
	Id   string `json:"id"`
	Path string `json:"path"`
	Data []byte `json:"-"`
}

func (obj *Object) compress() error {
	content, err := os.ReadFile(obj.Path)
	if err != nil {
		return fmt.Errorf("erro ao ler %s: %v", obj.Path, err)
	}

	buf := bytes.Buffer{}
	writer := gzip.NewWriter(&buf)

	_, err = writer.Write(content)
	if err != nil {
		return err
	}

	if err := writer.Close(); err != nil {
		return err
	}

	obj.Data = buf.Bytes()

	return nil
}

func (obj *Object) staging() error {
	dest := filepath.Join(AbsDir, ".backup/objects", obj.Id)

	err := os.WriteFile(dest, obj.Data, 0644)
	if err != nil {
		return fmt.Errorf("erro ao salvar %s: %v", obj.Path, err)

	}

	return nil
}

type Index struct {
	Objects []Object `json:"objects"`
}

func (i *Index) Unmarshal(pathToIndex string) error {
	file, err := os.Open(pathToIndex)
	if err != nil {
		return fmt.Errorf("erro ao abrir .json: %v", err)
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("erro ao ler .json: %v", err)
	}

	err = json.Unmarshal(data, i)
	if err != nil {
		return fmt.Errorf("erro ao desserializar .json: %w", err)
	}

	return nil
}
