package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/miltonbernhardt/go-web/cmd/server/handler"
	"github.com/miltonbernhardt/go-web/docs"
	"github.com/miltonbernhardt/go-web/internal/users"
	"github.com/miltonbernhardt/go-web/internal/utils"
	"github.com/miltonbernhardt/go-web/pkg/store"
	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"
	"log"
	"os"
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
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error al cargar archivo .env")
	}

	r := gin.Default() // add middleware (Logger & Recovery)

	userRepository := users.NewRepository(store.New(store.FileType, store.FileNameUsers))
	userService := users.NewService(userRepository, utils.New())
	userController := handler.NewUserController(userService)

	docs.SwaggerInfo.Host = os.Getenv("HOST")
	r.GET("/doc/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	userGroup := r.Group("/users")
	userGroup.GET("/", userController.ValidateToken, userController.GetAll)
	userGroup.GET("/:id", userController.ValidateToken, userController.GetById)
	userGroup.POST("/", userController.ValidateToken, userController.Store)
	userGroup.PUT("/:id", userController.ValidateToken, userController.Update())
	userGroup.DELETE("/:id", userController.ValidateToken, userController.Delete())
	userGroup.PATCH("/:id", userController.ValidateToken, userController.UpdateFields())

	_ = r.Run("127.0.0.1:8081") // listen and serve on 0. 0.0.0:8080 | "localhost:8080"
}
