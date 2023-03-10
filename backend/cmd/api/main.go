package main

import (
	"fmt"
	"houseware---backend-engineering-octernship-KunalSin9h/data"
	"log"
	"net/http"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Config struct {
	DB     *gorm.DB
	Models data.Models
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

	dbPool := connectDatabase()

	app := Config{
		DB:     dbPool,
		Models: data.New(dbPool),
	}

	server := &http.Server{
		Addr:    fmt.Sprintf(":%s", PORT),
		Handler: app.routes(),
	}

	log.Println("Starting Authentication server at port:", PORT)
	log.Fatal(server.ListenAndServe())
}

func connectDatabase() *gorm.DB {
	numberOfTry := 0
	numberOfTryLimit := 5

	for {
		db, err := gorm.Open(postgres.Open(DSN), &gorm.Config{})

		if err != nil {
			numberOfTry++
			log.Printf("@MAIN Trying to connect to Postgres Database...[%d/%d]", numberOfTry, numberOfTryLimit)
		} else {
			log.Println("@MAIN Successfully Connected to Postgres Database")
			return db
		}

		if numberOfTry >= numberOfTryLimit {
			log.Fatal("@MAIN Failed to Connect to Postgres Database")
		}

		holdTime := numberOfTry * numberOfTry // numberOfTry ^ 2
		log.Printf("@MAIN Retrying to connect in %d sec", holdTime)
		time.Sleep(time.Duration(holdTime) * time.Second)
	}
}
