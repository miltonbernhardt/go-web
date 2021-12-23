package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() // add middleware (Logger & Recovery)
	r.GET("/hello", Hello)
	userGroup := r.Group("/users")
	userGroup.GET("/", GetUsers)
	userGroup.GET("/:id", GetUserById)
	_ = r.Run() // listen and serve on 0.0.0.0:8080 | "localhost:8080"
}
