package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

func Hello(ctx *gin.Context) {
	ctx.String(200, "Hola Milton!")
}

func GetUsers(ctx *gin.Context) {
	users, err := GetUsersAsSlice()
	if err != nil {
		ctx.String(500, fmt.Sprint(err))
	}

	firstname := ctx.Query("firstname")
	fmt.Printf("\nfirstname: %v\n", firstname)
	if firstname != "" {
		users = GetUsersByFirstName(users, firstname)
	}

	lastname := ctx.Query("lastname")
	fmt.Printf("\nlastname: %v\n\n", lastname)
	if lastname != "" {
		users = GetUsersByLastname(users, lastname)
	}

	email := ctx.Query("email")
	fmt.Printf("\nemail: %v\n\n", email)
	if lastname != "" {
		users = GetUsersByEmail(users, email)
	}

	createdDate := ctx.Query("created_date")
	fmt.Printf("\ncreated_date: %v\n\n", createdDate)
	if createdDate != "" {
		users = GetUsersByCreatedDate(users, createdDate)
	}

	if activeString := ctx.Query("is_active"); activeString != "" {
		isActive, err := strconv.ParseBool(activeString)
		fmt.Printf("\nactive: %v - err: %v\n\n", isActive, err)
		if err == nil {
			users = GetUsersByIsActive(users, isActive)
		}
	}

	if ageString := ctx.Query("age"); ageString != "" {
		age, err := strconv.ParseInt(ageString, 10, 64)
		fmt.Printf("\nage: %v - err: %v\n\n", age, err)
		if err == nil {
			users = GetUsersByAge(users, age)
		}
	}

	if heightString := ctx.Query("height"); heightString != "" {
		height, err := strconv.ParseInt(heightString, 10, 64)
		fmt.Printf("\nheight: %v - err: %v\n\n", height, err)
		if err == nil {
			users = GetUsersByHeight(users, height)
		}
	}

	if len(users) == 0 {
		ctx.String(404, "error: no se encontraron usuarios que coincidan con la búsqueda")
		ctx.Abort()
		return
	}

	ctx.JSON(200, users)

}

func GetUserById(ctx *gin.Context) {
	//todo ver de tratar de obtener el user desde el map de json
	users, err := GetUsersAsSlice()
	if err != nil {
		ctx.String(500, fmt.Sprint(err))
		ctx.Abort()
		return
	}

	id := ctx.Param("id")

	user, err := GetUserByID(users, id)
	if err != nil {
		ctx.String(404, fmt.Sprint(err))
		ctx.Abort()
		return
	}

	ctx.JSON(200, user)
}
