package handler

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/miltonbernhardt/go-web/internal/users"
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
	usersSlice, err := c.getUsersByQuery(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprint(err)})
		ctx.Abort()
		return
	}

	if len(usersSlice) == 0 {
		ctx.String(http.StatusNotFound, "error: no se encontraron usuarios que coincidan con la búsqueda")
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, usersSlice)

}

func (c *User) getUsersByQuery(ctx *gin.Context) ([]users.User, error) {
	usersSlice, err := c.service.GetAll()

	if err != nil {
		return nil, err
	}

	if firstname := ctx.Query("firstname"); firstname != "" {
		usersSlice = c.service.GetAllByField(usersSlice, users.Firstname, firstname)
	}

	if lastname := ctx.Query("lastname"); lastname != "" {
		usersSlice = c.service.GetAllByField(usersSlice, users.Lastname, lastname)
	}

	if email := ctx.Query("email"); email != "" {
		usersSlice = c.service.GetAllByField(usersSlice, users.Email, email)
	}

	if createdDate := ctx.Query("created_date"); createdDate != "" {
		usersSlice = c.service.GetAllByField(usersSlice, users.CreatedDate, createdDate)
	}

	if activeString := ctx.Query("active"); activeString != "" {
		if isActive, err := strconv.ParseBool(activeString); err == nil {
			usersSlice = c.service.GetAllByField(usersSlice, users.Active, isActive)
		}
	}

	if ageString := ctx.Query("age"); ageString != "" {
		if age, err := strconv.ParseInt(ageString, 10, 64); err == nil {
			usersSlice = c.service.GetAllByField(usersSlice, users.Age, age)
		}
	}

	if heightString := ctx.Query("height"); heightString != "" {
		if height, err := strconv.ParseInt(heightString, 10, 64); err == nil {
			usersSlice = c.service.GetAllByField(usersSlice, users.Height, height)
		}
	}
	return usersSlice, nil
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

		validatorErrors := err.(validator.ValidationErrors)

		errorsToPrint := map[string]string{}
		for _, fieldError := range validatorErrors {
			errorsToPrint[fieldError.Field()] = fmt.Sprintf("el campo %v es requerido", fieldError.Field())
		}

		ctx.JSON(http.StatusBadRequest, gin.H{"error": errorsToPrint})
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
