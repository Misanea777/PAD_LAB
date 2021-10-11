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
