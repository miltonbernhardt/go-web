package users

import (
	"fmt"
	"github.com/miltonbernhardt/go-web/internal/domain"
	"github.com/miltonbernhardt/go-web/pkg/store"
)

type Repository interface {
	DeleteUser(id int64) error
	GetAll() ([]domain.User, error)
	Store(user domain.User) (domain.User, error)
	Update(id int64, user domain.User) (domain.User, error)
	UpdateName(id int64, name string) (domain.User, error)
	UpdateUser(id int64, lastname string, age int64) (domain.User, error)
	getUserLastID() (int64, error)
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

func (r *repository) getUserLastID() (int64, error) {
	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return 0, err
	}

	if len(users) == 0 {
		return 0, nil
	}

	return users[len(users)-1].ID, nil
}

func (r *repository) DeleteUser(id int64) error {
	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return err
	}

	isDeleted := false

	for i := range users {
		if users[i].ID == id {
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
		if users[i].ID == id {
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

	if user.ID == 0 {
		return domain.User{}, fmt.Errorf("no se encontró el usuario %v, puede no existir o estar dado de baja", id)
	}

	err = r.db.Write(users)
	if err != nil {
		return domain.User{}, err
	}

	return user, nil
}

func (r *repository) UpdateName(id int64, name string) (domain.User, error) {

	var ps []domain.User
	if err := r.db.Read(&ps); err != nil {
		return domain.User{}, err
	}

	var p domain.User
	updated := false

	for i := range ps {
		if ps[i].ID == id {
			ps[i].Firstname = name
			updated = true
			p = ps[i]
		}
	}

	if !updated {
		return domain.User{}, fmt.Errorf("user %d not found", id)
	}
	if err := r.db.Write(ps); err != nil {
		return domain.User{}, err
	}
	return p, nil
}

func (r *repository) UpdateUser(id int64, lastname string, age int64) (domain.User, error) {
	var users []domain.User

	err := r.db.Read(&users)
	if err != nil {
		return domain.User{}, err
	}

	user := domain.User{}

	for i := range users {
		if users[i].ID == id {
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

	if user.ID == 0 {
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

	lastID, err := r.getUserLastID()
	user.ID = lastID + 1

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
