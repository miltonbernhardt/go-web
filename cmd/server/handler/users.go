package handler

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
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

//GetAll 		godoc
//@Summary      List users
//@Tags         Users
//@Description  get users
//@Accept       json
//@Produce      json
//@Param        token  header    string  true  "token"
//@Success      200    {object}  web.Response
//@Router       /users [get]
func (c *user) GetAll(ctx *gin.Context) {
	allUsers, err := c.getUsersByFilters(ctx)

	if err != nil {
		fmt.Println(err)
		web.Error(ctx, http.StatusInternalServerError, web.InternalError)
		return
	}

	if len(allUsers) == 0 {
		allUsers = []domain.User{}
	}

	web.Success(ctx, http.StatusOK, allUsers)
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
//@Failure      400    {object}  web.ErrorResponse
//@Failure      401    {object}  web.ErrorResponse
//@Failure      404    {object}  web.ErrorResponse
//@Router       /users/{id} [get]
func (c *user) GetById(ctx *gin.Context) {
	id, done := c.getIdFromParams(ctx)
	if done {
		return
	}

	user, err := c.service.FetchUserByID(id)
	if c.checkError(ctx, err) {
		return
	}

	web.Success(ctx, http.StatusOK, user)
}

//Store 		godoc
//@Summary      Store user
//@Tags         Users
//@Description  store user
//@Accept       json
//@Produce      json
//@Param        token  header    string  true  "token"
//@Success      201    {object}  web.Response
//@Failure      400    {object}  web.ErrorResponse
//@Failure      401    {object}  web.ErrorResponse
//@Failure      500    {object}  web.ErrorResponse
//@Router       /users [post]
func (c *user) Store(ctx *gin.Context) {
	var userEntity domain.User

	if err := ctx.ShouldBindJSON(&userEntity); err != nil {
		web.ValidationError(ctx, http.StatusUnprocessableEntity, err)
		return
	}

	userEntity, err := c.service.StoreUser(userEntity)

	if c.checkError(ctx, err) {
		return
	}

	web.Success(ctx, http.StatusCreated, userEntity)
}

//Update 		godoc
//@Summary      Update   user
//@Tags         Users
//@Description  update user
//@Accept       json
//@Produce      json
//@Param        token  header    string       true  "token"
//@Param        user   body      domain.User  true  "user"
//@Param        id     path      string       true  "id user"
//@Success      200    {object}  web.Response
//@Failure      400    {object}  web.ErrorResponse
//@Failure      401    {object}  web.ErrorResponse
//@Failure      404    {object}  web.ErrorResponse
//@Router       /users/{id} [put]
func (c *user) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userEntity domain.User
		if err := ctx.ShouldBindJSON(&userEntity); err != nil {
			web.ValidationError(ctx, http.StatusUnprocessableEntity, err)
			return
		}

		id, done := c.getIdFromParams(ctx)
		if done {
			return
		}

		userEntity, err := c.service.UpdateUser(id, userEntity)
		if c.checkError(ctx, err) {
			return
		}

		web.Success(ctx, http.StatusOK, userEntity)
	}
}

//Delete 		godoc
//@Summary      Delete   user
//@Tags         Users
//@Description  unregisters a user
//@Accept       json
//@Produce      json
//@Param        token  header    string  true  "token"
//@Param        id     path      string  true  "id user"
//@Success      200    {object}  web.Response
//@Failure      400    {object}  web.ErrorResponse
//@Failure      401    {object}  web.ErrorResponse
//@Failure      404    {object}  web.ErrorResponse
//@Router       /users/{id} [delete]
func (c *user) Delete() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, done := c.getIdFromParams(ctx)
		if done {
			return
		}

		err := c.service.DeleteUser(id)
		if c.checkError(ctx, err) {
			return
		}

		web.Success(ctx, http.StatusOK, web.UserDeleted)
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
//@Failure      400       {object}  web.ErrorResponse
//@Failure      401       {object}  web.ErrorResponse
//@Failure      404       {object}  web.ErrorResponse
//@Router       /users/{id} [patch]
func (c *user) UpdateFields() gin.HandlerFunc {
	type userFields struct {
		Lastname string `json:"lastname"`
		Age      int    `json:"age"`
	}

	return func(ctx *gin.Context) {
		fields := userFields{}

		id, done := c.getIdFromParams(ctx)
		if done {
			return
		}

		bodyAsByteArray, _ := ioutil.ReadAll(ctx.Request.Body)
		err := json.Unmarshal(bodyAsByteArray, &fields)
		if err != nil || (fields.Lastname == "" && fields.Age == 0) {
			web.Error(ctx, http.StatusBadRequest, web.UserInvalidUpdate)
			return
		}

		user, err := c.service.UpdateFieldsUser(id, fields.Lastname, fields.Age)

		if c.checkError(ctx, err) {
			return
		}

		web.Success(ctx, http.StatusOK, user)
	}
}

func (c *user) getIdFromParams(ctx *gin.Context) (int, bool) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		web.Success(ctx, http.StatusBadRequest, web.InvalidID)
		return 0, true
	}
	return id, false
}

func (c *user) ValidateToken(ctx *gin.Context) {
	if !(ctx.GetHeader("token") != "" && os.Getenv("TOKEN") != "" && ctx.GetHeader("token") == os.Getenv("TOKEN")) {
		web.Error(ctx, http.StatusUnauthorized, web.UnauthorizedAction)
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (c *user) getUsersByFilters(ctx *gin.Context) ([]domain.User, error) {
	usersSlice, err := c.service.FetchAllUsers()
	fmt.Printf("\n\n\t%v\t\n\n", usersSlice)
	fmt.Printf("\n\n\t%v\t\n\n", err)

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

func (c *user) checkError(ctx *gin.Context, err error) bool {
	if err != nil {
		if err.Error() == web.UserNotFound {
			web.Error(ctx, http.StatusNotFound, web.UserNotFound)
			return true
		} else {
			fmt.Println(err)
			web.Error(ctx, http.StatusInternalServerError, web.InternalError)
			return true
		}
	}
	return false
}
