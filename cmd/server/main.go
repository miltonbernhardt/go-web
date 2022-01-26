package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/miltonbernhardt/go-web/cmd/server/handler"
	"github.com/miltonbernhardt/go-web/docs"
	"github.com/miltonbernhardt/go-web/internal/users"
	"github.com/miltonbernhardt/go-web/internal/utils"
	"github.com/miltonbernhardt/go-web/pkg/message"
	"github.com/miltonbernhardt/go-web/pkg/store"
	log "github.com/sirupsen/logrus"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"os"
	"runtime"
	"strings"
)

// @title           MELI Bootcamp API
// @version         1.0
// @description     This API Handle MELI Products.
// @termsOfService  https://developers.mercadolibre.com.ar/es_ar/terminos-y-condiciones
// @contact.name    API Support
// @contact.url     https://developers.mercadolibre.com.ar/support
// @license.name    Apache 2.0
// @license.url     http://www.apache.org/licenses/LICENSE-2.0.html
func main() {

	loadEnv()
	loadLog()

	r := gin.Default() // add middleware (Logger & Recovery)

	db, err := initDB()

	var userRepository users.Repository

	if err != nil {
		log.Error(err)
		userRepository = users.NewRepositoryFile(store.New(store.FileType, store.FileNameUsers))
	} else {
		userRepository = users.NewRepository(db)
	}

	userService := users.NewService(userRepository, utils.New())
	userController := handler.NewUserController(userService)
	loadRoutes(r, userController)

	docs.SwaggerInfo.Host = os.Getenv("HOST")

	ip := os.Getenv("IP_ADDRESS")

	if err = r.Run(ip); err != nil {
		panic(err)
	}

	log.Info("server started at ", ip)
}

func loadRoutes(r *gin.Engine, userController handler.UserController) {
	r.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userGroup := r.Group("/users")
	userGroup.GET("/", userController.ValidateToken, userController.GetAll())
	userGroup.GET("/:id", userController.ValidateToken, userController.GetById())
	userGroup.POST("/", userController.ValidateToken, userController.Store())
	userGroup.PUT("/:id", userController.ValidateToken, userController.Update())
	userGroup.DELETE("/:id", userController.ValidateToken, userController.Delete())
	userGroup.PATCH("/:id", userController.ValidateToken, userController.UpdateFields())
}

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(message.FailedToLoadEnv)
	}
}

func caller() func(*runtime.Frame) (function string, file string) {
	return func(f *runtime.Frame) (function string, file string) {
		p, _ := os.Getwd()

		return "", fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
	}
}

func loadLog() {
	//gin.SetMode(gin.ReleaseMode)
	log.SetReportCaller(true)
	log.SetFormatter(&log.JSONFormatter{
		//log.SetFormatter(&log.TextFormatter{
		//DisableColors: false,
		//FullTimestamp: true,
		CallerPrettyfier: caller(),
		FieldMap: log.FieldMap{
			log.FieldKeyFile: "caller",
		},
		PrettyPrint: true,
	})

}

func initDB() (*sql.DB, error) {
	database := os.Getenv("DATABASE")
	userDB := os.Getenv("DB_USER")
	passDB := os.Getenv("DB_PASS")

	dataSource := fmt.Sprintf("%v:%v/%v", userDB, passDB, database)
	db, err := sql.Open("mysql", dataSource)

	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}

	log.Println("database Configured")

	return db, nil
}
