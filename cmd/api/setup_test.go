package main

import (
	"houseware---backend-engineering-octernship-KunalSin9h/data"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

/*
	Setup Test
	This file will run before any other test will run.

	It must named be as `setup_test.go`

	We are going to setup the PostgresTestRepository which mock Postgres Database for testing
*/

var router *gin.Engine // package level variable used by `handlers_test.go`

func TestMain(m *testing.M) {

	testApp := Config{
		Repo: data.NewPostgresTestRepository(nil),
	}

	router = testApp.routes()

	os.Exit(m.Run())
}
