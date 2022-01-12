package users

import (
	"errors"
	"github.com/miltonbernhardt/go-web/internal/domain"
	"github.com/miltonbernhardt/go-web/internal/utils"
	"github.com/miltonbernhardt/go-web/pkg/web"
	"reflect"
	"strings"
)

type Service interface {
	DeleteUser(id int) error
	GetAll() ([]domain.User, error)
	FetchAllUsersByQuery(users []domain.User, attribute domain.UserTypes, value interface{}) []domain.User
	FetchUserByID(id int) (domain.User, error)
	StoreUser(user domain.User) (domain.User, error)
	UpdateFieldsUser(id int, lastname string, age int) (domain.User, error)
	UpdateUser(id int, user domain.User) (domain.User, error)
}
type service struct {
	repository Repository
	utils      utils.Functions
}

func NewService(r Repository, u utils.Functions) Service {
	return &service{
		repository: r,
		utils:      u,
	}
}

/*####################### GET #######################*/

func (s *service) GetAll() ([]domain.User, error) {
	return s.repository.GetAll()
}

func (s *service) FetchUserByID(id int) (domain.User, error) {
	users, err := s.GetAll()

	if err != nil {
		return domain.User{}, err
	}

	for _, user := range users {
		if user.ID == id {
			return user, nil
		}
	}

	return domain.User{}, errors.New(web.UserNotFound)
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
			case reflect.Int:
				valueToCompare = field.Interface().(int)
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
	user.CreatedDate = s.utils.GetNowAsString()
	return s.repository.Store(user)
}

/*####################### PUT #######################*/

func (s *service) UpdateUser(id int, user domain.User) (domain.User, error) {
	return s.repository.Update(id, user)
}

/*####################### DELETE #######################*/

func (s *service) DeleteUser(id int) error {
	return s.repository.DeleteUser(id)
}

/*####################### PATCH #######################*/

func (s *service) UpdateFieldsUser(id int, lastname string, age int) (domain.User, error) {
	return s.repository.UpdateUser(id, lastname, age)
}
