package store

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/livecodeforlife/go-simple-aws/pkg/gocloud/core"
)

// Loader is where you can save or write the resource store
type Loader interface {
	Write(map[core.ID][]byte) error
	Read() (map[core.ID][]byte, error)
}

// New creates a new resource store
func New(loader Loader) core.ResourceStorer {
	return &resourceStore{
		loader: loader,
		data:   make(map[core.ID][]byte),
		locker: &sync.Mutex{},
	}
}

type resourceStore struct {
	loader Loader
	data   map[core.ID][]byte
	locker sync.Locker
}

func (s *resourceStore) Exists(ID core.ID) (bool, error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	_, ok := s.data[ID]
	return ok, nil
}

func (s *resourceStore) Get(ID core.ID) ([]byte, error) {
	s.locker.Lock()
	defer s.locker.Unlock()
	return s.data[ID], nil
}

func (s *resourceStore) Set(ID core.ID, data []byte) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	s.data[ID] = data
	return nil
}

func (s *resourceStore) Delete(ID core.ID) error {
	s.locker.Lock()
	defer s.locker.Unlock()
	delete(s.data, ID)
	return nil
}

func (s *resourceStore) Save() error {
	return s.loader.Write(s.data)
}

func (s *resourceStore) Load() error {
	data, err := s.loader.Read()
	if err != nil {
		return err
	}
	s.data = data
	return nil
}

type fileLoader struct {
	path string
}

func (f *fileLoader) Write(data map[core.ID][]byte) error {
	file, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		return err
	}
	return os.WriteFile(f.path, file, 0644)
}

func (f *fileLoader) Read() (map[core.ID][]byte, error) {
	var data = make(map[core.ID][]byte)
	if _, err := os.Stat(f.path); err == nil {
		log.Printf("Reading resource store file: %s", f.path)
		file, err := os.ReadFile(f.path)
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(file, &data); err != nil {
			return nil, err
		}
	} else if os.IsNotExist(err) {
		log.Printf("Creating new resource store file: %s", f.path)
	} else {
		fmt.Println("Error:", err)
	}
	return data, nil
}

// NewFileLoader creates a new file loader
func NewFileLoader(path string) Loader {
	return &fileLoader{
		path: path,
	}
}
