package main

import (
	"github.com/PAD_LAB/controllers"
	"github.com/PAD_LAB/db"
	"github.com/PAD_LAB/validators"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	usersAPI := app.Party("/user")
	usersAPI.Post("/register", controllers.Register)
	usersAPI.Get("/login", controllers.Login)

	db.InitDB()
	validators.InitValidator()
	// db.InitRedis()

	app.Run(iris.Addr(":8080"))
}
