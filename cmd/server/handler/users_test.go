package handler

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/miltonbernhardt/go-web/internal/model"
	"github.com/miltonbernhardt/go-web/internal/users"
	"github.com/miltonbernhardt/go-web/internal/utils"
	"github.com/miltonbernhardt/go-web/pkg/message"
	"github.com/miltonbernhardt/go-web/pkg/store"
	"github.com/miltonbernhardt/go-web/pkg/web"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var (
	token   = "bearer 123456"
	urlPath = "/users"
)

type responseError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type responseOfUser struct {
	Data model.User `json:"data"`
}

type responseSliceOfUsers struct {
	Data []model.User `json:"data"`
}

type validationError struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

type responseValidationsErrors struct {
	Code    string            `json:"code"`
	Message string            `json:"message"`
	Fields  []validationError `json:"fields"`
}

func createServer() *gin.Engine {
	expectedUsers := []model.User{
		{
			ID:          1,
			Firstname:   "firstname",
			Lastname:    "lastname",
			Email:       "email",
			Age:         24,
			Height:      184,
			Active:      true,
			CreatedDate: "22/02/2021",
		},
		{
			ID:          2,
			Firstname:   "nombre",
			Lastname:    "apellido",
			Email:       "mail",
			Age:         25,
			Height:      185,
			Active:      false,
			CreatedDate: "23/03/2021",
		},
		{
			ID:          3,
			Firstname:   "firstname3",
			Lastname:    "lastname3",
			Email:       "email3",
			Age:         26,
			Height:      187,
			Active:      false,
			CreatedDate: "25/02/2021",
		},
	}

	dataJson, _ := json.Marshal(expectedUsers)

	dbStub := store.Mock{
		Data: dataJson,
		Err:  nil,
	}
	storeMocked := store.FileStore{
		FileName: "",
		Mock:     &dbStub,
	}

	_ = os.Setenv("TOKEN", token)
	repo := users.NewRepositoryFile(&storeMocked)
	util := utils.New()
	util.AddMock(&utils.Mock{Date: "02/01/2006 15:04:05"})
	service := users.NewService(repo, util)
	u := NewUserController(service)

	return initializeGin(u)
}

func createServerFailDB() *gin.Engine {
	_ = os.Setenv("TOKEN", token)
	db := store.New("error", "../products.json")
	repo := users.NewRepositoryFile(db)
	util := utils.New()
	util.AddMock(&utils.Mock{Date: "02/01/2006 15:04:05"})
	service := users.NewService(repo, util)
	u := NewUserController(service)

	return initializeGin(u)
}

func initializeGin(u UserController) *gin.Engine {
	log.SetLevel(log.FatalLevel)
	gin.DefaultWriter = ioutil.Discard //turn off gin logging
	r := gin.Default()

	ur := r.Group(urlPath)
	ur.GET("/", u.ValidateToken, u.GetAll())
	ur.GET("/:id", u.ValidateToken, u.GetById())
	ur.POST("/", u.ValidateToken, u.Store())
	ur.PUT("/:id", u.ValidateToken, u.Update())
	ur.DELETE("/:id", u.ValidateToken, u.Delete())
	ur.PATCH("/:id", u.ValidateToken, u.UpdateFields())
	return r
}

func createRequestTest(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, urlPath+url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("token", token)

	return req, httptest.NewRecorder()
}

func createRequestTestWithoutToken(method string, url string, body string) (*http.Request, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, urlPath+url, bytes.NewBuffer([]byte(body)))
	req.Header.Add("Content-Type", "application/json")

	return req, httptest.NewRecorder()
}

/*************************** GET ALL ***************************/

func Test_GetAll_AuthorizationTokenMissing(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTestWithoutToken(http.MethodGet, "/", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnauthorized), objReq.Code)
	assert.Equal(t, message.UnauthorizedAction, objReq.Message)
}

