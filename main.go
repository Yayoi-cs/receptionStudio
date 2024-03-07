package main

import (
	"receptionStudio/jsonHelper"
	"receptionStudio/server"
)

func main() {
	jsonHelper.DecodeTest()
	server.StartServer()
}
