package main

import (
	"os"
	"testing"
)

/*
	Setup Test
	This file will run before any other test will run.

	It must named be as `setup_test.go`

	We are going to setup the TestRepository which mock Postgres Database for testing
*/

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}
