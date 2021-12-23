package users

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

type User struct {
	Id          int64  `json:"id"`
	Firstname   string `json:"firstname" binding:"required"`
	Lastname    string `json:"lastname" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Age         int64  `json:"age" binding:"required"`
	Height      int64  `json:"height" binding:"required"`
	Active      bool   `json:"active"` //todo problemas para required bool con valor false
	CreatedDate string `json:"created_date"`
}

var users []User

type Repository interface {
	GetAll() ([]User, error)
	Store(user User) (User, error)
}

type repository struct{}

func NewRepository() Repository {
	return &repository{}
}



func (r *repository) getUsersAsJson() ([]byte, error) {
	usersJson, err := os.ReadFile("./users.json")

	if err != nil {
		return nil, errors.New("no se pudo leer el archivo")
	}

	return usersJson, nil
}

func (r *repository) GetAll() ([]User, error) {
	if users == nil {
		usersJson, err := r.getUsersAsJson()

		if err != nil {
			return nil, err
		}

		err = json.Unmarshal(usersJson, &users)

		if err != nil {
			return nil, err
		}
		fmt.Println(users)
	}

	return users, nil
}

func (r *repository) Store(user User) (User, error) {
	lastIndex := len(users) - 1
	lastId := (users[lastIndex].Id) + 1
	user.Id = lastId

	users = append(users, user)

	return user, nil
}
