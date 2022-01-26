package users

import (
	"database/sql"
	"github.com/miltonbernhardt/go-web/internal/model"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	Delete(id int) error
	GetAll() ([]model.User, error)
	Store(user model.User) (model.User, error)
	Update(id int, user model.User) (model.User, error)
	UpdateName(id int, name string) error
	UpdateUser(id int, lastname string, age int) error
}

type repositoryDB struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) Repository {
	return &repositoryDB{
		db: db,
	}
}

func (r *repositoryDB) GetAll() (users []model.User, err error) {
	return []model.User{}, nil
}

func (r *repositoryDB) Delete(id int) error {
	return nil
}

func (r *repositoryDB) Update(id int, user model.User) (model.User, error) {
	return user, nil
}

func (r *repositoryDB) UpdateName(id int, name string) error {
	return nil
}

func (r *repositoryDB) UpdateUser(id int, lastname string, age int) error {
	stmt, err := r.db.Prepare("UPDATE users SET lastname = ?, age = ? WHERE id = ?")
	if err != nil {
		log.Error(err)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	_, err = stmt.Exec(lastname, age, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryDB) Store(user model.User) (model.User, error) {
	stmt, err := r.db.Prepare("INSERT INTO users (username, lastname, email, age, height, active, created_at) VALUES( ?, ?, ?, ?, ?, ?, ?)")
	if err != nil {
		log.Error(err)
	}
	defer func(stmt *sql.Stmt) {
		_ = stmt.Close()
	}(stmt)

	result, err := stmt.Exec(user.Firstname, user.Lastname, user.Email, user.Age, user.Height, user.Active, user.CreatedDate)
	if err != nil {
		return model.User{}, err
	}
	insertedId, _ := result.LastInsertId()
	user.ID = int(insertedId)
	return user, nil
}
