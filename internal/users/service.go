package users

import (
	"fmt"
	"github.com/miltonbernhardt/go-web/internal/domain"
	"reflect"
	"strings"
)

type Service interface {
	DeleteUser(id int64) error
	FetchAllUsers() ([]domain.User, error)
	FetchAllUsersByQuery(users []domain.User, attribute domain.UserTypes, value interface{}) []domain.User
	FetchUserByID(id int64) (domain.User, error)
	StoreUser(user domain.User) (domain.User, error)
	UpdateFieldsUser(id int64, lastname string, age int64) (domain.User, error)
	UpdateUser(id int64, user domain.User) (domain.User, error)
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

func (s *service) FetchAllUsers() ([]domain.User, error) {
	return s.repository.FetchAllUsers()
}

func (s *service) FetchUserByID(id int64) (domain.User, error) {
	users, err := s.FetchAllUsers()

	if err != nil {
		return domain.User{}, err
	}

	for _, user := range users {
		if user.Id == id {
			return user, nil
		}
	}

	return domain.User{}, fmt.Errorf("no se encontró el usuario %v, bien puede estar eliminado o dado de baja", id)
}

func (s *service) FetchAllUsersByQuery(users []domain.User, fieldType domain.UserTypes, value interface{}) []domain.User {
	var sliceUsers []domain.User
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

	return sliceUsers
}

/*####################### POST #######################*/

func (s *service) StoreUser(user domain.User) (domain.User, error) {
	return s.repository.Store(user)
}

/*####################### PUT #######################*/

func (s *service) UpdateUser(id int64, user domain.User) (domain.User, error) {
	return s.repository.Update(id, user)
}

/*####################### DELETE #######################*/

func (s *service) DeleteUser(id int64) error {
	return s.repository.DeleteUser(id)
}

/*####################### PATCH #######################*/

func (s *service) UpdateFieldsUser(id int64, lastname string, age int64) (domain.User, error) {
	return s.repository.UpdateUser(id, lastname, age)
}
