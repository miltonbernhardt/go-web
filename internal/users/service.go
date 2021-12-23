package users

import (
	"fmt"
	"time"
)

var TokenAuth = "bearer 12345"

type Service interface {
	GetByID(id int64) (User, error)
	GetAll() ([]User, error)
	Store(user User) (User, error)
	//todo añadir los nuevos métodos
}
type service struct {
	repository Repository
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}

/*####################### GET #######################*/
// todo generalizar los get username by field - usar reflect

func (s *service) GetAll() ([]User, error) {
	return s.repository.GetAll()
}

func (s *service) GetByID(id int64) (User, error) {
	users, err := s.GetAll()

	if err != nil {
		return User{}, err
	}

	for _, user := range users {
		if user.Id == id {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("no se encontro un usuario con dicho id = %v", id)
}

//func (s *service) GetByFirstName(users []User, firstname string) []User {
//	var sliceUsers []User
//
//	for _, user := range users {
//		if strings.Contains(user.Firstname, firstname) {
//			sliceUsers = append(sliceUsers, user)
//		}
//	}
//
//	return sliceUsers
//}
//
//func (s *service) GetByLastname(users []User, lastname string) []User {
//	var sliceUsers []User
//
//	for _, user := range users {
//		if strings.Contains(user.Lastname, lastname) {
//			sliceUsers = append(sliceUsers, user)
//		}
//	}
//	return sliceUsers
//}
//
//func (s *service) GetByEmail(users []User, email string) []User {
//	var sliceUsers []User
//
//	for _, user := range users {
//		if strings.Contains(user.Email, email) {
//			sliceUsers = append(sliceUsers, user)
//		}
//	}
//	return sliceUsers
//}
//
//func (s *service) GetByCreatedDate(users []User, createdDate string) []User {
//	var sliceUsers []User
//
//	for _, user := range users {
//		if user.CreatedDate == createdDate {
//			sliceUsers = append(sliceUsers, user)
//		}
//	}
//	return sliceUsers
//}
//
//func (s *service) GetByIsActive(users []User, isActive bool) []User {
//	var sliceUsers []User
//
//	for _, user := range users {
//		if user.Active == isActive {
//			sliceUsers = append(sliceUsers, user)
//		}
//	}
//	return sliceUsers
//}
//
//func (s *service) GetByAge(users []User, age int64) []User {
//	var sliceUsers []User
//
//	for _, user := range users {
//		if user.Age == age {
//			sliceUsers = append(sliceUsers, user)
//		}
//	}
//	return sliceUsers
//}
//
//func (s *service) GetByHeight(users []User, height int64) []User {
//	var sliceUsers []User
//
//	for _, user := range users {
//		if user.Height == height {
//			sliceUsers = append(sliceUsers, user)
//		}
//	}
//	return sliceUsers
//}

/*####################### POST #######################*/

func (s *service) Store(user User) (User, error) {
	t := time.Now()
	user.CreatedDate = t.Format("02/01/2006")
	fmt.Printf("\nNew user: %v\n", user)

	return s.repository.Store(user)
}
