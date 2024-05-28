package main

import (
	"log"
	"one-way-ticket/db"
	"one-way-ticket/routers"
)

func main() {
	err := db.Connect()
	if err != nil {
		return
	}
	defer db.Close()

	r := routers.SetupRouter()
	// listen and serve on 0.0.0.0:8080
	err = r.Run(":8080")
	if err != nil {
		log.Fatal(err.Error())
		return
	}
}
