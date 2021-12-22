package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

func GetAll(ctx *gin.Context) {
	var users []User
	usersJson, err := os.ReadFile("./users.json")

	if err != nil {
		ctx.String(500, "error: no se pudo leer el archivo")
	}

	err = json.Unmarshal(usersJson, &users)

	if err != nil {
		ctx.String(500, "error: no se pudo descifrar el archivo")
	}

	fmt.Println(users)

	ctx.JSON(200, users)
}

func Greeting(ctx *gin.Context) {
	ctx.String(200, "Hola Milton!")

}
func main() {
	r := gin.Default() // add middleware (Logger & Recovery)
	r.GET("/greeting", Greeting)
	r.GET("/users", GetAll)
	_ = r.Run() // listen and serve on 0.0.0.0:8080 | "localhost:8080"
}

type User struct {
	Id            string
	Nombre        string `json:"nombre"`
	Apellido      string `json:"apellido"`
	Email         string `json:"email"`
	Altura        int    `json:"edad"`
	Activo        bool   `json:"activo"`
	FechaCreacion string `json:"fecha_creacion"`
}
