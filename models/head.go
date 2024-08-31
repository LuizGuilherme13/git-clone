package models

import (
	"io"
	"os"

	"github.com/LuizGuilherme13/git-clone/common"
)

type head struct {
	Path string
	Hash string
}

func OpenHeadFile() (*head, error) {
	file, err := os.OpenFile(common.HeadPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil, CheckError(err.Error())
	}
	defer file.Close()

	hash, err := io.ReadAll(file)
	if err != nil {
		return nil, CheckError(err.Error())
	}

	return &head{Hash: string(hash)}, nil
}

func (h *head) Write(newCommitHash string) error {
	file, err := os.OpenFile(common.HeadPath, os.O_RDWR, 0644)
	if err != nil {
		return CheckError(err.Error())
	}
	defer file.Close()

	file.Seek(0, io.SeekStart)
	file.Truncate(0)

	_, err = file.WriteString(newCommitHash)
	if err != nil {
		return CheckError(err.Error())
	}

	return nil
}
