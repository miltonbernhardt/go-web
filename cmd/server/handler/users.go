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

type UserController interface {
	GetAll(ctx *gin.Context)
	GetById(ctx *gin.Context)
	Store(ctx *gin.Context)
	Update() gin.HandlerFunc
	Delete() gin.HandlerFunc
	UpdateFields() gin.HandlerFunc
	ValidateToken(ctx *gin.Context)
}

type user struct {
	service users.Service
}

func NewUserController(s users.Service) UserController {
	return &user{
		service: s,
	}
}

//ListUsers godoc
//@Summary      List users
//@Tags         Users
//@Description  get users
//@Accept       json
//@Produce      json
//@Param        token  header    string  true  "token"
//@Success      200    {object}  web.Response
//@Router       /users [get]
func (c *user) GetAll(ctx *gin.Context) {
	allUsers, err := c.getUsersByQuery(ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, web.NewResponse(http.StatusInternalServerError, fmt.Sprint(err)))
		return
	}

	if len(allUsers) == 0 {
		ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "no se encontraron usuarios que coincidan con la búsqueda"))
		return
	}

	ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, allUsers))
}

//GetById 		godoc
//@Summary      GetById user
//@Tags         Users
//@Description  get user by id
//@Accept       json
//@Produce      json
//@Param        token  header    string  true  "token"
//@Param        id     path      string  true  "id user"
//@Success      201    {object}  web.Response
//@Failure      400    {object}  web.Error
//@Failure      401    {object}  web.Error
//@Failure      404    {object}  web.Error
//@Router       /users/{id} [get]
func (c *user) GetById(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "invalid ID"))
	}

	user, err := c.service.FetchUserByID(id)

	if err != nil {
		ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, err.Error()))
		return
	}

	ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, user))
}

//StoreUser 		godoc
//@Summary      Store user
//@Tags         Users
//@Description  store user
//@Accept       json
//@Produce      json
//@Param        token  header    string  true  "token"
//@Success      201    {object}  web.Response
//@Failure      400    {object}  web.Error
//@Failure      401    {object}  web.Error
//@Failure      500    {object}  web.Error
//@Router       /users [post]
func (c *user) Store(ctx *gin.Context) {
	var newUser domain.User

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		if errorsToPrint := web.Simple(err); len(errorsToPrint) > 0 {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, errorsToPrint))
		} else {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, err.Error()))
		}
		return
	}

	newUser, err := c.service.StoreUser(newUser)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, web.NewResponse(http.StatusInternalServerError, err.Error()))
		return
	}

	ctx.JSON(http.StatusCreated, web.NewResponse(http.StatusOK, newUser))
}

//UpdateUser 	godoc
//@Summary      Update   user
//@Tags         Users
//@Description  update user
//@Accept       json
//@Produce      json
//@Param        token  header    string       true  "token"
//@Param        user   body      domain.User  true  "user"
//@Param        id     path      string       true  "id user"
//@Success      200    {object}  web.Response
//@Failure      400    {object}  web.Error
//@Failure      401    {object}  web.Error
//@Failure      404    {object}  web.Error
//@Router       /users/{id} [put]
func (c *user) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var updatedUser domain.User
		if err := ctx.ShouldBindJSON(&updatedUser); err != nil {
			validatorErrors := err.(validator.ValidationErrors)

			errorsToPrint := map[string]string{}
			for _, fieldError := range validatorErrors {
				errorsToPrint[fieldError.Field()] = fmt.Sprintf("el campo %v es requerido", fieldError.Field())
			}

			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, errorsToPrint))
			return
		}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, "invalid ID"))
		}

		updatedUser, err = c.service.UpdateUser(id, updatedUser)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, updatedUser))
	}
}

