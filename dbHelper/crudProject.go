package dbHelper

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func InsertNewProject(mail, pna string) error {
	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	stmt, err := db.Prepare("select pnu from projectDB order by pk desc limit 1")
	if err != nil {
		fmt.Println("ERROR WHILE SELECT", err)
		return err
	}
	defer stmt.Close()
	var pn int
	err = stmt.QueryRow().Scan(&pn)
	if err != nil {
		fmt.Println("ERROR WHILE SELECT", err)
		return err
	}
	nextPn := pn + 1

	query := "INSERT INTO projectDB (pnu,pna,pd) VALUES (?,?,?)"

	_, err = db.Exec(query, nextPn, pna, nil)
	if err != nil {
		fmt.Println("ERROR WHILE INSERT", err)
		return err
	}
	err = AddAvailableProject(nextPn, mail)
	if err != nil {
		fmt.Println("ERROR", err)
		return err
	}
	return nil

}

func UpdateOldProject(num, name, data string) error {

	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		fmt.Println("ERROR WHILE CONNECT MYSQL:", err)
		return err
	}
	defer db.Close()
	query := "UPDATE projectDB set pna = ?,pd = ? where pnu = ?"

	_, err = db.Exec(query, name, nilOrString(data), num)
	if err != nil {
		fmt.Println("ERROR WHILE UPDATE", err)
		return err
	}
	return nil
}

func nilOrString(data string) sql.NullString {
	var nullStr sql.NullString
	if data == "" {
		nullStr.Valid = false
		return nullStr
	}
	nullStr = sql.NullString{data, true}
	return nullStr
}
