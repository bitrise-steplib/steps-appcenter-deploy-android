package util

import (
	"io/ioutil"
	"strings"
)

// LocalFile ...
type LocalFile struct {
	FilePath    string
	FileContent []byte
}

// FileName ...
func (lf LocalFile) FileName() string {
	pathParts := strings.Split(lf.FilePath, "/")

	return pathParts[len(pathParts)-1]
}

// OpenFile ...
func (lf *LocalFile) OpenFile() error {
	fb, err := ioutil.ReadFile(lf.FilePath)
	if err != nil {
		return err
	}

	lf.FileContent = fb
	return err
}

// FileSize ...
func (lf LocalFile) FileSize() int {
	return len(lf.FileContent)
}

// MakeChunks ...
func (lf LocalFile) MakeChunks(limit int) [][]byte {
	var retData [][]byte
	for i := 0; i < len(lf.FileContent); i += limit {
		batch := lf.FileContent[i:min(i+limit, len(lf.FileContent))]
		retData = append(retData, batch)
	}

	return retData
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
