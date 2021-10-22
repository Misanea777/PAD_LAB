package main

import (
	"time"

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
	usersAPI.Get("/auth", controllers.Auth)

	db.InitDB()
	validators.InitValidator()
	db.InitRedis()
	time.Sleep(time.Second * 60)
	db.PingEureka()

	app.Run(iris.Addr(":8090"))
}
