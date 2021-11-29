package models

import (
	"crypto"
	"database/sql"
	"encoding/hex"
	"errors"

	"github.com/PAD_LAB/db"
	"github.com/PAD_LAB/validators"
)

func RegisterUser(user validators.UserCredentials, tx *sql.Tx) (string, error) {
	var (
		count int
		ID    string
	)
	q := `
		SELECT COUNT(*) FROM Users
		WHERE
			Username = ?
	`
	res := tx.QueryRow(q, user.Username)
	res.Scan(&count)

	if count != 0 {
		return "", errors.New("username already registered")
	}

	hasher := crypto.SHA256.New()
	hasher.Write([]byte(user.Password))
	hp := hasher.Sum(nil)

	_, err := tx.Exec(`
		INSERT INTO Users (Username, Password)
		VALUES (?,?)
	
	`, user.Username, hex.EncodeToString(hp))
	if err != nil {
		return "", err
	}

	res = tx.QueryRow(`SELECT ID FROM Users WHERE Username = ?`, user.Username)
	err = res.Scan(&ID)
	if err != nil {
		return "", err
	}

	return ID, nil
}

func LoginUser(user validators.UserCredentials) (string, error) {
	var ID string
	q := `
		SELECT 
			ID
		FROM
			Users
		WHERE
			Username = ? AND Password = ?
	`
	hasher := crypto.SHA256.New()
	hasher.Write([]byte(user.Password))
	hp := hasher.Sum(nil)

	res := db.Mysql.QueryRow(q, user.Username, hex.EncodeToString(hp))
	err := res.Scan(&ID)
	if err != nil {
		return "", err
	}

	return ID, nil
}

func GetRegisteredUsersNumber() (int, error) {
	var count int
	q := `
		SELECT COUNT(*) FROM Users
	`
	res := db.Mysql.QueryRow(q)

	err := res.Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetUserBalance(userID string) (*validators.UserBalanceInfo, error) {
	var userBalance validators.UserBalanceInfo
	q := `
		SELECT ID, Balance FROM Users WHERE ID = ?
	`

	res := db.Mysql.QueryRow(q, userID)
	err := res.Scan(&userBalance.ID, &userBalance.Balance)
	if err != nil {
		return nil, err
	}

	return &userBalance, nil
}

func UpdateUserBalance(userID string, amount int) error {
	q := `
		UPDATE Users
		SET Balance=Balance + ?
		WHERE ID = ?
	`

	res, err := db.Mysql.Exec(q, amount, userID)
	affected, exc := res.RowsAffected()

	switch {
	case err != nil:
		return err
	case exc != nil:
		return exc
	case affected == 0:
		return errors.New("no rows affected when updating the balance")
	}

	return nil
}