func Test_GetAll_InternalError(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTest(http.MethodGet, "/", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusInternalServerError), objReq.Code)
	assert.Equal(t, message.InternalError, objReq.Message)
}

func Test_GetAll_OK(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, true, len(objReq.Data) > 0)
}

func Test_GetAll_OK_FilterByFirstname(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?firstname=firstname", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(objReq.Data))
}

func Test_GetAll_OK_FilterByLastname(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?lastname=apellido", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

func Test_GetAll_OK_FilterByEmail(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?email=email3", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

func Test_GetAll_OK_FilterByCreatedDate(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?created_date=25/02/2021", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

func Test_GetAll_OK_FilterByActive(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?active=true", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

func Test_GetAll_OK_FilterByAge(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?age=25", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

func Test_GetAll_OK_FilterByHeight(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?height=184", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

func Test_GetAll_OK_NoneUser(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?age=80", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(objReq.Data))
}

func Test_GetAll_OK_TwoFilters(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?firstname=firstname&lastname=lastname", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 2, len(objReq.Data))
}

func Test_GetAll_OK_ThreeFilters(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?firstname=firstname&lastname=lastname&age=24", "")

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

func Test_GetAll_OK_FourFilters(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?lastname=lastname&age=24&height=184&created_date=22/02/2021", "")

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

func Test_GetAll_OK_FiveFilters(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?email=email&age=24&height=184&active=true&created_date=22/02/2021", "")

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

func Test_GetAll_OK_SixFilters(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?firstname=firstname&lastname=lastname&email=email&age=24&height=184&active=true", "")

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

func Test_GetAll_OK_SevenFilters(t *testing.T) {
	objReq := responseSliceOfUsers{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/?firstname=firstname&lastname=lastname&email=email&age=24&height=184&active=true&created_date=22/02/2021", "")

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(objReq.Data))
}

/*************************** GET BY ID ***************************/

func Test_GetByID_AuthorizationTokenMissing(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTestWithoutToken(http.MethodGet, "/1", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnauthorized), objReq.Code)
	assert.Equal(t, message.UnauthorizedAction, objReq.Message)
}

func Test_GetByID_InternalError(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTest(http.MethodGet, "/1", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusInternalServerError), objReq.Code)
	assert.Equal(t, message.InternalError, objReq.Message)
}

func Test_GetByID_NotFound(t *testing.T) {
	objReq := responseError{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/100", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusNotFound), objReq.Code)
	assert.Equal(t, message.UserNotFound, objReq.Message)
}

func Test_GetByID_InvalidID(t *testing.T) {
	objReq := responseError{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/ID", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusBadRequest), objReq.Code)
	assert.Equal(t, message.InvalidID, objReq.Message)
}

func Test_GetByID_OK(t *testing.T) {
	objReq := responseOfUser{}
	r := createServer()

	req, rr := createRequestTest(http.MethodGet, "/1", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, objReq.Data.ID)
	assert.Equal(t, "firstname", objReq.Data.Firstname)
	assert.Equal(t, "lastname", objReq.Data.Lastname)
	assert.Equal(t, "email", objReq.Data.Email)
	assert.Equal(t, 24, objReq.Data.Age)
	assert.Equal(t, 184, objReq.Data.Height)
	assert.Equal(t, true, objReq.Data.Active)
	assert.Equal(t, "22/02/2021", objReq.Data.CreatedDate)
}

/*************************** STORE ***************************/

func Test_Store_AuthorizationTokenMissing(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTestWithoutToken(http.MethodPost, "/", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"age": 30,
		"height": 170,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnauthorized), objReq.Code)
	assert.Equal(t, message.UnauthorizedAction, objReq.Message)
}

func Test_Store_InternalError(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTest(http.MethodPost, "/", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"age": 30,
		"height": 170,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusInternalServerError), objReq.Code)
	assert.Equal(t, message.InternalError, objReq.Message)
}

func Test_Store_UnprocessableEntity_FirstnameMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPost, "/", `{
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"height": 170,
		"age": 40,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 1, len(objReq.Fields))
	assert.Equal(t, "Firstname", objReq.Fields[0].Field)
}

func Test_Store_UnprocessableEntity_LastnameMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPost, "/", `{
        "firstname": "Testing",
		"email": "testing@testing.org",
		"height": 170,
		"age": 40,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 1, len(objReq.Fields))
	assert.Equal(t, "Lastname", objReq.Fields[0].Field)
}

func Test_Store_UnprocessableEntity_EmailMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPost, "/", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"height": 170,
		"age": 40,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 1, len(objReq.Fields))
	assert.Equal(t, "Email", objReq.Fields[0].Field)
}

func Test_Store_UnprocessableEntity_AgeMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPost, "/", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"height": 170,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 1, len(objReq.Fields))
	assert.Equal(t, "Age", objReq.Fields[0].Field)
}

func Test_Store_UnprocessableEntity_HeightMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPost, "/", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"age": 40,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 1, len(objReq.Fields))
	assert.Equal(t, "Height", objReq.Fields[0].Field)
}

func Test_Store_UnprocessableEntity_TwoFieldsMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPost, "/", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 2, len(objReq.Fields))
	assert.Equal(t, "Age", objReq.Fields[0].Field)
	assert.Equal(t, "Height", objReq.Fields[1].Field)
}

