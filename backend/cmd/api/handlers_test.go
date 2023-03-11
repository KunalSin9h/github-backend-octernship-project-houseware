package main

import (
	"bytes"
	"encoding/json"
	"houseware---backend-engineering-octernship-KunalSin9h/data"
	"net/http"
	"net/http/httptest"
	"testing"
)

var db = data.ConnectDatabase("postgres://local:local@localhost:5432/local")

var testApp = Config{
	DB:     db,
	Models: data.New(db),
}

func Test_Login_Success(t *testing.T) {
	route := testApp.routes()
	route.POST("/login", testApp.login)

	var reqPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	reqPayload.Username = "user1"
	reqPayload.Password = "user1"

	reqBytes, _ := json.Marshal(reqPayload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBytes))
	req.Header.Add("Content-Type", "application/json")

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, req)

	if http.StatusOK != recorder.Code {
		t.Fail()
	}

	cookie := recorder.Result().Cookies()[0]

	if cookie.Name != "Authorization" {
		t.Fail()
	}
}

func Test_Login_Fail_Missing_Cred(t *testing.T) {
	route := testApp.routes()
	route.POST("/login", testApp.login)

	req, _ := http.NewRequest("POST", "/login", nil)

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, req)

	if http.StatusBadRequest != recorder.Code {
		t.Fail()
	}
}

func Test_Login_Fail_Invalid_Cred(t *testing.T) {
	route := testApp.routes()
	route.POST("/login", testApp.login)

	var reqPayload struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	reqPayload.Username = "user1"
	reqPayload.Password = "wrong-password"

	reqBytes, _ := json.Marshal(reqPayload)

	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(reqBytes))

	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, req)

	if http.StatusBadRequest != recorder.Code {
		t.Fail()
	}
}

func Test_Logout(t *testing.T) {
	route := testApp.routes()
	route.POST("/logout", testApp.logout)

	req, _ := http.NewRequest("POST", "/logout", nil)
	recorder := httptest.NewRecorder()

	route.ServeHTTP(recorder, req)

	if http.StatusOK != recorder.Code {
		t.Fail()
	}
}
