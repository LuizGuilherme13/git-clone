package models

import (
	"bytes"
	"compress/gzip"
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"

	"github.com/LuizGuilherme13/git-clone/common"
)

type Blob struct {
	Hash string `json:"hash"`
	Path string `json:"path"`
	Data []byte `json:"-"`
}

func CreateBlob(path string) (*Blob, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, CheckError(err.Error())
	}

	buf := bytes.Buffer{}

	writer := gzip.NewWriter(&buf)

	if _, err := writer.Write(content); err != nil {
		return nil, CheckError(err.Error())
	}

	if err := writer.Close(); err != nil {
		return nil, CheckError(err.Error())
	}

	blob := &Blob{Path: path}
	blob.Hash = fmt.Sprintf("%x", sha1.Sum(content))
	blob.Data = buf.Bytes()

	return blob, nil
}

func (b *Blob) Stage() error {

	err := os.WriteFile(filepath.Join(common.ObjPath, b.Hash), b.Data, 0644)
	if err != nil {
		return CheckError(err.Error())
	}

	return nil
}
