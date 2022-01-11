package users

import (
	"encoding/json"
	"github.com/miltonbernhardt/go-web/internal/domain"
	"github.com/miltonbernhardt/go-web/pkg/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetAllUsers(t *testing.T) {
	expectedUsers := []domain.User{
		{
			Id:          1,
			Firstname:   "firstname",
			Lastname:    "lastname",
			Email:       "email",
			Age:         24,
			Height:      184,
			Active:      true,
			CreatedDate: "22/02/2021",
			DeletedDate: "",
		},
		{
			Id:          2,
			Firstname:   "firstname2",
			Lastname:    "lastname2",
			Email:       "email2",
			Age:         24,
			Height:      184,
			Active:      false,
			CreatedDate: "23/02/2021",
			DeletedDate: "",
		},
		{
			Id:          3,
			Firstname:   "firstname3",
			Lastname:    "lastname3",
			Email:       "email3",
			Age:         26,
			Height:      187,
			Active:      false,
			CreatedDate: "25/02/2021",
			DeletedDate: "26/02/2021",
		},
	}

	repository := NewRepository(&stubStore{})

	users, err := repository.FetchAllUsers()

	assert.Nil(t, err)
	assert.Equal(t, expectedUsers, users)
}

type stubStore struct {
	FileName store.FileName
}

func (stub *stubStore) Read(data interface{}) error {
	file, err := json.Marshal([]domain.User{
		{
			Id:          1,
			Firstname:   "firstname",
			Lastname:    "lastname",
			Email:       "email",
			Age:         24,
			Height:      184,
			Active:      true,
			CreatedDate: "22/02/2021",
			DeletedDate: "",
		},
		{
			Id:          2,
			Firstname:   "firstname2",
			Lastname:    "lastname2",
			Email:       "email2",
			Age:         24,
			Height:      184,
			Active:      false,
			CreatedDate: "23/02/2021",
			DeletedDate: "",
		},
		{
			Id:          3,
			Firstname:   "firstname3",
			Lastname:    "lastname3",
			Email:       "email3",
			Age:         26,
			Height:      187,
			Active:      false,
			CreatedDate: "25/02/2021",
			DeletedDate: "26/02/2021",
		},
	})

	if err != nil {
		return err
	}

	return json.Unmarshal(file, &data)
}

func (stub *stubStore) Write(data interface{}) error {
	return nil
}
