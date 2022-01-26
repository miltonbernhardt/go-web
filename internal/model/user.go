package model

type User struct {
	ID          int    `json:"id"`
	Firstname   string `json:"firstname" binding:"required"`
	Lastname    string `json:"lastname" binding:"required"`
	Email       string `json:"email" binding:"required"`
	Age         int    `json:"age" binding:"required,min=18"`
	Height      int    `json:"height" binding:"required"`
	Active      bool   `json:"active"`
	CreatedDate string `json:"created_date"`
}

type UserTypes string

// ToDo: refactor these const
const (
	Id          UserTypes = "ID"
	Firstname   UserTypes = "Firstname"
	Lastname    UserTypes = "Lastname"
	Email       UserTypes = "Email"
	Age         UserTypes = "Age"
	Height      UserTypes = "Height"
	Active      UserTypes = "Active"
	CreatedDate UserTypes = "CreatedDate"
)
