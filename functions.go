package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
)

func GetUsersAsJson() ([]byte, error) {
	usersJson, err := os.ReadFile("./users.json")

	if err != nil {
		return nil, errors.New("error: no se pudo leer el archivo")
	}

	return usersJson, nil
}

func GetUsersAsSlice() ([]User, error) {
	var users []User
	usersJson, err := GetUsersAsJson()

	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(usersJson, &users)

	if err != nil {
		return nil, err
	}

	fmt.Println(users)
	return users, nil
}

func GetUserByID(users []User, id string) (User, error) {
	for _, user := range users {
		if user.Id == id {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("error: no se encontro un usuario con dicho id = %v", id)
}

func GetUsersByFirstName(users []User, firstname string) []User {
	var sliceUsers []User

	for _, user := range users {
		if strings.Contains(user.Firstname, firstname) {
			sliceUsers = append(sliceUsers, user)
		}
	}

	return sliceUsers
}

func GetUsersByLastname(users []User, lastname string) []User {
	var sliceUsers []User

	for _, user := range users {
		if strings.Contains(user.Lastname, lastname) {
			sliceUsers = append(sliceUsers, user)
		}
	}
	return sliceUsers
}

func GetUsersByEmail(users []User, email string) []User {
	var sliceUsers []User

	for _, user := range users {
		if strings.Contains(user.Email, email) {
			sliceUsers = append(sliceUsers, user)
		}
	}
	return sliceUsers
}

func GetUsersByCreatedDate(users []User, createdDate string) []User {
	var sliceUsers []User

	for _, user := range users {
		if user.CreatedDate == createdDate {
			sliceUsers = append(sliceUsers, user)
		}
	}
	return sliceUsers
}

func GetUsersByIsActive(users []User, isActive bool) []User {
	var sliceUsers []User

	for _, user := range users {
		if user.Active == isActive {
			sliceUsers = append(sliceUsers, user)
		}
	}
	return sliceUsers
}

func GetUsersByAge(users []User, age int64) []User {
	var sliceUsers []User

	for _, user := range users {
		if user.Age == age {
			sliceUsers = append(sliceUsers, user)
		}
	}
	return sliceUsers
}

func GetUsersByHeight(users []User, height int64) []User {
	var sliceUsers []User

	for _, user := range users {
		if user.Height == height {
			sliceUsers = append(sliceUsers, user)
		}
	}
	return sliceUsers
}
