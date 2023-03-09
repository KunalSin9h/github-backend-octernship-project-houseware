package main

import (
	"fmt"
	"log"
	"net/http"
)

type Config struct{}

var PORT = "3000"

func main() {

	app := Config{}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	log.Println("Starting Authentication server at port:", PORT)
	log.Fatal(server.ListenAndServe())
}
