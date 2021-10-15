package controllers

import (
	"log"
	"os"
	"time"

	"github.com/PAD_LAB/db"
	"github.com/PAD_LAB/models"
	"github.com/PAD_LAB/validators"
	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
	"github.com/twinj/uuid"
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

func CreateToken(ID string) (*validators.TokenDetails, error) {
	td := &validators.TokenDetails{}
	td.AtExpires = time.Now().Add(time.Minute * 60).Unix()
	td.AccessUuid = uuid.NewV4().String()

	td.RtExpires = time.Now().Add(time.Hour * 24 * 7).Unix()
	td.RefreshUuid = uuid.NewV4().String()

	var err error
	os.Setenv("ACCESS_SECRET", "jdnfksdmfksd")
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["access_uuid"] = td.AccessUuid
	atClaims["user_id"] = ID
	atClaims["exp"] = td.AtExpires
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	td.AccessToken, err = at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return nil, err
	}

	os.Setenv("REFRESH_SECRET", "mcmvmkmsdnfsdmfdsjf")
	rtClaims := jwt.MapClaims{}
	rtClaims["refresh_uuid"] = td.RefreshUuid
	rtClaims["user_id"] = ID
	rtClaims["exp"] = td.RtExpires
	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	td.RefreshToken, err = rt.SignedString([]byte(os.Getenv("REFRESH_SECRET")))
	if err != nil {
		return nil, err
	}
	err = db.SaveToken(ID, td)
	if err != nil {
		return nil, err
	}

	return td, nil
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

	token, err := CreateToken(uuid)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.WriteString("error when generating token")
		return
	}

	ctx.StatusCode(iris.StatusOK)
	ctx.JSON(token)
}
