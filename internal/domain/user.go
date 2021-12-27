package domain

type User struct {
	Id          int64  `json:"id"`
	Firstname   string `json:"firstname" binding:"required"`
	Lastname    string `json:"lastname" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Age         int64  `json:"age" binding:"required"`
	Height      int64  `json:"height" binding:"required"`
	Active      bool   `json:"active"`
	CreatedDate string `json:"created_date"`
	DeletedDate string `json:"deleted_date"`
}

type UserTypes string

const (
	Id          UserTypes = "Id"
	Firstname   UserTypes = "Firstname"
	Lastname    UserTypes = "Lastname"
	Email       UserTypes = "Email"
	Age         UserTypes = "Age"
	Height      UserTypes = "Height"
	Active      UserTypes = "Active"
	CreatedDate UserTypes = "CreatedDate"
	DeletedDate UserTypes = "DeletedDate"
)
