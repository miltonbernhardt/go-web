package users

import (
	"context"
	"github.com/miltonbernhardt/go-web/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

func (r *repositoryMock) GetAll() ([]model.User, error) {
	args := r.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (r *repositoryMock) GetByFirstname(firstname string) (model.User, error) {
	args := r.Called(firstname)
	return args.Get(0).(model.User), args.Error(1)
}

func (r *repositoryMock) Store(user model.User) (model.User, error) {
	args := r.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (r *repositoryMock) Update(user model.User) (model.User, error) {
	args := r.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (r *repositoryMock) UpdateUserFirstname(id int, name string) error {
	args := r.Called(id, name)
	return args.Error(0)
}

func (r *repositoryMock) UpdateUserAgeLastname(id int, lastname string, age int) error {
	args := r.Called(id, lastname, age)
	return args.Error(0)
}

func (r *repositoryMock) UpdateWithContext(ctx context.Context, user model.User) (model.User, error) {
	args := r.Called(ctx, user)
	return args.Get(0).(model.User), args.Error(1)
}

func Test_DB_Service_GetAll_Success(t *testing.T) {
	repository := new(repositoryMock)

	repository.On("GetAll").Return(
		[]model.User{
			{ID: 1, Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"},
			{ID: 2, Firstname: "firstname2", Lastname: "lastname2", Email: "email2", Age: 25, Height: 185, Active: false, CreatedDate: "23/02/2021"},
			{ID: 3, Firstname: "firstname3", Lastname: "lastname3", Email: "email3", Age: 26, Height: 186, Active: false, CreatedDate: "24/02/2021"},
		},
		nil,
	)

	input := []model.User{
		{ID: 1, Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"},
		{ID: 2, Firstname: "firstname2", Lastname: "lastname2", Email: "email2", Age: 25, Height: 185, Active: false, CreatedDate: "23/02/2021"},
		{ID: 3, Firstname: "firstname3", Lastname: "lastname3", Email: "email3", Age: 26, Height: 186, Active: false, CreatedDate: "24/02/2021"},
	}

	service := NewService(repository, nil)
	expected, err := service.GetAll()

	assert.Equal(t, input, expected)
	assert.Nil(t, err)
}

//func TestStore(t *testing.T) {
//	testUser := model.User{
//		Firstname: "firstname",
//		Lastname:  "lastname",
//		Email:     "email",
//		Age:       24,
//		Height:    184,
//		Active:    true,
//	}
//	dbMock := store.Mock{}
//
//	storeStub := store.FileStore{
//		FileName: "",
//		Mock:     &dbMock,
//	}
//	myRepo := NewRepository(&storeStub)
//	util := utils.New()
//	util.AddMock(&utils.Mock{Date: "02/01/2006 15:04:05"})
//	myService := NewService(myRepo, util)
//
//	result, _ := myService.Store(testUser)
//
//	assert.Equal(t, testUser.Firstname, result.Firstname)
//	assert.Equal(t, testUser.Lastname, result.Lastname)
//	assert.Equal(t, testUser.Email, result.Email)
//	assert.Equal(t, testUser.Age, result.Age)
//	assert.Equal(t, testUser.Height, result.Height)
//	assert.Equal(t, testUser.Active, result.Active)
//	assert.Equal(t, "02/01/2006 15:04:05", result.CreatedDate)
//	assert.Equal(t, 4, result.ID)
//}
//
//func TestStoreError(t *testing.T) {
//	testProduct := model.User{
//		Firstname: "firstname",
//		Lastname:  "lastname",
//		Email:     "email",
//		Age:       24,
//		Height:    184,
//		Active:    true,
//	}
//	expectedError := errors.New(web.)
//	dbMock := store.Mock{
//		Err: expectedError,
//	}
//	storeStub := store.FileStore{
//		FileName: "",
//		Mock:     &dbMock,
//	}
//	myRepo := NewRepository(&storeStub)
//	myService := NewService(myRepo, utils.New())
//	result, err := myService.Store(testProduct)
//
//	assert.Equal(t, expectedError, err)
//	assert.Equal(t, Product{}, result)
//}
//
////func TestSum(t *testing.T) {
////	expectedResult := float64(6)
////	myDummyRepo := RepositoryMock{}
////	myService := NewService(&myDummyRepo)
////
////	result := myService.Sum(1, 2, 3)
////
////	assert.Equal(t, expectedResult, result)
////}
