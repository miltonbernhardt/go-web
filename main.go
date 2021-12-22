package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() // add middleware (Logger & Recovery)
	r.GET("/hello", Hello)
	r.GET("/users", GetUsers)
	r.GET("/users/:id", GetUserById)
	_ = r.Run() // listen and serve on 0.0.0.0:8080 | "localhost:8080"
}
