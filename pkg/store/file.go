package store

import (
	"encoding/json"
	"errors"
	"os"
)

type Store interface {
	Read(data interface{}) error
	Write(data interface{}) error
}

type Type string

type FileName string

const (
	FileType            Type     = "file"
	FileNameUsers       FileName = "./users-db.json"
	fileNameUsersBackup FileName = "./users-backup.json"
)

func New(store Type, fileName FileName) Store {
	switch store {
	case FileType:
		return &FileStore{fileName}
	}
	return nil
}

type FileStore struct {
	FileName FileName
}

func (fs *FileStore) Write(data interface{}) error {
	fileData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	f, err := os.OpenFile(string(fs.FileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	_, err = f.Write(fileData)
	if err != nil {
		return err
	}
	return nil
}

func (fs *FileStore) Read(data interface{}) error {
	file, err := os.ReadFile(string(fs.FileName))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return err
		}
		file, err = os.ReadFile(string(fileNameUsersBackup))

		if err != nil {
			return err
		}
	}
	return json.Unmarshal(file, &data)
}
