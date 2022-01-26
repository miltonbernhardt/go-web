package users

import (
	"encoding/json"
	"errors"
	"github.com/miltonbernhardt/go-web/internal/domain"
	"github.com/miltonbernhardt/go-web/internal/utils"
	"github.com/miltonbernhardt/go-web/pkg/store"
	"testing"

	"github.com/stretchr/testify/assert"
)

type DummyRepo struct{}

func (dr *DummyRepo) Delete(id int) error {
	return nil
}

func (dr *DummyRepo) GetAll() ([]domain.User, error) {
	return []domain.User{}, nil
}

func (dr *DummyRepo) Store(user domain.User) (domain.User, error) {
	return domain.User{}, nil
}

func (dr *DummyRepo) Update(id int, user domain.User) (domain.User, error) {
	return domain.User{}, nil
}

func (dr *DummyRepo) UpdateName(id int, name string) (domain.User, error) {
	return domain.User{}, nil
}

func (dr *DummyRepo) UpdateUser(id int, lastname string, age int) (domain.User, error) {
	return domain.User{}, nil
}

func TestServiceGetAll(t *testing.T) {
	input := []domain.User{
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
	dataJson, _ := json.Marshal(input)
	dbMock := store.Mock{
		Data: dataJson,
	}
	storeStub := store.FileStore{
		FileName: "",
		Mock:     &dbMock,
	}
	myRepo := NewRepository(&storeStub)
	myService := NewService(myRepo, utils.New())

	result, err := myService.GetAll()

	assert.Equal(t, input, result)
	assert.Nil(t, err)
}

func TestStore(t *testing.T) {
	testUser := domain.User{
		Firstname: "firstname",
		Lastname:  "lastname",
		Email:     "email",
		Age:       24,
		Height:    184,
		Active:    true,
	}
	dbMock := store.Mock{}

	storeStub := store.FileStore{
		FileName: "",
		Mock:     &dbMock,
	}
	myRepo := NewRepository(&storeStub)
	util := utils.New()
	util.AddMock(&utils.Mock{Date: "02/01/2006 15:04:05"})
	myService := NewService(myRepo, util)

	result, _ := myService.Store(testUser)

	assert.Equal(t, testUser.Firstname, result.Firstname)
	assert.Equal(t, testUser.Lastname, result.Lastname)
	assert.Equal(t, testUser.Email, result.Email)
	assert.Equal(t, testUser.Age, result.Age)
	assert.Equal(t, testUser.Height, result.Height)
	assert.Equal(t, testUser.Active, result.Active)
	assert.Equal(t, "02/01/2006 15:04:05", result.CreatedDate)
	assert.Equal(t, 4, result.ID)
}

func TestStoreError(t *testing.T) {
	testProduct := domain.User{
		Firstname: "firstname",
		Lastname:  "lastname",
		Email:     "email",
		Age:       24,
		Height:    184,
		Active:    true,
	}
	expectedError := errors.New(web.)
	dbMock := store.Mock{
		Err: expectedError,
	}
	storeStub := store.FileStore{
		FileName: "",
		Mock:     &dbMock,
	}
	myRepo := NewRepository(&storeStub)
	myService := NewService(myRepo, utils.New())
	result, err := myService.Store(testProduct)

	assert.Equal(t, expectedError, err)
	assert.Equal(t, Product{}, result)
}

//func TestSum(t *testing.T) {
//	expectedResult := float64(6)
//	myDummyRepo := DummyRepo{}
//	myService := NewService(&myDummyRepo)
//
//	result := myService.Sum(1, 2, 3)
//
//	assert.Equal(t, expectedResult, result)
//}
