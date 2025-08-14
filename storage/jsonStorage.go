package storage

import (
	"os"
)

type JsonStorage struct {
	*Storage
}

func NewJsonStorage(filename string) *JsonStorage {
	return &JsonStorage{
		&Storage{filename: filename},
	}
}

func (s *JsonStorage) Save(data []byte) error {
	err := os.WriteFile(s.GetFilename(), data, 0644)
	return err
}

func (s *JsonStorage) Load() ([]byte, error) {
	d, err := os.ReadFile(s.GetFilename())
	return d, err
}
