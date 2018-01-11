// money
package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/leekchan/accounting"
)

type Money struct {
	MoneyId   int
	UserId    int
	Amount    float64
	Category  string
	CreatedAt int64
}

type Balance struct {
	CurrentBalance float64
}

func InsertAmount(db *sql.DB, userId int, amount float64, category string) (int, error) {
	stmtIns, err := db.Prepare("INSERT INTO money (user_id, amount, category) VALUES (?, ?, ?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	res, err := stmtIns.Exec(userId, amount, category)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInsertId), nil
}

func GetBalance(db *sql.DB, userId int) (*Balance, error) {
	var b Balance
	err := db.QueryRow("SELECT COALESCE(SUM(amount), 0) AS current_balance FROM money WHERE user_id = ?", userId).Scan(&b.CurrentBalance)
	if err != nil {
		return nil, err
	}

	return &b, nil
}

func ConvertBalance(amount float64) string {
	ac := accounting.Accounting{Symbol: "$", Precision: 2}
	return ac.FormatMoneyFloat64(amount)
}
