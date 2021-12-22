package main

type User struct {
	Id          string `json:"id"`
	Firstname   string `json:"firstname"`
	Lastname    string `json:"lastname"`
	Email       string `json:"email"`
	Age         int64  `json:"age"`
	Height      int64  `json:"height"`
	Active      bool   `json:"active"`
	CreatedDate string `json:"created_date"`
}