func Test_Store_UnprocessableEntity_ThreeFieldsMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPost, "/", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 3, len(objReq.Fields))
	assert.Equal(t, "Email", objReq.Fields[0].Field)
	assert.Equal(t, "Age", objReq.Fields[1].Field)
	assert.Equal(t, "Height", objReq.Fields[2].Field)
}

func Test_Store_UnprocessableEntity_FourFieldsMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPost, "/", `{
		"active": false,
		"height": 30
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 4, len(objReq.Fields))
	assert.Equal(t, "Firstname", objReq.Fields[0].Field)
	assert.Equal(t, "Lastname", objReq.Fields[1].Field)
	assert.Equal(t, "Email", objReq.Fields[2].Field)
	assert.Equal(t, "Age", objReq.Fields[3].Field)
}

func Test_Store_UnprocessableEntity_FiveFieldsMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPost, "/", `{
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 5, len(objReq.Fields))
	assert.Equal(t, "Firstname", objReq.Fields[0].Field)
	assert.Equal(t, "Lastname", objReq.Fields[1].Field)
	assert.Equal(t, "Email", objReq.Fields[2].Field)
	assert.Equal(t, "Age", objReq.Fields[3].Field)
	assert.Equal(t, "Height", objReq.Fields[4].Field)
}

func Test_Store_OK(t *testing.T) {
	objReq := responseOfUser{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPost, "/", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"age": 30,
		"height": 170,
		"active": false
    }`)

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusCreated, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, "Testing", objReq.Data.Firstname)
	assert.Equal(t, "Save Product", objReq.Data.Lastname)
	assert.Equal(t, "testing@testing.org", objReq.Data.Email)
	assert.Equal(t, 30, objReq.Data.Age)
	assert.Equal(t, 170, objReq.Data.Height)
	assert.Equal(t, false, objReq.Data.Active)
	assert.Equal(t, "02/01/2006 15:04:05", objReq.Data.CreatedDate)
}

/*************************** UPDATE ***************************/

func Test_Update_AuthorizationTokenMissing(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTestWithoutToken(http.MethodPut, "/1", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"age": 30,
		"height": 170,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnauthorized), objReq.Code)
	assert.Equal(t, message.UnauthorizedAction, objReq.Message)
}

func Test_Update_InternalError(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"age": 30,
		"height": 170,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusInternalServerError), objReq.Code)
	assert.Equal(t, message.InternalError, objReq.Message)
}

func Test_Update_NotFound(t *testing.T) {
	objReq := responseError{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/100", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"age": 30,
		"height": 170,
		"active": false
    }`)

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusNotFound), objReq.Code)
	assert.Equal(t, message.UserNotFound, objReq.Message)
}

