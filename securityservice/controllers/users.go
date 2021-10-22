package controllers

import (
	"log"

	"github.com/PAD_LAB/db"
	"github.com/PAD_LAB/models"
	"github.com/PAD_LAB/validators"
	"github.com/kataras/iris/v12"
)

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

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(map[string]interface{}{
		"status":  "success",
		"message": userID,
	})
}
