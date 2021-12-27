package main

import (
	"github.com/gin-gonic/gin"
	"github.com/miltonbernhardt/go-web/cmd/server/handler"
	"github.com/miltonbernhardt/go-web/internal/users"
)

func main() {
	r := gin.Default() // add middleware (Logger & Recovery)

	userRepository := users.NewRepository()
	userService := users.NewService(userRepository)
	userController := handler.NewUserController(userService)

	userGroup := r.Group("/users")
	userGroup.GET("/", userController.GetAll)
	userGroup.GET("/:id", userController.GetById)
	userGroup.POST("/save", userController.Save)
	userGroup.PUT("/:id", userController.Update())
	userGroup.DELETE("/:id", userController.Delete())
	userGroup.PATCH("/:id", userController.UpdateFields())

	_ = r.Run() // listen and serve on 0.0.0.0:8080 | "localhost:8080"
}
