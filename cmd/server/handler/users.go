package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/miltonbernhardt/go-web/internal/domain"
	"github.com/miltonbernhardt/go-web/internal/users"
	"github.com/miltonbernhardt/go-web/pkg/web"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type User interface {
	GetAll(ctx *gin.Context)
	GetById(ctx *gin.Context)
	Save(ctx *gin.Context)
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc
	UpdateFields() gin.HandlerFunc
}

type user struct {
	service users.Service
}

func NewUserController(s users.Service) User {
	return &user{
		service: s,
	}
}

func (c *user) GetAll(ctx *gin.Context) {
	allUsers, err := c.getUsersByQuery(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, web.NewResponse(http.StatusInternalServerError, fmt.Sprint(err)))
		ctx.Abort()
		return
	}

	if len(allUsers) == 0 {
		ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "no se encontraron usuarios que coincidan con la búsqueda"))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, allUsers))

}

func (c *user) getUsersByQuery(ctx *gin.Context) ([]domain.User, error) {
	usersSlice, err := c.service.GetAll()

	if err != nil {
		return nil, err
	}

	if firstname := ctx.Query("firstname"); firstname != "" {
		usersSlice = c.service.GetAllByField(usersSlice, domain.Firstname, firstname)
	}

	if lastname := ctx.Query("lastname"); lastname != "" {
		usersSlice = c.service.GetAllByField(usersSlice, domain.Lastname, lastname)
	}

	if email := ctx.Query("email"); email != "" {
		usersSlice = c.service.GetAllByField(usersSlice, domain.Email, email)
	}

	if createdDate := ctx.Query("created_date"); createdDate != "" {
		usersSlice = c.service.GetAllByField(usersSlice, domain.CreatedDate, createdDate)
	}

	if activeString := ctx.Query("active"); activeString != "" {
		if isActive, err := strconv.ParseBool(activeString); err == nil {
			usersSlice = c.service.GetAllByField(usersSlice, domain.Active, isActive)
		}
	}

	if ageString := ctx.Query("age"); ageString != "" {
		if age, err := strconv.ParseInt(ageString, 10, 64); err == nil {
			usersSlice = c.service.GetAllByField(usersSlice, domain.Age, age)
		}
	}

	if heightString := ctx.Query("height"); heightString != "" {
		if height, err := strconv.ParseInt(heightString, 10, 64); err == nil {
			usersSlice = c.service.GetAllByField(usersSlice, domain.Height, height)
		}
	}
	return usersSlice, nil
}

func (c *user) GetById(ctx *gin.Context) {

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "invalid ID"))
	}

	user, err := c.service.GetByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, err.Error()))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, user))
}

func (c *user) Save(ctx *gin.Context) {
	if !(ctx.GetHeader("token") != "" && os.Getenv("TOKEN") != "" && ctx.GetHeader("token") == os.Getenv("TOKEN")) {
		ctx.JSON(http.StatusUnauthorized, web.ResponseUnauthorized())
		ctx.Abort()
		return
	}

	var newUser domain.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {

		validatorErrors := err.(validator.ValidationErrors)

		errorsToPrint := map[string]string{}
		for _, fieldError := range validatorErrors {
			errorsToPrint[fieldError.Field()] = fmt.Sprintf("el campo %v es requerido", fieldError.Field())
		}

		ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, errorsToPrint))
		ctx.Abort()
		return
	}

	newUser, err := c.service.Store(newUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, web.NewResponse(http.StatusInternalServerError, err.Error()))
		ctx.Abort()
		return
	}

	ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, newUser))
}

func (c *user) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !(ctx.GetHeader("token") != "" && os.Getenv("TOKEN") != "" && ctx.GetHeader("token") == os.Getenv("TOKEN")) {
			ctx.JSON(http.StatusUnauthorized, web.ResponseUnauthorized())
			ctx.Abort()
			return
		}

		var updatedUser domain.User
		if err := ctx.ShouldBindJSON(&updatedUser); err != nil {
			validatorErrors := err.(validator.ValidationErrors)

			errorsToPrint := map[string]string{}
			for _, fieldError := range validatorErrors {
				errorsToPrint[fieldError.Field()] = fmt.Sprintf("el campo %v es requerido", fieldError.Field())
			}

			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, errorsToPrint))
			ctx.Abort()
			return
		}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "invalid ID"))
		}

		updatedUser, err = c.service.Update(id, updatedUser)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, err.Error()))
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, updatedUser))
	}
}

func (c *user) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if !(ctx.GetHeader("token") != "" && os.Getenv("TOKEN") != "" && ctx.GetHeader("token") == os.Getenv("TOKEN")) {
			ctx.JSON(http.StatusUnauthorized, web.ResponseUnauthorized())
			ctx.Abort()
			return
		}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "invalid ID"))
		}

		err = c.service.Delete(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, err.Error()))
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, fmt.Sprintf("El producto %d ha sido eliminado", id)))
	}
}

func (c *user) UpdateFields() gin.HandlerFunc {
	type userFields struct {
		Lastname string `json:"lastname"`
		Age      int64  `json:"age"`
	}

	return func(ctx *gin.Context) {
		if !(ctx.GetHeader("token") != "" && os.Getenv("TOKEN") != "" && ctx.GetHeader("token") == os.Getenv("TOKEN")) {
			ctx.JSON(http.StatusUnauthorized, web.ResponseUnauthorized())
			ctx.Abort()
			return
		}

		fields := userFields{}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "invalid ID"))
		}

		bodyAsByteArray, _ := ioutil.ReadAll(ctx.Request.Body)
		err = json.Unmarshal(bodyAsByteArray, &fields)
		if err != nil || (fields.Lastname == "" && fields.Age == 0) {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, fmt.Sprintf("modificación invalida del usuario %d", id)))
			ctx.Abort()
			return
		}

		user, err := c.service.UpdateFields(id, fields.Lastname, fields.Age)

		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, err.Error()))
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, user))
	}
}
