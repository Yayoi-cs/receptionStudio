package main

import (
	"receptionStudio/server"
)

func main() {
	//jsonHelper.DecodeTest()

	//err := dbHelper.UpdateOldProject("2", "Testing from Golang", "")
	//if err != nil {
	//fmt.Println(err)
	//}
	server.StartServer()
}
