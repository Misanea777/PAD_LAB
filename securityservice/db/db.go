package db

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

var Mysql *sql.DB
var dbname = "users"

func InitDB() {
	// var (
	// 	migrationDir = flag.String("migration.files", "./migrations", "Directory where the migration files are located?")
	// )

	flag.Parse()

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
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		panic(err.Error())
	}
	_, exc := res.RowsAffected()
	if exc != nil {
		panic(exc.Error())
	}
	db.Close()

	db, err = sql.Open("mysql", "root:root@tcp(goDB)/users")
	if err != nil {
		panic(err.Error())
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	err = SelfMadeShittyTransaction(db)
	if err != nil {
		panic(err.Error())
	}

	Mysql = db
}

func SelfMadeShittyTransaction(sql *sql.DB) error {
	q := `
	CREATE TABLE IF NOT EXISTS
	Users(
		ID VARCHAR(255) primary key DEFAULT (uuid()),
		Username Text NOT NULL,
		Password Text NOT NULL,
		DateCreated INT DEFAULT (unix_timestamp())
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

func CheckDBConnection() {
	for {
		err := Mysql.Ping()
		if err != nil {
			fmt.Println("db down, try to reconnect...")
			InitDB()
		}
		fmt.Println("dp up, sleep 30s...")

		time.Sleep(30 * time.Second)
	}
}
