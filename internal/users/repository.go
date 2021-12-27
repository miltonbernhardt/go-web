package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/miltonbernhardt/go-web/internal/domain"
	"os"
)

var users []domain.User

type Repository interface {
	GetAll() ([]domain.User, error)
	Store(user domain.User) domain.User
	Update(id int64, user domain.User) (domain.User, error)
	Delete(id int64) error
	UpdateFields(id int64, lastname string, age int64) (domain.User, error)
}

type repository struct{}

func NewRepository() Repository {
	repo := &repository{}
	_, _ = repo.GetAll()
	return repo
}

func (r *repository) getUsersAsJson() ([]byte, error) {
	usersJson, err := os.ReadFile("./users.json")

	if err != nil {
		return nil, errors.New("no se pudo leer el archivo")
	}

	return usersJson, nil
}

func (r *repository) GetAll() ([]domain.User, error) {
	if users == nil {
		usersJson, err := r.getUsersAsJson()

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(usersJson, &users)

		if err != nil {
			return nil, err
		}
	}

	return users, nil
}

func (r *repository) Store(user domain.User) domain.User {
	user.CreatedDate = GetNowAsString()

	lastIndex := len(users) - 1
	lastId := (users[lastIndex].Id) + 1
	user.Id = lastId

	users = append(users, user)

	return user
}

func (r *repository) Update(id int64, userToUpdate domain.User) (domain.User, error) {
	for i := range users {
		if users[i].Id == id {
			users[i].Active = userToUpdate.Active
			users[i].Age = userToUpdate.Age
			users[i].Email = userToUpdate.Email
			users[i].Firstname = userToUpdate.Firstname
			users[i].Height = userToUpdate.Height
			users[i].Lastname = userToUpdate.Lastname
			return users[i], nil
		}
	}

	return domain.User{}, fmt.Errorf("no se encontro un usuario con dicho id = %v", id)

}

func (r *repository) Delete(id int64) error {
	for i := range users {
		if users[i].Id == id && users[i].DeletedDate == "" {
			users[i].DeletedDate = GetNowAsString()
			return nil
		}
	}

	return fmt.Errorf("usuario %d no encontrado", id)
}

func (r *repository) UpdateFields(id int64, lastname string, age int64) (domain.User, error) {
	for i := range users {
		if users[i].Id == id {
			if age != 0 {
				users[i].Age = age
			}
			if lastname != "" {
				users[i].Lastname = lastname
			}
			return users[i], nil
		}
	}

	return domain.User{}, fmt.Errorf("no se encontro un usuario con dicho id = %v", id)

}
