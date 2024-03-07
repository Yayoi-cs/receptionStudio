package jsonHelper

import (
	"encoding/json"
	"fmt"
)

func DecodeTest() {
	jsonData := []byte(`{
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
	}`)
	var project Project
	err := json.Unmarshal(jsonData, &project)
	if err != nil {
		fmt.Println("JSON DECODE ERROR:", err)
		return
	}
	fmt.Printf("%+v\n", project)
}
