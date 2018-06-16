package main

import (
	"log"
	"net/http"

	"github.com/Dacode45/addressbook/server"
)

func init() {
	server.SetupDB()
}

func main() {
	if err := http.ListenAndServe(":8080", server.SetupRouter()); err != nil {
		log.Fatal(err)
	}

}
