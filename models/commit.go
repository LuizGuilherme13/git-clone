package models

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/LuizGuilherme13/git-clone/common"
)

type Commit struct {
	Hash    string `json:"hash"`
	Message string `json:"message"`
	Index   index  `json:"index"`
	Parent  string `json:"parent"`
}

func NewCommit(msg string) *Commit {
	return &Commit{Message: msg}
}

func (c *Commit) Commit() error {
	index, err := OpenIndexFile()
	if err != nil {
		return err
	}

	if err := index.Read(); err != nil {
		return err
	}
	c.Index = *index

	hash, err := common.CheckSum(common.IndexPath)
	if err != nil {
		return CheckError(err.Error())
	}
	c.Hash = hash

	head, err := OpenHeadFile()
	if err != nil {
		return err
	}
	c.Parent = head.Hash

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return CheckError(err.Error())
	}

	err = os.WriteFile(filepath.Join(common.ObjPath, c.Hash), data, 0644)
	if err != nil {
		return CheckError(err.Error())
	}

	if err := head.Write(c.Hash); err != nil {
		return err
	}

	return nil
}
