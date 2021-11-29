package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	// "time"

	"github.com/PAD_LAB/db"
	"github.com/PAD_LAB/models"
	"github.com/PAD_LAB/validators"
	"github.com/kataras/iris/v12"
	"github.com/valyala/fasthttp"
	eureka_client "github.com/xuanbo/eureka-client"
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
	// time.Sleep(4 * time.Second)

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

func BalanceTransaction(ctx iris.Context) {
	var userBalance *validators.UserBalanceInfo

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

	amount, err := ctx.URLParamInt("amount")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "couldnt get amount param, err msg: " + err.Error(),
		})
		return
	}

	userBalance, err = models.GetUserBalance(userID)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "couldnt get amount param, err msg: " + err.Error(),
		})
		return
	}

	fmt.Println(amount)
	fmt.Println(userBalance.Balance)
	var gatewayApp eureka_client.Application
	for _, app := range db.EurekaClient.Applications.Applications {
		if app.Name == "GATEWAY" {
			gatewayApp = app
		}
	}

	if amount < 0 {
		if userBalance.Balance < -amount {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(map[string]interface{}{
				"status":  "failed",
				"message": "unsufficient funds",
			})
			return
		} else if userBalance.Balance >= -amount {
			err = models.UpdateUserBalance(userID, amount)
			if err != nil {
				ctx.StatusCode(iris.StatusBadRequest)
				ctx.JSON(map[string]interface{}{
					"status":  "failed",
					"message": "failed to update the balance, err msg: " + err.Error(),
				})
				return
			}
		}
	} else if amount > 0 {
		// TODO call payment service
		req := fasthttp.AcquireRequest()
		res := fasthttp.AcquireResponse()
		defer fasthttp.ReleaseRequest(req)
		defer fasthttp.ReleaseResponse(res)

		userBalance.Balance = amount
		bodyBytes, err := json.Marshal(userBalance)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(map[string]interface{}{
				"status":  "failed",
				"message": "failed to marshal json, err msg : " + err.Error(),
			})
			return
		}

		req.AppendBody(bodyBytes)
		req.Header.SetMethod(fasthttp.MethodPost)
		req.SetRequestURI(gatewayApp.Instances[0].HomePageURL + "payment/init")

		err = fasthttp.Do(req, res)
		if err != nil {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(map[string]interface{}{
				"status":  "failed",
				"message": "failed to send request to payment service, err msg: " + err.Error(),
			})
			return
		}

		if res.StatusCode() != iris.StatusOK {
			ctx.StatusCode(iris.StatusInternalServerError)
			ctx.JSON(map[string]interface{}{
				"status":  "failed",
				"message": "failed to init payment request, err msg: " + string(res.Body()),
			})
			return
		}

	} else {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "invalid amount value",
		})
		return
	}
}

func PaymentConfirm(ctx iris.Context) {
	body, err := ctx.GetBody()
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "failed to get pay confirm body, err msg : " + err.Error(),
		})
		return
	}

	var userBalance validators.UserBalanceInfo
	err = json.Unmarshal(body, &userBalance)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "failed to unmarshal pay confirm body, err msg : " + err.Error(),
		})
		return
	}

	err = models.UpdateUserBalance(userBalance.ID, userBalance.Balance)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "failed to update user balance, err msg : " + err.Error(),
		})
		return
	}

	ctx.StatusCode(iris.StatusOK)
}
