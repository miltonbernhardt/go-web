package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/miltonbernhardt/go-web/internal/users"
	"log"
	"net/http"
	"strconv"
)

type User struct {
	service users.Service
}

func NewUserController(u users.Service) *User {
	return &User{
		service: u,
	}
}

func (c *User) GetUsers(ctx *gin.Context) {
	usersSlice, err := c.service.GetAll()

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprint(err)})
		ctx.Abort()
		return
	}
	//todo generalizar
	//firstname := ctx.Query("firstname")
	//fmt.Printf("\nfirstname: %v\n", firstname)
	//if firstname != "" {
	//	usersSlice = c.service.GetByFirstName(usersSlice, firstname)
	//}
	//
	//lastname := ctx.Query("lastname")
	//fmt.Printf("\nlastname: %v\n\n", lastname)
	//if lastname != "" {
	//	usersSlice = c.service.GetByLastname(usersSlice, lastname)
	//}
	//
	//email := ctx.Query("email")
	//fmt.Printf("\nemail: %v\n\n", email)
	//if lastname != "" {
	//	usersSlice = c.service.GetByEmail(usersSlice, email)
	//}
	//
	//createdDate := ctx.Query("created_date")
	//fmt.Printf("\ncreated_date: %v\n\n", createdDate)
	//if createdDate != "" {
	//	usersSlice = c.service.GetByCreatedDate(usersSlice, createdDate)
	//}
	//
	//if activeString := ctx.Query("is_active"); activeString != "" {
	//	isActive, err := strconv.ParseBool(activeString)
	//	fmt.Printf("\nactive: %v - err: %v\n\n", isActive, err)
	//	if err == nil {
	//		usersSlice = c.service.GetByIsActive(usersSlice, isActive)
	//	}
	//}
	//
	//if ageString := ctx.Query("age"); ageString != "" {
	//	age, err := strconv.ParseInt(ageString, 10, 64)
	//	fmt.Printf("\nage: %v - err: %v\n\n", age, err)
	//	if err == nil {
	//		usersSlice = c.service.GetByAge(usersSlice, age)
	//	}
	//}
	//
	//if heightString := ctx.Query("height"); heightString != "" {
	//	height, err := strconv.ParseInt(heightString, 10, 64)
	//	fmt.Printf("\nheight: %v - err: %v\n\n", height, err)
	//	if err == nil {
	//		usersSlice = c.service.GetByHeight(usersSlice, height)
	//	}
	//}

	if len(usersSlice) == 0 {
		ctx.String(http.StatusNotFound, "error: no se encontraron usuarios que coincidan con la búsqueda")
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, usersSlice)

}

func (c *User) GetUserById(ctx *gin.Context) {
	if idString := ctx.Param("id"); idString != "" {
		if id, err := strconv.ParseInt(idString, 10, 64); err == nil {
			user, err := c.service.GetByID(id)

			if err != nil {
				ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err)})
				ctx.Abort()
				return
			}
			ctx.JSON(http.StatusOK, user)
			ctx.Abort()
			return
		} else {
			ctx.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprint(err)})
		}
	}
}

func (c *User) SaveUser(ctx *gin.Context) {
	if !(ctx.GetHeader("token") != "" && ctx.GetHeader("token") == users.TokenAuth) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "no tiene permisos para realizar la petición solicitada"})
		ctx.Abort()
		return
	}

	var newUser users.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {

		//field, erro := main.GetField(&err, "Key")
		//log.Printf("\t1 - valor: %v - error: %v\n", field, erro)
		//
		//field, erro = main.GetField(err, "Key")
		//log.Printf("\t2 - valor: %v - error: %v\n", field, erro)

		//field, erro =GetField(err, "Key")
		//log.Printf("\t3 - valor: %v - error: %v\n", field, erro)

		//v := reflect.ValueOf(err)
		//tipoObtenidoDeReflection := v.Type()
		//
		//fmt.Printf("%v\n", tipoObtenidoDeReflection) //validator.ValidationErrors
		//var h = validator.ValidationErrors{}
		//
		//fmt.Printf("%v\n", v.NumField()) //validator.ValidationErrors
		//
		//
		//for i := 0; i < len(validator.ValidationErrors(v)); i++ {
		//	fmt.Printf("El campo %s tiene como valor: %v\n", tipoObtenidoDeReflection.Field(i).Tag, v.Field(i).Interface())
		//}

		log.Printf("error: %v\n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	newUser, err := c.service.Store(newUser)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, newUser)
}
