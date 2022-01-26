package users

import (
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/miltonbernhardt/go-web/internal/model"
	"github.com/miltonbernhardt/go-web/pkg/message"
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_DB_Mock_GetAll_Success(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	expectedSections := []model.User{
		{ID: 1, Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"},
		{ID: 2, Firstname: "firstname2", Lastname: "lastname2", Email: "email2", Age: 25, Height: 185, Active: false, CreatedDate: "23/02/2021"},
		{ID: 3, Firstname: "firstname3", Lastname: "lastname3", Email: "email3", Age: 26, Height: 186, Active: false, CreatedDate: "24/02/2021"},
	}

	rows := mockDB.
		NewRows([]string{"id", "firstname", "lastname", "email", "age", "height", "active", "created_date"}).
		AddRow(1, "firstname1", "lastname1", "email1", 24, 184, true, "22/02/2021").
		AddRow(2, "firstname2", "lastname2", "email2", 25, 185, false, "23/02/2021").
		AddRow(3, "firstname3", "lastname3", "email3", 26, 186, false, "24/02/2021")
	mockDB.ExpectQuery(getAllQuery).WillReturnRows(rows)

	repository := NewRepository(db)
	actualSections, err := repository.GetAll()
	assert.Nil(t, err)
	assert.Equal(t, expectedSections, actualSections)
}

func Test_DB_Mock_GetAll_Fail(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	repository := NewRepository(db)
	actualSections, err := repository.GetAll()
	assert.Error(t, err)
	assert.Nil(t, actualSections)
}

func Test_DB_Mock_GetByFirstname_Success(t *testing.T) {
	query := `SELECT id, firstname, lastname, email, age, height, active, created_at FROM users WHERE firstname = \?`
	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	expectedSection := model.User{ID: 1, Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"}

	rows := mockDB.
		NewRows([]string{"id", "firstname", "lastname", "email", "age", "height", "active", "created_date"}).
		AddRow(1, "firstname1", "lastname1", "email1", 24, 184, true, "22/02/2021")
	mockDB.ExpectQuery(query).WillReturnRows(rows)

	repository := NewRepository(db)
	actualSection, err := repository.GetByFirstname("firstname1")
	assert.Nil(t, err)
	assert.Equal(t, expectedSection, actualSection)
}

func Test_DB_Mock_GetByFirstname_Fail(t *testing.T) {
	query := `SELECT id, firstname, lastname, email, age, height, active, created_at FROM users WHERE firstname = \?`
	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	rows := mockDB.
		NewRows([]string{"id", "firstname", "lastname", "email", "age", "height", "active", "created_date"})
	mockDB.ExpectQuery(query).WillReturnRows(rows)

	repository := NewRepository(db)
	actualSection, err := repository.GetByFirstname("firstname1")
	assert.Equal(t, errors.New(message.UserNotFound), err)
	assert.Equal(t, model.User{}, actualSection)
}

func Test_DB_Mock__GetByFirstname_Fail(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	repository := NewRepository(db)
	actualSections, err := repository.GetAll()
	assert.Error(t, err)
	assert.Nil(t, actualSections)
}

func Test_DB_Mock_Delete_Success(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	mockDB.ExpectPrepare(deleteQuery).
		WillBeClosed().
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewRepository(db)
	err := repository.Delete(1)

	assert.Nil(t, err)
}

func Test_DB_Mock_Delete_Fail_Prepare(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	repository := NewRepository(db)
	err := repository.Delete(1)

	assert.Error(t, err)
}

func Test_DB_Mock_Delete_Fail_Exec(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	mockDB.ExpectPrepare(deleteQuery)

	repository := NewRepository(db)
	err := repository.Delete(1)

	assert.Error(t, err)
}

func Test_DB_Mock_Update_Success(t *testing.T) {
	query := `UPDATE users SET firstname = \?, lastname = \?, email = \?, age = \?, height = \?, active = \?, created_at = \? WHERE id = \?`
	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	mockDB.ExpectPrepare(query).
		WillBeClosed().
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewRepository(db)
	actualUser, err := repository.Update(model.User{ID: 1, Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"})

	expectedUser := model.User{ID: 1, Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"}

	assert.Nil(t, err)
	assert.Equal(t, expectedUser, actualUser)
}

func Test_DB_Mock_Update_Fail_Prepare(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	repository := NewRepository(db)
	actualUser, err := repository.Update(model.User{ID: 1, Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"})
	assert.Error(t, err)
	assert.Equal(t, model.User{}, actualUser)
}

func Test_DB_Mock_Update_Fail_Exec(t *testing.T) {
	query := `UPDATE users SET firstname = \?, lastname = \?, email = \?, age = \?, height = \?, active = \?, created_at = \? WHERE id = \?`
	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	mockDB.ExpectPrepare(query)

	repository := NewRepository(db)
	actualUser, err := repository.Update(model.User{ID: 1, Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"})

	assert.Error(t, err)
	assert.Equal(t, model.User{}, actualUser)
}

func Test_DB_Mock_UpdateUserAgeLastname_Success(t *testing.T) {
	query := `UPDATE users SET lastname = \?, age = \? WHERE id = \?`

	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	mockDB.ExpectPrepare(query).
		WillBeClosed().
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewRepository(db)
	err := repository.UpdateUserAgeLastname(1, "lastname", 24)

	assert.Nil(t, err)
}

func Test_DB_Mock_UpdateUserAgeLastname_Fail_Prepare(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	repository := NewRepository(db)
	err := repository.UpdateUserAgeLastname(1, "lastname", 24)

	assert.Error(t, err)
}

func Test_DB_Mock_UpdateUserAgeLastname_Fail_Exec(t *testing.T) {
	query := `UPDATE users SET lastname = \?, age = \? WHERE id = \?`

	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	mockDB.ExpectPrepare(query)

	repository := NewRepository(db)
	err := repository.UpdateUserAgeLastname(1, "lastname", 24)

	assert.Error(t, err)
}

func Test_DB_Mock_UpdateUserFirstname_Success(t *testing.T) {
	query := `UPDATE users SET firstname = \? WHERE id = \?`

	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	mockDB.ExpectPrepare(query).
		WillBeClosed().
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewRepository(db)
	err := repository.UpdateUserFirstname(1, "firstname")

	assert.Nil(t, err)
}

func Test_DB_Mock_UpdateUserFirstname_Fail_Prepare(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	repository := NewRepository(db)
	err := repository.UpdateUserFirstname(1, "firstname")

	assert.Error(t, err)
}

func Test_DB_Mock_UpdateUserFirstname_Fail_Exec(t *testing.T) {
	query := `UPDATE users SET firstname = \? WHERE id = \?`

	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	mockDB.ExpectPrepare(query)

	repository := NewRepository(db)
	err := repository.UpdateUserFirstname(1, "firstname")
	assert.Error(t, err)
}

func Test_DB_Mock_Store_Success(t *testing.T) {
	query := `INSERT INTO users \(firstname, lastname, email, age, height, active, created_at\) VALUES\( \?, \?, \?, \?, \?, \?, \?\)`

	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	mockDB.ExpectPrepare(query).
		WillBeClosed().
		ExpectExec().
		WillReturnResult(sqlmock.NewResult(1, 1))

	repository := NewRepository(db)
	actualUser, err := repository.Store(model.User{Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"})

	expectedUser := model.User{ID: 1, Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"}

	assert.Nil(t, err)
	assert.Equal(t, expectedUser, actualUser)
}

func Test_DB_Mock_Store_Fail_Prepare(t *testing.T) {
	db, _, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	repository := NewRepository(db)
	actualUser, err := repository.Store(model.User{Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"})
	assert.Error(t, err)
	assert.Equal(t, model.User{}, actualUser)
}

func Test_DB_Mock_Store_Fail_Exec(t *testing.T) {
	db, mockDB, _ := sqlmock.New()
	defer func(db *sql.DB) {
		_ = db.Close()
	}(db)

	mockDB.ExpectPrepare(storeQuery)

	repository := NewRepository(db)
	actualUser, err := repository.Store(model.User{Firstname: "firstname1", Lastname: "lastname1", Email: "email1", Age: 24, Height: 184, Active: true, CreatedDate: "22/02/2021"})

	assert.Error(t, err)
	assert.Equal(t, model.User{}, actualUser)
}
