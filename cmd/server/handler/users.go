package handler

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/miltonbernhardt/go-web/internal/model"
	"github.com/miltonbernhardt/go-web/internal/users"
	"github.com/miltonbernhardt/go-web/pkg/message"
	"github.com/miltonbernhardt/go-web/pkg/web"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
)

type UserController interface {
	GetAll() gin.HandlerFunc
	GetById() gin.HandlerFunc
	Store() gin.HandlerFunc
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
func (c *user) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Info("trying to get all users")
		allUsers, err := c.getUsersByFilters(ctx)

		if err != nil {
			log.Error(err)
			web.Error(ctx, http.StatusInternalServerError, message.InternalError)
			return
		}

		if len(allUsers) == 0 {
			allUsers = []model.User{}
		}

		log.WithFields(log.Fields{
			"users": len(allUsers),
		}).Info("success get all users")
		web.Success(ctx, http.StatusOK, allUsers)
	}
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
func (c *user) GetById() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, done := c.getIdFromParams(ctx)
		if done {
			return
		}

		user, err := c.service.GetByID(id)
		if c.checkError(ctx, err) {
			return
		}

		log.WithFields(log.Fields{
			"user": user,
		}).Info("success get user by ID")
		web.Success(ctx, http.StatusOK, user)
	}
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
func (c *user) Store() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		log.Info("trying to store user")
		var userEntity model.User

		if err := ctx.ShouldBindJSON(&userEntity); err != nil {
			web.ValidationError(ctx, http.StatusUnprocessableEntity, err)
			return
		}

		userEntity, err := c.service.Store(userEntity)

		if c.checkError(ctx, err) {
			return
		}

		log.WithFields(log.Fields{
			"user": userEntity,
		}).Info("success store user")
		web.Success(ctx, http.StatusCreated, userEntity)
	}
}

//Update 		godoc
//@Summary      Update   user
//@Tags         Users
//@Description  update user
//@Accept       json
//@Produce      json
//@Param        token  header    string       true  "token"
//@Param        user   body      model.User  true  "user"
//@Param        id     path      string       true  "id user"
//@Success      200    {object}  web.Response
//@Failure      400    {object}  web.ErrorResponse
//@Failure      401    {object}  web.ErrorResponse
//@Failure      404    {object}  web.ErrorResponse
//@Router       /users/{id} [put]
func (c *user) Update() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var userEntity model.User
		if err := ctx.ShouldBindJSON(&userEntity); err != nil {
			web.ValidationError(ctx, http.StatusUnprocessableEntity, err)
			return
		}

		id, done := c.getIdFromParams(ctx)
		if done {
			return
		}

		userEntity, err := c.service.Update(id, userEntity)
		if c.checkError(ctx, err) {
			return
		}

		log.WithFields(log.Fields{
			"user": userEntity,
		}).Info("success update user")
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

		err := c.service.Delete(id)
		if c.checkError(ctx, err) {
			return
		}

		log.WithFields(log.Fields{
			"user_id": id,
		}).Info("success delete user")
		web.Success(ctx, http.StatusOK, message.UserDeleted)
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
			log.Info("error: bad request for update fields")
			web.Error(ctx, http.StatusBadRequest, message.UserInvalidUpdate)
			return
		}

		user, err := c.service.UpdateFields(id, fields.Lastname, fields.Age)

		if c.checkError(ctx, err) {
			return
		}

		log.WithFields(log.Fields{
			"user_id": id,
		}).Info("success update fields user")
		web.Success(ctx, http.StatusOK, user)
	}
}

func (c *user) getIdFromParams(ctx *gin.Context) (int, bool) {
	id, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		log.Info("error: invalid ID")
		web.Error(ctx, http.StatusBadRequest, message.InvalidID)
		return 0, true
	}
	return id, false
}

func (c *user) ValidateToken(ctx *gin.Context) {
	//ToDo: add a const TOKEN
	if !(ctx.GetHeader("token") != "" && os.Getenv("TOKEN") != "" && ctx.GetHeader("token") == os.Getenv("TOKEN")) {
		web.Error(ctx, http.StatusUnauthorized, message.UnauthorizedAction)
		log.Info(http.StatusUnauthorized, "invalid token")
		ctx.Abort()
		return
	}

	ctx.Next()
}

func (c *user) getUsersByFilters(ctx *gin.Context) ([]model.User, error) {
	usersSlice, err := c.service.GetAll()

	if err != nil {
		log.Error(err)
		return nil, err
	}

	if firstname := ctx.Query("firstname"); firstname != "" {
		usersSlice = c.service.GetAllWithFilters(usersSlice, model.Firstname, firstname)
	}

	if lastname := ctx.Query("lastname"); lastname != "" {
		usersSlice = c.service.GetAllWithFilters(usersSlice, model.Lastname, lastname)
	}

	if email := ctx.Query("email"); email != "" {
		usersSlice = c.service.GetAllWithFilters(usersSlice, model.Email, email)
	}

	if createdDate := ctx.Query("created_date"); createdDate != "" {
		usersSlice = c.service.GetAllWithFilters(usersSlice, model.CreatedDate, createdDate)
	}

	if activeString := ctx.Query("active"); activeString != "" {
		if isActive, err := strconv.ParseBool(activeString); err == nil {
			usersSlice = c.service.GetAllWithFilters(usersSlice, model.Active, isActive)
		}
	}

	if ageString := ctx.Query("age"); ageString != "" {
		if age, err := strconv.Atoi(ageString); err == nil {
			usersSlice = c.service.GetAllWithFilters(usersSlice, model.Age, age)
		}
	}

	if heightString := ctx.Query("height"); heightString != "" {
		if height, err := strconv.Atoi(heightString); err == nil {
			usersSlice = c.service.GetAllWithFilters(usersSlice, model.Height, height)
		}
	}
	return usersSlice, nil
}

func (c *user) checkError(ctx *gin.Context, err error) bool {
	if err != nil {
		if err.Error() == message.UserNotFound {
			log.Info(err)
			web.Error(ctx, http.StatusNotFound, message.UserNotFound)
			return true
		} else {
			log.Error(err)
			web.Error(ctx, http.StatusInternalServerError, message.InternalError)
			return true
		}
	}
	return false
}
