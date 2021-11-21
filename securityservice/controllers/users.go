package controllers

import (
	"fmt"
	"log"
	"time"

	"github.com/PAD_LAB/db"
	"github.com/PAD_LAB/models"
	"github.com/PAD_LAB/validators"
	"github.com/kataras/iris/v12"
)

var PlayingPlayers int

func Register(ctx iris.Context) {
	var (
		user validators.UserCredentials
	)
	err := ctx.ReadJSON(&user)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	err = validators.Validate.Struct(user)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("invalid input: username must be at least 9 characters and password at least 6")
		return
	}

	tx, err := db.Mysql.Begin()
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString(err.Error())
		return
	}

	defer func() {
		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
		}
	}()

	_, err = models.RegisterUser(user, tx)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("username taken")
		return
	}

	err = tx.Commit()
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("transaction not executed: " + err.Error())
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.WriteString("success")
}

func Login(ctx iris.Context) {
	var (
		user validators.UserCredentials
	)
	err := ctx.ReadJSON(&user)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString(err.Error())
		return
	}

	err = validators.Validate.Struct(user)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("invalid input: username must be at least 9 characters and password at least 6")
		return
	}

	uuid, err := models.LoginUser(user)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.WriteString("invalid credentials")
		return
	}

	token, err := models.CreateToken(uuid)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("error when generating token")
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(token)
}

func Auth(ctx iris.Context) {
	tokenAuth, err := models.ExtractTokenMetadata(ctx.GetHeader("Authorization"))
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "unauthorized",
		})
		return
	}

	userID, err := db.FetchAuth(tokenAuth)
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "unauthorized",
		})
		return
	}

	PlayingPlayers++

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(map[string]interface{}{
		"status":  "success",
		"message": userID,
	})
}

func EndpointStatus(ctx iris.Context) {
	var (
		count int
		err   error
	)

	res1 := db.RedisClient.ClientGetName().Name()
	fmt.Println(res1)
	time.Sleep(4 * time.Second)

	count, err = models.GetRegisteredUsersNumber()
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": err.Error(),
		})
		return
	}

	err = db.Mysql.Ping()
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "mysql db is down",
		})
		return
	}

	_, err = db.RedisClient.Ping().Result()
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "redis is down",
		})
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(map[string]interface{}{
		"status": "success",
		"message": map[string]interface{}{
			"users_count": count,
			"mysql":       "up",
			"redis":       "up",
		},
	})
}

func OnlinePlayers(ctx iris.Context) {
	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(map[string]interface{}{
		"status": "success",
		"message": map[string]interface{}{
			"current_playing_players": PlayingPlayers,
		},
	})
}
