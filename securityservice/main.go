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
	usersAPI.Get("/status", controllers.EndpointStatus)
	usersAPI.Get("/count", controllers.OnlinePlayers)

	db.InitDB()
	validators.InitValidator()
	db.InitRedis()
	time.Sleep(time.Second * 20)
	db.PingEureka()

	go db.CheckDBConnection()

	app.Run(iris.Addr(":8090"))
}
