package users

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/miltonbernhardt/go-web/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

type repositoryMock struct {
	mock.Mock
}

func (r *repositoryMock) GetAll() ([]model.User, error) {
	args := r.Called()
	return args.Get(0).([]model.User), args.Error(1)
}

func (r *repositoryMock) Delete(id int) error {
	args := r.Called(id)
	return args.Error(0)
}

func (r *repositoryMock) Update(id int, userToUpdate model.User) (model.User, error) {
	args := r.Called(id, userToUpdate)
	return args.Get(0).(model.User), args.Error(1)
}

func (r *repositoryMock) UpdateName(id int, name string) (model.User, error) {
	args := r.Called(id, name)
	return args.Get(0).(model.User), args.Error(1)
}

func (r *repositoryMock) UpdateUser(id int, lastname string, age int) (model.User, error) {
	args := r.Called(id, lastname, age)
	return args.Get(0).(model.User), args.Error(1)
}

func (r *repositoryMock) Store(user model.User) (model.User, error) {
	args := r.Called(user)
	return args.Get(0).(model.User), args.Error(1)
}

func (r *repositoryMock) getUserLastID() (int, error) {
	args := r.Called()
	return args.Get(0).(int), args.Error(1)
}

func Test_DB_GetAll_Success(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	expectedSections := []model.User{
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
			Age:         25,
			Height:      185,
			Active:      false,
			CreatedDate: "23/02/2021",
		},
		{
			ID:          3,
			Firstname:   "firstname3",
			Lastname:    "lastname3",
			Email:       "email3",
			Age:         26,
			Height:      186,
			Active:      false,
			CreatedDate: "24/02/2021",
		},
	}

	rows := mockDB.
		NewRows([]string{"id", "firstname", "lastname", "email", "age", "height", "active", "created_date"}).
		AddRow(1, "firstname1", "lastname2", "email", 24, 184, true, "22/02/2021").
		AddRow(2, "firstname2", "lastname2", "email", 25, 185, false, "23/02/2021").
		AddRow(3, "firstname3", "lastname2", "email", 26, 186, false, "24/02/2021")
	mockDB.ExpectQuery("SELECT").WillReturnRows(rows)

	repository := NewRepository(db)
	actualSections, err := repository.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, expectedSections, actualSections)
}
