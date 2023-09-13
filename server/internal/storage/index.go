package storage

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"sync"

	"goupdater/server/internal/log"

	"github.com/MasterDimmy/jsonc"
)

/*
	index of stored items
*/

type IndexItem struct {
	Path string `json:",omitempty"` //path to the stored file

	Uploaded string

	Crc32 int64 //of stored body, calculated on client
	Size  int64 //in bytes
}

type indexBranch map[string]IndexItem

type index struct {
	file string

	Projects map[string]indexBranch
	m        sync.RWMutex
}

func (i *index) read(file string) error {
	i.m.Lock()
	defer i.m.Unlock()
	i.file = file

	buf, err := ioutil.ReadFile(file)
	if err != nil {
		log.Error_log.Printf("cant read storage index: %s\ncreating new.", err.Error())
	} else {
		return jsonc.Unmarshal(buf, i)
	}
	return nil
}

func (i *index) save() error {
	i.m.Lock()
	defer i.m.Unlock()

	// save to tmp first
	buf, err := json.MarshalIndent(i, "", " ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(i.file+".tmp", buf, 0644)
	if err != nil {
		return err
	}

	// rename if ok
	return os.Rename(i.file+".tmp", i.file)
}

// find file in database
func (i *index) find(project string, branch string) (IndexItem, error) {
	i.m.RLock()
	defer i.m.RUnlock()

	b, ok := i.Projects[project]
	if !ok {
		return IndexItem{}, nil
	}
	f, ok := b[branch]
	if !ok {
		return IndexItem{}, nil
	}
	return f, nil
}

func (i *index) Add(project string, branch string, item *IndexItem) {
	i.m.Lock()
	defer i.m.Unlock()

	if i.Projects == nil {
		i.Projects = make(map[string]indexBranch)
	}
	b := i.Projects[project]

	if b == nil {
		b = make(map[string]IndexItem)
	}
	b[branch] = *item

	i.Projects[project] = b
}
