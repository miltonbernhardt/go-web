package users

import (
	"errors"
	"github.com/miltonbernhardt/go-web/internal/model"
	"github.com/miltonbernhardt/go-web/internal/utils"
	"github.com/miltonbernhardt/go-web/pkg/message"
	"reflect"
	"strings"
)

type Service interface {
	Delete(id int) error
	GetAll() ([]model.User, error)
	GetAllWithFilters(users []model.User, attribute model.UserTypes, value interface{}) []model.User
	GetByID(id int) (model.User, error)
	Store(user model.User) (model.User, error)
	UpdateFields(id int, lastname string, age int) error
	Update(id int, user model.User) (model.User, error)
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

func (s *service) GetAll() ([]model.User, error) {
	return s.repository.GetAll()
}

func (s *service) GetByID(id int) (model.User, error) {
	allUsers, err := s.GetAll()

	if err != nil {
		return model.User{}, err
	}

	for _, user := range allUsers {
		if user.ID == id {
			return user, nil
		}
	}

	return model.User{}, errors.New(message.UserNotFound)
}

func (s *service) GetAllWithFilters(users []model.User, fieldType model.UserTypes, value interface{}) []model.User {
	var sliceUsers []model.User
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

func (s *service) Store(user model.User) (model.User, error) {
	user.CreatedDate = s.utils.GetNowAsString()
	return s.repository.Store(user)
}

/*####################### PUT #######################*/

func (s *service) Update(id int, user model.User) (model.User, error) {
	//todo change model.user to some dto
	return s.repository.Update(id, user)
}

/*####################### DELETE #######################*/

func (s *service) Delete(id int) error {
	return s.repository.Delete(id)
}

/*####################### PATCH #######################*/

func (s *service) UpdateFields(id int, lastname string, age int) error {
	return s.repository.UpdateUser(id, lastname, age)
}
