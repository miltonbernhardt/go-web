package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/miltonbernhardt/go-web/internal/model"
	"github.com/miltonbernhardt/go-web/pkg/message"
	log "github.com/sirupsen/logrus"
)

type Repository interface {
	Delete(id int) error
	GetAll() ([]model.User, error)
	GetByFirstname(firstname string) (model.User, error)
	Store(user model.User) (model.User, error)
	Update(user model.User) (model.User, error)
	UpdateUserAgeLastname(id int, lastname string, age int) error
	UpdateUserFirstname(id int, firstname string) error
	UpdateWithContext(ctx context.Context, user model.User) (model.User, error)
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
	rows, err := r.db.Query(getAllQuery)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	users = []model.User{}
	for rows.Next() {
		var user model.User
		if err = rows.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Age, &user.Height, &user.Active, &user.CreatedDate); err != nil {
			log.Error(err)
			return nil, err
		}

		users = append(users, user)
	}
	return users, nil
}

func (r *repositoryDB) GetByFirstname(firstname string) (model.User, error) {
	row := r.db.QueryRow(getByFirstnameQuery, firstname)

	var user model.User

	if err := row.Scan(&user.ID, &user.Firstname, &user.Lastname, &user.Email, &user.Age, &user.Height, &user.Active, &user.CreatedDate); err != nil {
		log.Error(err)
		return model.User{}, errors.New(message.UserNotFound)
	}
	return user, nil
}

func (r *repositoryDB) Delete(id int) error {
	stmt, err := r.db.Prepare(deleteQuery)
	if err != nil {
		log.Error(err)
		return err
	}
	defer r.stmtClose(stmt)

	_, err = stmt.Exec(id)

	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryDB) Update(user model.User) (model.User, error) {
	stmt, err := r.db.Prepare(updateQuery)
	if err != nil {
		log.Error(err)
		return model.User{}, err
	}
	defer r.stmtClose(stmt)

	_, err = stmt.Exec(user.Firstname, user.Lastname, user.Email, user.Age, user.Height, user.Active, user.CreatedDate, user.ID)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *repositoryDB) UpdateUserAgeLastname(id int, lastname string, age int) error {
	stmt, err := r.db.Prepare(updateUserAgeLastnameQuery)
	if err != nil {
		log.Error(err)
		return err
	}
	defer r.stmtClose(stmt)

	_, err = stmt.Exec(lastname, age, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryDB) UpdateUserFirstname(id int, firstname string) error {
	stmt, err := r.db.Prepare(updateUserFirstnameQuery)
	if err != nil {
		log.Error(err)
		return err
	}
	defer r.stmtClose(stmt)

	_, err = stmt.Exec(firstname, id)
	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryDB) UpdateWithContext(ctx context.Context, user model.User) (model.User, error) {
	stmt, err := r.db.Prepare(updateQuery)
	if err != nil {
		log.Error(err)
		return model.User{}, err
	}
	defer r.stmtClose(stmt)

	_, err = stmt.ExecContext(ctx, user.Firstname, user.Lastname, user.Email, user.Age, user.Height, user.Active, user.CreatedDate, user.ID)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *repositoryDB) Store(user model.User) (model.User, error) {
	stmt, err := r.db.Prepare(storeQuery)
	if err != nil {
		log.Error(err)
		return model.User{}, err
	}
	defer r.stmtClose(stmt)

	result, err := stmt.Exec(user.Firstname, user.Lastname, user.Email, user.Age, user.Height, user.Active, user.CreatedDate)
	if err != nil {
		return model.User{}, err
	}
	insertedId, _ := result.LastInsertId()
	user.ID = int(insertedId)
	return user, nil
}

func (r *repositoryDB) stmtClose(stmt *sql.Stmt) {
	_ = stmt.Close()
}
