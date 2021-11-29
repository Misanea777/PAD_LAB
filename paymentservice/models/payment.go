package models

import (
	"errors"

	"github.com/PAD_LAB/db"
)

func InitPayment(args PaymentInitArgs) error {
	q := `
		INSERT INTO Payments (UserID, Amount, Status)
		VALUES (?,?,?)
	`

	_, err := db.MySQL.Exec(q, args.UserID, args.Amount, "init")
	if err != nil {
		return err
	}
	return nil
}

type PaymentInitArgs struct {
	UserID string `json:"id"`
	Amount int    `json:"balance"`
}

type Payment struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Amount      int    `json:"amount"`
	DateCreated int    `json:"date_created"`
	Status      string `json:"status"`
}

func UpdatePaymentStatus(payID, status string) error {
	q := `
		UPDATE Payments
		SET Status = ?
		WHERE ID = ?
	`

	res, err := db.MySQL.Exec(q, status, payID)
	affected, exc := res.RowsAffected()

	switch {
	case err != nil:
		return err
	case exc != nil:
		return exc
	case affected == 0:
		return errors.New("no payment updated")
	}

	return nil
}

func GetPayment(payID string) (*PaymentInitArgs, error) {
	var payment PaymentInitArgs

	q := `
		SELECT UserID, Amount FROM Payments WHERE ID = ?
	`

	res := db.MySQL.QueryRow(q, payID)

	err := res.Scan(
		&payment.UserID,
		&payment.Amount,
	)
	if err != nil {
		return nil, err
	}

	return &payment, nil
}