//DeleteUser 	godoc
//@Summary      Delete   user
//@Tags         Users
//@Description  unregister a user
//@Accept       json
//@Produce      json
//@Param        token  header    string  true  "token"
//@Param        id     path      string  true  "id user"
//@Success      200    {object}  web.Response
//@Failure      400    {object}  web.Error
//@Failure      401    {object}  web.Error
//@Failure      404    {object}  web.Error
//@Router       /users/{id} [delete]
func (c *user) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "invalid ID"))
		}

		err = c.service.DeleteUser(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, fmt.Sprintf("El producto %d ha sido eliminado", id)))
	}
}

//UpdateFields 	godoc
//@Summary      UpdateFields user
//@Tags         Users
//@Description  update Lastname and/or age from user
//@Accept       json
//@Produce      json
//@Param        token     header    string  true  "token"
//@Param        lastname  body      int     true  "lastname"
//@Param        age       body      string  true  "age"
//@Param        id        path      string  true  "id user"
//@Success      201       {object}  web.Response
//@Failure      400       {object}  web.Error
//@Failure      401       {object}  web.Error
//@Failure      404       {object}  web.Error
//@Router       /users/{id} [patch]
func (c *user) UpdateFields() gin.HandlerFunc {
	type userFields struct {
		Lastname string `json:"lastname"`
		Age      int64  `json:"age"`
	}

	return func(ctx *gin.Context) {
		fields := userFields{}

		id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, "invalid ID"))
			return
		}

		bodyAsByteArray, _ := ioutil.ReadAll(ctx.Request.Body)
		err = json.Unmarshal(bodyAsByteArray, &fields)
		if err != nil || (fields.Lastname == "" && fields.Age == 0) {
			ctx.JSON(http.StatusBadRequest, web.NewResponse(http.StatusBadRequest, fmt.Sprintf("modificación invalida del usuario %d", id)))
			return
		}

		user, err := c.service.UpdateFieldsUser(id, fields.Lastname, fields.Age)

		if err != nil {
			ctx.JSON(http.StatusNotFound, web.NewResponse(http.StatusNotFound, err.Error()))
			return
		}

		ctx.JSON(http.StatusOK, web.NewResponse(http.StatusOK, user))
	}
}

func (c *user) ValidateToken(ctx *gin.Context) {
	if !(ctx.GetHeader("token") != "" && os.Getenv("TOKEN") != "" && ctx.GetHeader("token") == os.Getenv("TOKEN")) {
		ctx.JSON(http.StatusUnauthorized, web.ResponseUnauthorized())
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (c *user) getUsersByQuery(ctx *gin.Context) ([]domain.User, error) {
	usersSlice, err := c.service.FetchAllUsers()

	if err != nil {
		return nil, err
	}

	if firstname := ctx.Query("firstname"); firstname != "" {
		usersSlice = c.service.FetchAllUsersByQuery(usersSlice, domain.Firstname, firstname)
	}

	if lastname := ctx.Query("lastname"); lastname != "" {
		usersSlice = c.service.FetchAllUsersByQuery(usersSlice, domain.Lastname, lastname)
	}

	if email := ctx.Query("email"); email != "" {
		usersSlice = c.service.FetchAllUsersByQuery(usersSlice, domain.Email, email)
	}

	if createdDate := ctx.Query("created_date"); createdDate != "" {
		usersSlice = c.service.FetchAllUsersByQuery(usersSlice, domain.CreatedDate, createdDate)
	}

	if activeString := ctx.Query("active"); activeString != "" {
		if isActive, err := strconv.ParseBool(activeString); err == nil {
			usersSlice = c.service.FetchAllUsersByQuery(usersSlice, domain.Active, isActive)
		}
	}

	if ageString := ctx.Query("age"); ageString != "" {
		if age, err := strconv.ParseInt(ageString, 10, 64); err == nil {
			usersSlice = c.service.FetchAllUsersByQuery(usersSlice, domain.Age, age)
		}
	}

	if heightString := ctx.Query("height"); heightString != "" {
		if height, err := strconv.ParseInt(heightString, 10, 64); err == nil {
			usersSlice = c.service.FetchAllUsersByQuery(usersSlice, domain.Height, height)
		}
	}
	return usersSlice, nil
}
