package dbHelper

import (
	"database/sql"
	"fmt"
	"log"
)

func Seed() error {
	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	var pk int
	err = db.QueryRow("SELECT pk FROM projectDB WHERE pk = 1").Scan(&pk)
	if err == sql.ErrNoRows {
		fmt.Println("Insert Seed Data")
		query := "INSERT INTO projectDB (pk,creater,pnu,pna,pd) VALUES (?,?,?,?,?)"
		_, err = db.Exec(query, 1, "sample@example.com", 1, "Sample Project", nil)
		if err != nil {
			return err
		}
		return nil
	} else if err != nil {
		return err
	}

	return nil

}
