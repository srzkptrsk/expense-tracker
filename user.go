// user
package main

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

type User struct {
	UserId int
}

func GetUserById(db *sql.DB, userId int) (*User, error) {
	var u User
	err := db.QueryRow("SELECT * FROM user WHERE user_id = ?", userId).Scan(&u.UserId)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

func InsertUser(db *sql.DB, userId int) (int, error) {
	stmtIns, err := db.Prepare("INSERT INTO user (user_id) VALUES (?)")
	if err != nil {
		panic(err.Error())
	}
	defer stmtIns.Close()

	res, err := stmtIns.Exec(userId)
	if err != nil {
		return 0, err
	}

	lastInsertId, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(lastInsertId), nil
}
