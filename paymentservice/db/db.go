package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var (
	paymentDBName = "payments"
	MySQL         *sql.DB
)

func InitDB() {
	db, _ := sql.Open("mysql", "root:root@tcp(goDB)/")
	for {
		err := db.Ping()
		if err != nil {
			fmt.Println("db not up, wait 10 sec")
			time.Sleep(5 * time.Second)
			continue
		} else {
			break
		}
	}

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+paymentDBName)
	if err != nil {
		panic(err.Error())
	}
	_, exc := res.RowsAffected()
	if exc != nil {
		panic(exc.Error())
	}
	db.Close()

	db, err = sql.Open("mysql", "root:root@tcp(goDB)/payments")
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = PaymentSelfMadeShitTransaction(db)
	if err != nil {
		panic(err.Error())
	}

	MySQL = db
}

func PaymentSelfMadeShitTransaction(sql *sql.DB) error {
	q := `
		CREATE TABLE IF NOT EXISTS
		Payments(
			ID VARCHAR(255) primary key DEFAULT (uuid()),
			UserID Text NOT NULL,
			Amount int NOT NULL,
			DateCreated INT DEFAULT (unix_timestamp()),
			Status VARCHAR(16) NOT NULL
		)
	`

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()

	_, err := sql.ExecContext(ctx, q)
	if err != nil {
		return err
	}

	return nil
}
