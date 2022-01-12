package store

import (
	"encoding/json"
	"errors"
	"github.com/miltonbernhardt/go-web/internal/domain"
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
		return &FileStore{fileName, nil}
	}
	return nil
}

type FileStore struct {
	FileName FileName
	Mock     *Mock
}

type Mock struct {
	Data        []byte
	Err         error
	ReadWasUsed bool
}

func (fs *FileStore) AddMock(mock *Mock) {
	fs.Mock = mock
}

func (fs *FileStore) ClearMock() {
	fs.Mock = nil
}

func (fs *FileStore) Write(data interface{}) error {
	if fs.Mock != nil {
		return fs.Mock.Err
	}

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
	if fs.Mock != nil {
		return mockRead(fs.Mock, data)
	}

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

func mockRead(mock *Mock, data interface{}) error {
	mock.ReadWasUsed = true
	if mock.Err != nil {
		return mock.Err
	}

	if mock.Data == nil {
		file, _ := json.Marshal([]domain.User{
			{
				ID:          1,
				Firstname:   "firstname",
				Lastname:    "lastname",
				Email:       "email",
				Age:         24,
				Height:      184,
				Active:      true,
				CreatedDate: "22/02/2021",
			},
			{
				ID:          2,
				Firstname:   "firstname2",
				Lastname:    "lastname2",
				Email:       "email2",
				Age:         24,
				Height:      184,
				Active:      false,
				CreatedDate: "23/02/2021",
			},
			{
				ID:          3,
				Firstname:   "firstname3",
				Lastname:    "lastname3",
				Email:       "email3",
				Age:         26,
				Height:      187,
				Active:      false,
				CreatedDate: "25/02/2021",
			},
		})

		return json.Unmarshal(file, data)
	}
	return json.Unmarshal(mock.Data, data)
}
