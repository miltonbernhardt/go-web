package users

import (
	"encoding/json"
	"github.com/miltonbernhardt/go-web/internal/model"
	"github.com/miltonbernhardt/go-web/pkg/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_File_GetAll_Success(t *testing.T) {
	expectedUsers := []model.User{
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
	}

	dataJson, _ := json.Marshal(expectedUsers)
	dbStub := store.Mock{
		Data: dataJson,
		Err:  nil,
	}
	storeMocked := store.FileStore{
		FileName: "",
		Mock:     &dbStub,
	}

	repository := NewRepositoryFile(&storeMocked)

	users, err := repository.GetAll()

	assert.Nil(t, err)
	assert.Equal(t, expectedUsers, users)
}

func Test_File_UpdateName_Success(t *testing.T) {
	mock := store.Mock{}

	storeMocked := store.FileStore{
		FileName: "",
		Mock:     &mock,
	}

	repository := NewRepositoryFile(&storeMocked)

	actualUser, err := repository.UpdateName(1, "After Update")
	expectedUser := model.User{
		ID:          1,
		Firstname:   "After Update",
		Lastname:    "lastname",
		Email:       "email",
		Age:         24,
		Height:      184,
		Active:      true,
		CreatedDate: "22/02/2021",
	}

	assert.Nil(t, err)
	assert.Equal(t, actualUser, expectedUser)
	assert.True(t, mock.ReadWasUsed)
}
