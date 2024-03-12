package dbHelper

import (
	"database/sql"
	"fmt"
	"log"
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

/*
mysql> CREATE TABLE userTable (
    -> id INT AUTO_INCREMENT PRIMARY KEY,
    -> email VARCHAR(255) NOT NULL,
    -> hash VARCHAR(255) NOT NULL,
    -> availableProject VARCHAR(255) NULL
    -> );
*/

func CheckExistUserTable(email string) (bool, error) {
	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer func() {
		if err := recover(); err != nil {
			_ = tx.Rollback()
		}
	}()
	var existingEmail string
	err = tx.QueryRow("SELECT email FROM userTable WHERE email = ?", email).Scan(&existingEmail)
	switch {
	case err == sql.ErrNoRows:
		return false, nil
		break
	case err != nil:
		_ = tx.Rollback()
		log.Fatal(err)
		return false, err
	default:
		return true, nil
	}
	return false, nil
}

func InsertIntoUserTable(email, hash string, isOauth bool) error {
	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer db.Close()

	query := "INSERT INTO userTable (email,hash,oauth,availableProject) VALUES (?,?,?,?)"

	_, err = db.Exec(query, email, hash, isOauth, "")
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

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

func CheckMailWithHash(mail, hash string) (bool, error) {
	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT COUNT(*) FROM userTable WHERE email = ? AND hash = ? AND oauth = 0")
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(mail, hash).Scan(&count)
	if err != nil {
		log.Fatal(err)
		return false, err
	}
	return count > 0, nil
}

func MoveFromVerifyDBToUserDB(mail string) error {
	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("SELECT hash FROM verifyDB WHERE email = ?")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	var hash string
	err = stmt.QueryRow(mail).Scan(&hash)
	if err != nil {
		log.Fatal(err)
		return err
	}
	query := "DELETE FROM verifyDB WHERE email = ?"
	_, err = db.Exec(query, mail)
	if err != nil {
		return err
	}
	InsertIntoUserTable(mail, hash, false)
	return nil
}

func AddAvailableProject(pn int, mail string) error {
	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("select availableProject from userTable where email = ?")
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer stmt.Close()
	var currentAvailableProject string
	err = stmt.QueryRow(mail).Scan(&currentAvailableProject)
	if err != nil {
		log.Fatal(err)
		return err
	}
	newAvailableProject := currentAvailableProject + "," + strconv.Itoa(pn)

	query := "UPDATE userTable SET availableProject = ? where email = ?"
	_, err = db.Exec(query, newAvailableProject, mail)
	if err != nil {
		return err
	}
	return nil
}