func Test_Update_InvalidID(t *testing.T) {
	objReq := responseError{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/id", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"age": 30,
		"height": 170,
		"active": false
    }`)

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusBadRequest), objReq.Code)
	assert.Equal(t, message.InvalidID, objReq.Message)
}

func Test_Update_UnprocessableEntity_FirstnameMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"height": 170,
		"age": 40,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 1, len(objReq.Fields))
	assert.Equal(t, "Firstname", objReq.Fields[0].Field)
}

func Test_Update_UnprocessableEntity_LastnameMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
        "firstname": "Testing",
		"email": "testing@testing.org",
		"height": 170,
		"age": 40,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 1, len(objReq.Fields))
	assert.Equal(t, "Lastname", objReq.Fields[0].Field)
}

func Test_Update_UnprocessableEntity_EmailMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"height": 170,
		"age": 40,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 1, len(objReq.Fields))
	assert.Equal(t, "Email", objReq.Fields[0].Field)
}

func Test_Update_UnprocessableEntity_AgeMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"height": 170,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 1, len(objReq.Fields))
	assert.Equal(t, "Age", objReq.Fields[0].Field)
}

func Test_Update_UnprocessableEntity_HeightMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"age": 40,
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 1, len(objReq.Fields))
	assert.Equal(t, "Height", objReq.Fields[0].Field)
}

func Test_Update_UnprocessableEntity_TwoFieldsMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 2, len(objReq.Fields))
	assert.Equal(t, "Age", objReq.Fields[0].Field)
	assert.Equal(t, "Height", objReq.Fields[1].Field)
}

func Test_Update_UnprocessableEntity_ThreeFieldsMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 3, len(objReq.Fields))
	assert.Equal(t, "Email", objReq.Fields[0].Field)
	assert.Equal(t, "Age", objReq.Fields[1].Field)
	assert.Equal(t, "Height", objReq.Fields[2].Field)
}

func Test_Update_UnprocessableEntity_FourFieldsMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
		"active": false,
		"height": 30
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 4, len(objReq.Fields))
	assert.Equal(t, "Firstname", objReq.Fields[0].Field)
	assert.Equal(t, "Lastname", objReq.Fields[1].Field)
	assert.Equal(t, "Email", objReq.Fields[2].Field)
	assert.Equal(t, "Age", objReq.Fields[3].Field)
}

func Test_Update_UnprocessableEntity_FiveFieldsMissing(t *testing.T) {
	objReq := responseValidationsErrors{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
		"active": false
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnprocessableEntity, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnprocessableEntity), objReq.Code)
	assert.Equal(t, message.InvalidFields, objReq.Message)
	assert.Equal(t, 5, len(objReq.Fields))
	assert.Equal(t, "Firstname", objReq.Fields[0].Field)
	assert.Equal(t, "Lastname", objReq.Fields[1].Field)
	assert.Equal(t, "Email", objReq.Fields[2].Field)
	assert.Equal(t, "Age", objReq.Fields[3].Field)
	assert.Equal(t, "Height", objReq.Fields[4].Field)
}

func Test_Update_OK(t *testing.T) {
	objReq := responseOfUser{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPut, "/1", `{
        "firstname": "Testing",
		"lastname": "Save Product",
		"email": "testing@testing.org",
		"age": 30,
		"height": 170,
		"active": false
    }`)

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 1, objReq.Data.ID)
	assert.Equal(t, "Testing", objReq.Data.Firstname)
	assert.Equal(t, "Save Product", objReq.Data.Lastname)
	assert.Equal(t, "testing@testing.org", objReq.Data.Email)
	assert.Equal(t, 30, objReq.Data.Age)
	assert.Equal(t, 170, objReq.Data.Height)
	assert.Equal(t, false, objReq.Data.Active)
}

/*************************** DELETE  ***************************/

func Test_Delete_AuthorizationTokenMissing(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTestWithoutToken(http.MethodDelete, "/1", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnauthorized), objReq.Code)
	assert.Equal(t, message.UnauthorizedAction, objReq.Message)
}

func Test_Delete_InternalError(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTest(http.MethodDelete, "/1", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusInternalServerError), objReq.Code)
	assert.Equal(t, message.InternalError, objReq.Message)
}

func Test_Delete_NotFound(t *testing.T) {
	objReq := responseError{}
	r := createServer()

	req, rr := createRequestTest(http.MethodDelete, "/100", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusNotFound), objReq.Code)
	assert.Equal(t, message.UserNotFound, objReq.Message)
}

func Test_Delete_InvalidID(t *testing.T) {
	objReq := responseError{}
	r := createServer()

	req, rr := createRequestTest(http.MethodDelete, "/ID", "")
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusBadRequest), objReq.Code)
	assert.Equal(t, message.InvalidID, objReq.Message)
}

func Test_Delete_OK(t *testing.T) {
	objReq := struct {
		Data string `json:"data"`
	}{}
	r := createServer()

	req, rr := createRequestTest(http.MethodDelete, "/1", "")

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, message.UserDeleted, objReq.Data)
}

/*************************** UPDATE FIELDS ***************************/

func Test_UpdateFields_AuthorizationTokenMissing(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTestWithoutToken(http.MethodPatch, "/1", `{
		"lastname": "Save Product",
		"age": 30
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusUnauthorized, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusUnauthorized), objReq.Code)
	assert.Equal(t, message.UnauthorizedAction, objReq.Message)
}

