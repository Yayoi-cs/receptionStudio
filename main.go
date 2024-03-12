package main

import (
	"fmt"
	"receptionStudio/dbHelper"
	"receptionStudio/server"
)

func main() {
	//jsonHelper.DecodeTest()

	/*jsonData := []byte(`{
		"ProjectName": "Sample",
		"ProjectID": "1",
		"ProjectData": {
			"Name": "Parent1",
			"IsRecep": false,
			"Child": [
				{
					"Name": "Child1",
					"IsRecep": true,
					"Child": null
				},
				{
					"Name": "Child1",
					"IsRecep": true,
					"Child": null
				}
			]
		}
	}`)*/
	err := dbHelper.UpdateOldProject("2", "Testing from Golang", "")
	if err != nil {
		fmt.Println(err)
	}
	server.StartServer()
}
