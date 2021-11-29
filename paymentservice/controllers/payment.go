package controllers

import (
	"encoding/json"

	"github.com/PAD_LAB/db"
	"github.com/PAD_LAB/models"
	"github.com/kataras/iris/v12"
	"github.com/valyala/fasthttp"
	eureka_client "github.com/xuanbo/eureka-client"
)

func Init(ctx iris.Context) {
	var (
		args models.PaymentInitArgs
	)
	body, err := ctx.GetBody()
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	err = json.Unmarshal(body, &args)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		return
	}

	err = models.InitPayment(args)
}

func Confirm(ctx iris.Context) {
	id := ctx.URLParam("pay_id")

	err := models.UpdatePaymentStatus(id, "complete")
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "payment id is invalid",
		})
		return
	}

	var payment *models.PaymentInitArgs

	payment, err = models.GetPayment(id)
	if err != nil {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "failed getting payment, err msg: " + err.Error(),
		})
	}

	var gatewayApp eureka_client.Application
	for _, app := range db.EurekaClient.Applications.Applications {
		if app.Name == "GATEWAY" {
			gatewayApp = app
		}
	}

	req := fasthttp.AcquireRequest()
	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseRequest(req)
	defer fasthttp.ReleaseResponse(res)

	bodyBytes, err := json.Marshal(payment)
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
	req.SetRequestURI(gatewayApp.Instances[0].HomePageURL + "user/payment/confirm")

	err = fasthttp.Do(req, res)
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "failed to send confirmation request, err msg : " + err.Error(),
		})
		return
	}

	if res.StatusCode() != iris.StatusOK {
		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(map[string]interface{}{
			"status":  "failed",
			"message": "failed to add funds to user acc, err msg" + string(res.Body()),
		})
		_ = models.UpdatePaymentStatus(id, "reversed")
		return
	}
}
