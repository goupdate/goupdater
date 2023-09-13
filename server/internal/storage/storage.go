package storage

import (
	"os"
	"path/filepath"
	"time"

	"goupdater/server/internal/config"
)

type storage struct {
	path string //root of the storage directory

	idx *index
}

func Create() error {
	Storage = &storage{
		path: config.Config.Storage,
		idx:  &index{},
	}

	return Storage.idx.read(filepath.Join(Storage.path, "index.json"))
}

// search in existing database of uploaded files
func (s *storage) Find(projectName string, projectBranch string) (IndexItem, error) {
	return s.idx.find(projectName, projectBranch)
}

// already came gzipped
func (s *storage) Store(project string, branch string, crc int64, file []byte, szOrig int64) error {
	project = filepath.Clean(filepath.Base(project))
	branch = filepath.Clean(filepath.Base(branch))

	fpath := filepath.Join(s.path, project, branch, "last")

	os.MkdirAll(filepath.Dir(fpath), 0644)

	f, err := os.Create(fpath)
	if err != nil {
		return err
	}

	_, err = f.Write(file)
	if err != nil {
		return err
	}

	item := &IndexItem{
		Path:     fpath,
		Crc32:    crc,
		Size:     szOrig,
		Uploaded: time.Now().String(),
	}

	s.idx.Add(project, branch, item)
	return s.idx.save()
}

var Storage *storage
