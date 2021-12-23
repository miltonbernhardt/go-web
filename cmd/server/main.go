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
	userGroup.GET("/", userController.GetUsers)
	userGroup.GET("/:id", userController.GetUserById)
	userGroup.POST("/save", userController.SaveUser)

	_ = r.Run() // listen and serve on 0.0.0.0:8080 | "localhost:8080"
}
