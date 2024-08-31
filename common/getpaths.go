package common

import (
	"os"
	"path/filepath"
)

func GetPaths() {
	rootPath, err := os.Getwd()
	if err != nil {
		println("Erro ao obter raíz do diretório")
		os.Exit(1)
	}

	RootPath = rootPath
	ObjPath = filepath.Join(rootPath, ".backup", "objects")
	IndexPath = filepath.Join(rootPath, ".backup", "index.json")
	HeadPath = filepath.Join(rootPath, ".backup", "HEAD.txt")
}
