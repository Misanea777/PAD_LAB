package main

import (
	"fmt"
	"time"

	"github.com/PAD_LAB/controllers"
	"github.com/PAD_LAB/db"
	"github.com/PAD_LAB/models"
	"github.com/kataras/iris/v12"
)

func main() {
	app := iris.New()

	paymentAPI := app.Party("/payment")
	paymentAPI.Post("/init", controllers.Init)
	paymentAPI.Post("/confirm", controllers.Confirm)

	db.InitDB()
	time.Sleep(time.Second * 20)
	db.PingEureka()

	go Tmp()

	app.Run(iris.Addr(":8091"))
}

func Tmp() {
	for {
		fmt.Println("Payment parser start")
		var payments []models.Payment
		q := `
			SELECT
				ID,
				UserID,
				Amount,
				DateCreated,
				Status
			FROM Payments
			WHERE
			Status = ?
		`

		res, err := db.MySQL.Query(q, "init")
		if err != nil {
			fmt.Println("error selecting payments, error msg: ", err.Error())
			continue
		}

		for res.Next() {
			var payment models.Payment
			err = res.Scan(
				&payment.ID,
				&payment.UserID,
				&payment.Amount,
				&payment.DateCreated,
				&payment.Status,
			)

			if err != nil {
				fmt.Println("error scanning payments, err msg :", err.Error())
				continue
			}

			payments = append(payments, payment)
		}

		for _, payment := range payments {
			tsTime := time.Unix(int64(payment.DateCreated), 0)
			diff := time.Since(tsTime).Seconds()
			if diff > 60 {
				err := models.UpdatePaymentStatus(payment.ID, "expired")
				if err != nil {
					fmt.Println("error expiring the payment, err msg: ", err.Error())
					continue
				}
			}
		}

		time.Sleep(time.Second * 5)
	}
}
