package users

import (
	"fmt"
	"github.com/miltonbernhardt/go-web/internal/domain"
	"github.com/miltonbernhardt/go-web/pkg/store"
)

type Repository interface {
	GetAll() ([]domain.User, error)
	GetLastID() (int64, error)
	Store(user domain.User) (domain.User, error)
	Update(id int64, user domain.User) (domain.User, error)
	Delete(id int64) error
	UpdateFields(id int64, lastname string, age int64) (domain.User, error)
}

type repository struct {
	db store.Store
}

func NewRepository(db store.Store) Repository {
	return &repository{db: db}
}

func (r *repository) GetAll() ([]domain.User, error) {
	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *repository) GetLastID() (int64, error) {
	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, nil
	}

	return users[len(users)-1].Id, nil
}

func (r *repository) Delete(id int64) error {
	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return err
	}

	isDeleted := false

	for i := range users {
		if users[i].Id == id {
			if users[i].DeletedDate == "" {
				users[i].DeletedDate = GetNowAsString()
				isDeleted = true
			}
			break
		}
	}

	err = r.db.Write(&users)
	if err != nil {
		return err
	}

	if !isDeleted {
		return fmt.Errorf("usuario %d no encontrado", id)
	}

	return nil
}

func (r *repository) Update(id int64, userToUpdate domain.User) (domain.User, error) {
	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{}

	for i := range users {
		if users[i].Id == id {
			if users[i].DeletedDate == "" {
				users[i].Active = userToUpdate.Active
				users[i].Age = userToUpdate.Age
				users[i].Email = userToUpdate.Email
				users[i].Firstname = userToUpdate.Firstname
				users[i].Height = userToUpdate.Height
				users[i].Lastname = userToUpdate.Lastname
				user = users[i]
			}
			break
		}
	}

	if user.Id == 0 {
		return domain.User{}, fmt.Errorf("no se encontró el usuario %v, puede no existir o estar dado de baja", id)
	}

	err = r.db.Write(users)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *repository) UpdateFields(id int64, lastname string, age int64) (domain.User, error) {
	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{}

	for i := range users {
		if users[i].Id == id {
			if users[i].DeletedDate == "" {
				if age != 0 {
					users[i].Age = age
				}
				if lastname != "" {
					users[i].Lastname = lastname
				}
				user = users[i]
			}
			break
		}
	}

	if user.Id == 0 {
		return domain.User{}, fmt.Errorf("no se encontró el usuario %v", id)
	}

	err = r.db.Write(&users)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *repository) Store(user domain.User) (domain.User, error) {
	user.CreatedDate = GetNowAsString()

	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return domain.User{}, err
	}

	lastID, err := r.GetLastID()
	user.Id = lastID + 1

	if err != nil {
		return domain.User{}, err
	}

	users = append(users, user)

	err = r.db.Write(&users)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}
