package models

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/LuizGuilherme13/git-clone/common"
)

// Index é responsável por armazenar as versões mais atualizadas de cada arquivo.
type index struct {
	Path    string `json:"-"`
	Objects []Blob `json:"objects"`
}

// OpenIndexFile verifica se o índice já existe, se não, o cria.
func OpenIndexFile() (*index, error) {
	if _, err := os.Stat(common.IndexPath); errors.Is(err, os.ErrNotExist) {
		file, err := os.Create(common.IndexPath)
		if err != nil {
			return nil, CheckError(err.Error())
		}
		defer file.Close()

		_, err = file.WriteString("{}")
		if err != nil {
			return nil, CheckError(err.Error())
		}
	}

	return &index{}, nil
}

// Read lê o índice e popula o 'Index.Objects' com o conteúdo atual.
func (i *index) Read() error {
	content, err := os.ReadFile(common.IndexPath)
	if err != nil {
		return CheckError(err.Error())
	}

	if err := json.Unmarshal(content, i); err != nil {
		return CheckError(err.Error())
	}

	return nil
}

// Write atualiza o índice com o conteúdo atual.
func (i *index) Write() error {
	data, err := json.MarshalIndent(*i, "", "  ")
	if err != nil {
		return CheckError(err.Error())
	}

	err = os.WriteFile(common.IndexPath, data, 0644)
	if err != nil {
		return CheckError(err.Error())
	}

	return nil
}
