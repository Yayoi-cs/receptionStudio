package dbHelper

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"strings"

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

	query := "INSERT INTO projectDB (creater,pnu,pna,pd) VALUES (?,?,?,?)"

	_, err = db.Exec(query, mail, nextPn, pna, nil)
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

func DeleteOldProject(mail, num string) error {

	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer db.Close()

	query := "DELETE FROM projectDB WHERE pnu = ? AND creater = ?"

	_, err = db.Exec(query, num, mail)
	if err != nil {
		fmt.Println("ERROR WHILE INSERT", err)
		return err
	}

	return nil
}

func ReadOldProject(num string) (string, error) {

	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer db.Close()

	stmt, err := db.Prepare("select pd from projectDB where pnu = ?")
	if err != nil {
		fmt.Println("ERROR WHILE SELECT", err)
		return "", err
	}
	defer stmt.Close()
	var pd string
	err = stmt.QueryRow(num).Scan(&pd)
	if err != nil {
		fmt.Println("ERROR WHILE SELECT", err)
		return "", err
	}
	return pd, nil
}
func pnaFromPnu(num string) (string, error) {

	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	defer db.Close()

	stmt, err := db.Prepare("select pna from projectDB where pnu = ?")
	if err != nil {
		fmt.Println("ERROR WHILE SELECT", err)
		return "", err
	}
	defer stmt.Close()
	var pna string
	err = stmt.QueryRow(num).Scan(&pna)
	if err != nil {
		fmt.Println("ERROR WHILE SELECT", err)
		return "", err
	}
	return pna, nil
}

type ProjectObject struct {
	Pnu string `json:"pnu"`
	Pna string `json:"pna"`
}

func AvailableProjectInformation(mail string) ([]byte, error) {
	DbUser, DbPassWord := DBconfig()
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DbUser, DbPassWord, DbName))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare("select availableProject from userTable where email = ?")
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer stmt.Close()
	var currentAvailableProject sql.NullString
	err = stmt.QueryRow(mail).Scan(&currentAvailableProject)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	if !currentAvailableProject.Valid {
		return nil, nil
	}
	projectList := strings.Split(currentAvailableProject.String, ",")
	var returnObjects []ProjectObject
	for _, s := range projectList {
		if s != "" {
			pna, err := pnaFromPnu(s)
			if err != nil {
				return nil, err
			}
			returnObjects = append(returnObjects, ProjectObject{
				Pnu: s,
				Pna: pna,
			})
		}
	}
	jsonData, err := json.Marshal(returnObjects)
	if err != nil {
		return nil, err
	}
	return jsonData, nil
}
