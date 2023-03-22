package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func Test_LoginSuccess(t *testing.T) {

	var testPayload = map[string]any{
		"username": "username",
		"password": "password",
	}

	testPayloadBytes, err := json.Marshal(testPayload)

	if err != nil {
		t.Errorf("Failed to marshal testPayload: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, "/v1/login", bytes.NewReader(testPayloadBytes))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Errorf("Failed to create post request: %s", err.Error())
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusOK {
		t.Errorf("FAILED: Expected 200 get %d", reqRecorder.Code)
	}

	/*
		Testing Authorization Cookie Presence
	*/
	cookies := reqRecorder.Result().Cookies()[0]
	allCookies := cookies.Raw
	if !strings.Contains(allCookies, "Authorization") {
		t.Error("FAILED: Authorization Cookie absent")
	}
}

func Test__LoginBadRequest(t *testing.T) {

	var testPayload = map[string]any{} // empty request payload

	testPayloadBytes, err := json.Marshal(testPayload)

	if err != nil {
		t.Errorf("Failed to marshal testPayload: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, "/v1/login", bytes.NewReader(testPayloadBytes))
	req.Header.Add("Content-Type", "application/json")

	if err != nil {
		t.Errorf("Failed to create post request: %s", err.Error())
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusBadRequest {
		t.Errorf("FAILED: Expected %d get %d", http.StatusBadRequest, reqRecorder.Code)
	}
}

func Test_Logout(t *testing.T) {

	req, err := http.NewRequest(http.MethodPost, "/v1/logout", nil)

	if err != nil {
		t.Errorf("Failed to create request: %s", err.Error())
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusOK {
		t.Errorf("FAILED: Expected 200 get %d", reqRecorder.Code)
	}

	/*
		Testing Authorization Cookie Absence
	*/
	cookies := reqRecorder.Result().Cookies()[0]
	allCookies := cookies.Raw

	// means auth cookies is emptied out
	if !strings.Contains(allCookies, "Authorization=;") {
		t.Error("FAILED: Authorization Cookie still present after logout")
	}
}

/*
Testing GET /v1/users

	-> Success
*/
func Test_AllUsersSuccess(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/v1/users", nil)

	if err != nil {
		t.Errorf("Failed to create request: %s", err.Error())
	}
	// test jwt-auth-token
	jwtToken, err := getJWTTestToken()

	if err != nil {
		t.Errorf("Failed to generate test JWT Token: %s", err.Error())
	}

	req.Header.Set("Cookie", fmt.Sprintf("Authorization=%s", jwtToken))

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusOK {
		t.Errorf("FAILED: Expected %d get %d", http.StatusOK, reqRecorder.Code)
	}

}

/*
Testing GET /v1/users

	-> Missing JWT Auth Token in Cookies

This Endpoint gives every other user in the same org.
If the Authorization cookie is not set then this endpoint will return  un-authorized code
*/
func Test_AllUsersUnAuthorized(t *testing.T) {
	req, err := http.NewRequest(http.MethodGet, "/v1/users", nil)

	if err != nil {
		t.Errorf("Failed to create request: %s", err.Error())
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusUnauthorized {
		t.Errorf("FAILED: Expected %d get %d", http.StatusUnauthorized, reqRecorder.Code)
	}
}

/*
Testing  POST /v1/add
*/
func Test_AddUserSuccess(t *testing.T) {
	testPayload := map[string]any{
		"username": "user-to-add",
		"password": "user-password",
	}

	testPayloadBytes, err := json.Marshal(testPayload)

	if err != nil {
		t.Errorf("Failed to marshal request payload: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, "/v1/add", bytes.NewReader(testPayloadBytes))
	if err != nil {
		t.Errorf("Failed to create test JWT Token: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	// test jwt-token
	jwtToken, err := getJWTTestToken()
	req.Header.Set("Cookie", fmt.Sprintf("Authorization=%s", jwtToken))

	if err != nil {
		t.Errorf("Failed to create request: %s", err.Error())
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusOK {
		t.Errorf("FAILED: Expected %d get %d", http.StatusOK, reqRecorder.Code)
	}
}

/*
Testing  POST /v1/add

	-> Missing JWT Auth Token in Cookies
*/
func Test_AddUserUnAuthorized(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "/v1/add", nil)

	if err != nil {
		t.Errorf("Failed to create request: %s", err.Error())
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusUnauthorized {
		t.Errorf("FAILED: Expected %d get %d", http.StatusUnauthorized, reqRecorder.Code)
	}
}

/*
Testing  POST /v1/add

	-> Missing request payload to add user
*/
func Test_AddUserBadRequest_MissingCred(t *testing.T) {
	testPayload := map[string]any{} // missing username and password

	testPayloadBytes, err := json.Marshal(testPayload)

	if err != nil {
		t.Errorf("Failed to marshal request payload: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodPost, "/v1/add", bytes.NewReader(testPayloadBytes))
	if err != nil {
		t.Errorf("Failed to create test JWT Token: %s", err.Error())
	}

	// test jwt-token
	jwtToken, err := getJWTTestToken()
	req.Header.Set("Cookie", fmt.Sprintf("Authorization=%s", jwtToken))

	if err != nil {
		t.Errorf("Failed to create request: %s", err.Error())
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusBadRequest {
		t.Errorf("FAILED: Expected %d get %d", http.StatusBadRequest, reqRecorder.Code)
	}
}

/*
Testing  DELETE /v1/delete
*/
func Test_DeleteUserSuccess(t *testing.T) {
	testPayload := map[string]any{
		"username": "user-to-add",
	}

	testPayloadBytes, err := json.Marshal(testPayload)

	if err != nil {
		t.Errorf("Failed to marshal request payload: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodDelete, "/v1/delete", bytes.NewReader(testPayloadBytes))
	if err != nil {
		t.Errorf("Failed to create test JWT Token: %s", err.Error())
	}

	req.Header.Set("Content-Type", "application/json")

	// test jwt-token
	jwtToken, err := getJWTTestToken()
	req.Header.Set("Cookie", fmt.Sprintf("Authorization=%s", jwtToken))

	if err != nil {
		t.Errorf("Failed to create request: %s", err.Error())
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusOK {
		t.Errorf("FAILED: Expected %d get %d", http.StatusOK, reqRecorder.Code)
	}
}

/*
Testing  DELETE /v1/delete

	-> Missing JWT Auth Token in Cookies
*/
func Test_DeleteUserUnAuthorized(t *testing.T) {
	req, err := http.NewRequest(http.MethodDelete, "/v1/delete", nil)

	if err != nil {
		t.Errorf("Failed to create request: %s", err.Error())
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusUnauthorized {
		t.Errorf("FAILED: Expected %d get %d", http.StatusUnauthorized, reqRecorder.Code)
	}
}

/*
Testing  DELETE /v1/delete

	-> Missing request payload to delete user
*/
func Test_DeleteUserBadRequest_MissingCred(t *testing.T) {
	testPayload := map[string]any{} // missing username

	testPayloadBytes, err := json.Marshal(testPayload)

	if err != nil {
		t.Errorf("Failed to marshal request payload: %s", err.Error())
	}

	req, err := http.NewRequest(http.MethodDelete, "/v1/delete", bytes.NewReader(testPayloadBytes))
	if err != nil {
		t.Errorf("Failed to create test JWT Token: %s", err.Error())
	}

	// test jwt-token
	jwtToken, err := getJWTTestToken()
	req.Header.Set("Cookie", fmt.Sprintf("Authorization=%s", jwtToken))

	if err != nil {
		t.Errorf("Failed to create request: %s", err.Error())
	}

	reqRecorder := httptest.NewRecorder()

	router.ServeHTTP(reqRecorder, req)

	if reqRecorder.Code != http.StatusBadRequest {
		t.Errorf("FAILED: Expected %d get %d", http.StatusBadRequest, reqRecorder.Code)
	}
}

/*
Function to get a dummy JWT Token for testing AllUsers endpoint
*/
func getJWTTestToken() (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userId": "test-user-id",
		"exp":    time.Now().Add(time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(JWT_SECRET))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}
