package users

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

var TokenAuth = "bearer 12345"

type Service interface {
	GetByID(id int64) (User, error)
	GetAll() ([]User, error)
	Store(user User) (User, error)
	GetAllByField(users []User, attribute UserTypes, value interface{}) []User
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

func (s *service) GetAllByField(users []User, fieldType UserTypes, value interface{}) []User {
	var sliceUsers []User
	fmt.Printf("\tCampo: %s - valor: %v\n", fieldType, value)

	for _, user := range users {

		userReflected := reflect.ValueOf(&user)
		userReflected = userReflected.Elem()
		field := userReflected.FieldByName(string(fieldType))

		var valueToCompare interface{}
		if field.IsValid() && field.CanSet() {
			switch field.Kind() {
			case reflect.Int64:
				valueToCompare = field.Interface().(int64)
				if valueToCompare == value {
					sliceUsers = append(sliceUsers, user)
				}
			case reflect.Bool:
				valueToCompare = field.Interface().(bool)
				if valueToCompare == value {
					sliceUsers = append(sliceUsers, user)
				}
			case reflect.String:
				valueToCompare = field.Interface().(string)
				if valueToCompare == value || strings.Contains(valueToCompare.(string), value.(string)) {
					sliceUsers = append(sliceUsers, user)
				}
			}

		}
	}

	fmt.Printf("\tcantidad de usuarios que cumplen: %v\n", len(sliceUsers))
	return sliceUsers
}

/*####################### POST #######################*/

func (s *service) Store(user User) (User, error) {
	t := time.Now()
	user.CreatedDate = t.Format("02/01/2006")
	fmt.Printf("\nNew user: %v\n", user)

	return s.repository.Store(user)
}
