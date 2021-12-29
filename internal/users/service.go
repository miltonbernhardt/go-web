package users

import (
	"fmt"
	"github.com/miltonbernhardt/go-web/internal/domain"
	"reflect"
	"strings"
)

type Service interface {
	GetAll() ([]domain.User, error)
	GetAllByField(users []domain.User, attribute domain.UserTypes, value interface{}) []domain.User
	GetByID(id int64) (domain.User, error)
	Store(user domain.User) (domain.User, error)
	Update(id int64, user domain.User) (domain.User, error)
	UpdateFields(id int64, lastname string, age int64) (domain.User, error)
	Delete(id int64) error
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

func (s *service) GetAll() ([]domain.User, error) {
	return s.repository.GetAll()
}

func (s *service) GetByID(id int64) (domain.User, error) {
	users, err := s.GetAll()

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

func (s *service) GetAllByField(users []domain.User, fieldType domain.UserTypes, value interface{}) []domain.User {
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

func (s *service) Store(user domain.User) (domain.User, error) {
	return s.repository.Store(user)
}

/*####################### PUT #######################*/

func (s *service) Update(id int64, user domain.User) (domain.User, error) {
	return s.repository.Update(id, user)
}

/*####################### DELETE #######################*/

func (s *service) Delete(id int64) error {
	return s.repository.Delete(id)
}

/*####################### PATCH #######################*/

func (s *service) UpdateFields(id int64, lastname string, age int64) (domain.User, error) {
	return s.repository.UpdateFields(id, lastname, age)
}
