package users

import (
	"database/sql"
	"github.com/miltonbernhardt/go-web/internal/model"
)

type Repository interface {
	Delete(id int) error
	GetAll() ([]model.User, error)
	Store(user model.User) (model.User, error)
	Update(id int, user model.User) (model.User, error)
	UpdateName(id int, name string) (model.User, error)
	UpdateUser(id int, lastname string, age int) (model.User, error)
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

func (r *repositoryDB) Update(id int, userToUpdate model.User) (model.User, error) {
	return model.User{}, nil
}

func (r *repositoryDB) UpdateName(id int, name string) (model.User, error) {
	return model.User{}, nil
}

func (r *repositoryDB) UpdateUser(id int, lastname string, age int) (model.User, error) {
	return model.User{}, nil
}

func (r *repositoryDB) Store(user model.User) (model.User, error) {
	return model.User{}, nil
}

func (r *repositoryDB) getUserLastID() (int, error) {
	return 0, nil
}
