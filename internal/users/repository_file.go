package users

import (
	"context"
	"errors"
	"github.com/miltonbernhardt/go-web/internal/model"
	"github.com/miltonbernhardt/go-web/pkg/message"
	"github.com/miltonbernhardt/go-web/pkg/store"
	log "github.com/sirupsen/logrus"
)

type repositoryFile struct {
	db store.Store
}

func NewRepositoryFile(db store.Store) Repository {
	return &repositoryFile{db: db}
}

func (r *repositoryFile) GetAll() (users []model.User, err error) {
	defer func() {
		panicError := recover()

		if panicError != nil {
			users = nil
			err = panicError.(error)
		}
	}()

	err = r.db.Read(&users)
	if err != nil {
		return []model.User{}, err
	}

	return users, nil
}

func (r *repositoryFile) GetByFirstname(firstname string) (model.User, error) {
	allUsers, err := r.GetAll()

	if err != nil {
		return model.User{}, err
	}

	for i := range allUsers {
		if allUsers[i].Firstname == firstname {
			return allUsers[i], nil
		}
	}

	return model.User{}, errors.New("user not found")
}

func (r *repositoryFile) Delete(id int) error {
	allUsers, err := r.GetAll()

	if err != nil {
		return err
	}

	index := -1

	for i := range allUsers {
		if allUsers[i].ID == id {
			index = i
			break
		}
	}

	if index == -1 {
		return errors.New(message.UserNotFound)
	} else {
		allUsers = append(allUsers[:index], allUsers[index+1:]...)
		err = r.db.Write(&allUsers)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *repositoryFile) Update(userToUpdate model.User) (model.User, error) {
	allUsers, err := r.GetAll()
	if err != nil {
		return model.User{}, err
	}

	user := model.User{}

	for i := range allUsers {
		if allUsers[i].ID == userToUpdate.ID {
			allUsers[i].Active = userToUpdate.Active
			allUsers[i].Age = userToUpdate.Age
			allUsers[i].Email = userToUpdate.Email
			allUsers[i].Firstname = userToUpdate.Firstname
			allUsers[i].Height = userToUpdate.Height
			allUsers[i].Lastname = userToUpdate.Lastname
			user = allUsers[i]
			break
		}
	}

	if user.ID == 0 {
		return model.User{}, errors.New(message.UserNotFound)
	}

	err = r.db.Write(allUsers)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *repositoryFile) UpdateUserFirstname(id int, name string) error {
	allUsers, err := r.GetAll()
	if err != nil {
		return err
	}

	updated := false

	for i := range allUsers {
		if allUsers[i].ID == id {
			allUsers[i].Firstname = name
			updated = true
		}
	}

	if !updated {
		return errors.New(message.UserNotFound)
	}
	if err := r.db.Write(allUsers); err != nil {
		return err
	}
	return nil
}

func (r *repositoryFile) UpdateUserAgeLastname(id int, lastname string, age int) error {
	allUsers, err := r.GetAll()
	if err != nil {
		return err
	}

	user := model.User{}

	for i := range allUsers {
		if allUsers[i].ID == id {
			if age != 0 {
				allUsers[i].Age = age
			}
			if lastname != "" {
				allUsers[i].Lastname = lastname
			}
			user = allUsers[i]
			break
		}
	}

	if user.ID == 0 {
		return errors.New(message.UserNotFound)
	}

	err = r.db.Write(&allUsers)
	if err != nil {
		return err
	}

	return nil
}

func (r *repositoryFile) UpdateWithContext(ctx context.Context, user model.User) (model.User, error) {
	return r.Update(user)
}

func (r *repositoryFile) Store(user model.User) (model.User, error) {
	allUsers, err := r.GetAll()
	if err != nil {
		return model.User{}, err
	}

	lastID, _ := r.getUserLastID()
	user.ID = lastID + 1

	if err != nil {
		return model.User{}, err
	}

	allUsers = append(allUsers, user)

	err = r.db.Write(&allUsers)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (r *repositoryFile) getUserLastID() (int, error) {
	allUsers, err := r.GetAll()
	if err != nil {
		log.Error(err)
		return 0, err
	}

	if len(allUsers) == 0 {
		return 0, nil
	}

	return allUsers[len(allUsers)-1].ID, nil
}
