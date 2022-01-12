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

func (r *repository) GetAll() ([]domain.User, error) {
	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) getUserLastID() (int, error) {
	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, nil
	}

	return users[len(users)-1].ID, nil
}

func (r *repository) DeleteUser(id int) error {
	var users []domain.User

	err := r.db.Read(&users)
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
		return errors.New(web.UserNotExists)
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
	var users []domain.User

	err := r.db.Read(&users)
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

	var ps []domain.User
	if err := r.db.Read(&ps); err != nil {
		return domain.User{}, err
	}

	var p domain.User
	updated := false

	for i := range ps {
		if ps[i].ID == id {
			ps[i].Firstname = name
			updated = true
			p = ps[i]
		}
	}

	if !updated {
		return domain.User{}, fmt.Errorf("user %d not found", id)
	}
	if err := r.db.Write(ps); err != nil {
		return domain.User{}, err
	}
	return p, nil
}

func (r *repository) UpdateUser(id int, lastname string, age int) (domain.User, error) {
	var users []domain.User

	err := r.db.Read(&users)
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
	user.CreatedDate = GetNowAsString()

	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return domain.User{}, err
	}

	lastID, err := r.getUserLastID()
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
