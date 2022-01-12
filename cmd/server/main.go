package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/miltonbernhardt/go-web/cmd/server/handler"
	"github.com/miltonbernhardt/go-web/docs"
	"github.com/miltonbernhardt/go-web/internal/users"
	"github.com/miltonbernhardt/go-web/internal/utils"
	"github.com/miltonbernhardt/go-web/pkg/store"
	"github.com/miltonbernhardt/go-web/pkg/web"
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

	userRepository := users.NewRepository(store.New(store.FileType, store.FileNameUsers))
	userService := users.NewService(userRepository, utils.New())
	userController := handler.NewUserController(userService)
	loadRoutes(r, userController)

	docs.SwaggerInfo.Host = os.Getenv("HOST")

	_ = r.Run("127.0.0.1:8081") // listen and serve on 0. 0.0.0:8080 | "localhost:8080"
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
		log.Fatal(web.FailedToLoadEnv)
	}
}

func caller() func(*runtime.Frame) (function string, file string) {
	return func(f *runtime.Frame) (function string, file string) {
		p, _ := os.Getwd()

		return "", fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
	}
}

func loadLog() {
	gin.SetMode(gin.ReleaseMode)
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
