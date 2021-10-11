package db

import (
	"context"
	"database/sql"
	"flag"
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

	db, err := sql.Open("mysql", "root:root@/")
	if err != nil {
		panic(err.Error())
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

	db, err = sql.Open("mysql", "root:root@/users")
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