func Test_UpdateFields_InternalError(t *testing.T) {
	objReq := responseError{}
	r := createServerFailDB()

	req, rr := createRequestTest(http.MethodPatch, "/1", `{
		"lastname": "Save Product",
		"age": 30
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusInternalServerError), objReq.Code)
	assert.Equal(t, message.InternalError, objReq.Message)
}

func Test_UpdateFields_NotFound(t *testing.T) {
	objReq := responseError{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPatch, "/100", `{
		"lastname": "Save Product",
		"age": 30
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusNotFound, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusNotFound), objReq.Code)
	assert.Equal(t, message.UserNotFound, objReq.Message)
}

func Test_UpdateFields_BadRequest(t *testing.T) {
	objReq := responseError{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPatch, "/1", `{
		"firstname": "Save Product",
		"height": 30
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusBadRequest), objReq.Code)
	assert.Equal(t, message.UserInvalidUpdate, objReq.Message)
}

func Test_UpdateFields_InvalidID(t *testing.T) {
	objReq := responseError{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPatch, "/ID", `{
		"firstname": "Save Product",
		"height": 30
    }`)
	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusBadRequest, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, web.StatusMsg(http.StatusBadRequest), objReq.Code)
	assert.Equal(t, message.InvalidID, objReq.Message)
}

func Test_UpdateFields_OK(t *testing.T) {
	objReq := responseOfUser{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPatch, "/1", `{
		"lastname": "Save Product",
		"age": 30
    }`)

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, "Save Product", objReq.Data.Lastname)
	assert.Equal(t, 30, objReq.Data.Age)
}

func Test_UpdateFields_OK_OnlyAge(t *testing.T) {
	objReq := responseOfUser{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPatch, "/1", `{
		"age": 30
    }`)

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, 30, objReq.Data.Age)
}

func Test_UpdateFields_OK_OnlyLastname(t *testing.T) {
	objReq := responseOfUser{}
	r := createServer()

	req, rr := createRequestTest(http.MethodPatch, "/1", `{
		"lastname": "Save Product"
    }`)

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)

	err := json.Unmarshal(rr.Body.Bytes(), &objReq)
	assert.Nil(t, err)
	assert.Equal(t, "Save Product", objReq.Data.Lastname)
}
