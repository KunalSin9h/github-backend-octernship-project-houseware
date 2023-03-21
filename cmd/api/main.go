package main

import (
	"fmt"
	"houseware---backend-engineering-octernship-KunalSin9h/data"
	"log"
	"net/http"
	"os"
)

/*
data.Repository is an interface which we will use to put a real Postgres database when deploying in producing
and put a mock postgres database for testing
*/
type Config struct {
	Repo data.Repository
}

var (
	PORT       = os.Getenv("PORT")
	DSN        = os.Getenv("DSN")
	JWT_SECRET = os.Getenv("JWT_SECRET")
)

func init() {
	if PORT == "" {
		log.Println("@MAIN  Missing PORT in Env. Using 5000")
		PORT = "5000"
	}

	if DSN == "" {
		log.Println("@MAIN Missing Database Connection String (DSN) in Env. Using postgres://local:local@localhost:5432/local")
		DSN = "postgres://local:local@localhost:5432/local"
	}

	if JWT_SECRET == "" {
		log.Println("@MAIN Missing JWT Secret in Env. Using $ecret")
		JWT_SECRET = "$ecret"
	}

}

func main() {

	pool := data.ConnectDatabase(DSN)

	app := Config{
		Repo: data.NewPostgresRepository(pool), // Real Postgres database connection
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	log.Println("Starting Authentication server at port:", PORT)
	log.Fatal(server.ListenAndServe())
}
