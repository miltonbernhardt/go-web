package users

import (
	"errors"
	"fmt"
	"github.com/miltonbernhardt/go-web/internal/domain"
	"github.com/miltonbernhardt/go-web/pkg/store"
	"github.com/miltonbernhardt/go-web/pkg/web"
)

type Repository interface {
	DeleteUser(id int) error
	GetAll() ([]domain.User, error)
	Store(user domain.User) (domain.User, error)
	Update(id int, user domain.User) (domain.User, error)
	UpdateName(id int, name string) (domain.User, error)
	UpdateUser(id int, lastname string, age int) (domain.User, error)
	getUserLastID() (int, error)
}

type repository struct {
	db store.Store
}

func NewRepository(db store.Store) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() (users []domain.User, err error) {
	defer func() {
		panicError := recover()

		if panicError != nil {
			fmt.Printf("\n\n\t\t-------------------\n\n")
			users = nil
			err = panicError.(error)
		}
	}()

	err = r.db.Read(&users)
	if err != nil {
		return []domain.User{}, err
	}

	return users, nil
}

func (r *repository) getUserLastID() (int, error) {
	users, err := r.GetAll()
	if err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, nil
	}

	return users[len(users)-1].ID, nil
}

func (r *repository) DeleteUser(id int) error {
	users, err := r.GetAll()

	if err != nil {
		return err
	}

	index := -1

	for i := range users {
		if users[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New(web.UserNotFound)
	} else {
		users = append(users[:index], users[index+1:]...)
		err = r.db.Write(&users)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repository) Update(id int, userToUpdate domain.User) (domain.User, error) {
	users, err := r.GetAll()
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{}

	for i := range users {
		if users[i].ID == id {
			users[i].Active = userToUpdate.Active
			users[i].Age = userToUpdate.Age
			users[i].Email = userToUpdate.Email
			users[i].Firstname = userToUpdate.Firstname
			users[i].Height = userToUpdate.Height
			users[i].Lastname = userToUpdate.Lastname
			user = users[i]
			break
		}
	}

	if user.ID == 0 {
		return domain.User{}, errors.New(web.UserNotFound)
	}

	err = r.db.Write(users)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *repository) UpdateName(id int, name string) (domain.User, error) {

	users, err := r.GetAll()
	if err != nil {
		return domain.User{}, err
	}

	var p domain.User
	updated := false

	for i := range users {
		if users[i].ID == id {
			users[i].Firstname = name
			updated = true
			p = users[i]
		}
	}

	if !updated {
		return domain.User{}, errors.New(web.UserNotFound)
	}
	if err := r.db.Write(users); err != nil {
		return domain.User{}, err
	}
	return p, nil
}

func (r *repository) UpdateUser(id int, lastname string, age int) (domain.User, error) {
	users, err := r.GetAll()
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{}

	for i := range users {
		if users[i].ID == id {
			if age != 0 {
				users[i].Age = age
			}
			if lastname != "" {
				users[i].Lastname = lastname
			}
			user = users[i]
			break
		}
	}

	if user.ID == 0 {
		return domain.User{}, errors.New(web.UserNotFound)
	}

	err = r.db.Write(&users)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *repository) Store(user domain.User) (domain.User, error) {
	users, err := r.GetAll()
	if err != nil {
		return domain.User{}, err
	}

	lastID, _ := r.getUserLastID()
	user.ID = lastID + 1

	if err != nil {
		return domain.User{}, err
	}

	users = append(users, user)

	err = r.db.Write(&users)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
