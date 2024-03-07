package dbHelper

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

const (
	DbName = "receptionStudio"
)

func InsertIntoVerifyDB(email, hash, verifyCode string) error {
	/*
			mysql> CREATE TABLE verifyDB (
		       ->     id INT AUTO_INCREMENT PRIMARY KEY,
		       ->     email VARCHAR(255) NOT NULL,
		       ->     hash VARCHAR(255) NOT NULL,
		       ->     verify_code VARCHAR(255) NOT NULL
		       -> );
	*/
	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
		}
	}()

	var existingEmail string
	err = tx.QueryRow("SELECT email FROM verifyDB WHERE email = ?", email).Scan(&existingEmail)
	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		_ = tx.Rollback()
		log.Fatal(err)
		return err
	default:
		_, err := tx.Exec("DELETE FROM verifyDB WHERE email = ?", email)
		if err != nil {
			_ = tx.Rollback()
			log.Fatal(err)
			return err
		}
	}

	stmt, err := tx.Prepare("INSERT INTO verifyDB (email, hash, verify_code) VALUES (?, ?, ?)")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(email, hash, verifyCode)
	if err != nil {
		_ = tx.Rollback()
		log.Fatal(err)
		return err
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

func CheckVerifyCode(mail, verifyCode string) (bool, error) {
	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM verifyDB WHERE email = ? AND verify_code = ?")
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(mail, verifyCode).Scan(&count)
	if err != nil {
		log.Fatal(err)
		return false, err
	}

	return count > 0, nil
}

func MoveFromVerifyCodeToUserDB() {

}
